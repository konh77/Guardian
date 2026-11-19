package authorize

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strings"
	"time"

	"Guardian_/internal/store"
	"Guardian_/internal/store/entity"
	"Guardian_/internal/userinfo"

	"github.com/google/uuid"
)

// Service handles authorization business logic
type Service struct {
	ClientStore      store.ClientStore
	UserinfoProvider userinfo.Provider
	AuthCodeStore    store.AuthCodeStore
	NonceStore       store.NonceStore
}

// NewService creates a new authorization service
func NewService(
	clientStore store.ClientStore,
	userinfoProvider userinfo.Provider,
	authCodeStore store.AuthCodeStore,
	nonceStore store.NonceStore,
) *Service {
	return &Service{
		ClientStore:      clientStore,
		UserinfoProvider: userinfoProvider,
		AuthCodeStore:    authCodeStore,
		NonceStore:       nonceStore,
	}
}

// AuthorizeRequest represents the data needed for authorization
type AuthorizeRequest struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	Scope               string
	State               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
	UserID              string
}

// Authorize processes authorization request (pure business logic)
func (s *Service) Authorize(ctx context.Context, req *AuthorizeRequest) (*Result, ErrorType, string, error) {
	if req.ClientID == "" {
		return nil, ErrorInvalidRequest, "client_id is required", nil
	}

	// TODO: Step 2 - Validate state parameter (CSRF protection)
	if req.State == "" {
		return nil, ErrorInvalidRequest, "state is required", nil
	}
	fmt.Println("1")
	// Step 3 - Nonce replay protection
	if req.Nonce != "" {
		err := s.NonceStore.MarkNonceSeen(ctx, req.Nonce, req.ClientID)
		if errors.Is(err, entity.ErrNonceAlreadySeen) {
			return nil, ErrorInvalidRequest, "nonce has already been used", nil
		}
		if err != nil {
			return nil, ErrorServerError, "", err
		}
	}
	// Step 4 - Validate response_type
	fmt.Println("2")
	if req.ResponseType != "code" {
		return nil, ErrorUnsupportedResponse, "response_type must be \"code\"", nil
	}

	// Step 5 - Validate scope
	// - scope is a space-delimited list and MUST contain "openid"
	if !slices.Contains(strings.Fields(req.Scope), "openid") {
		return nil, ErrorInvalidScope, "scope must include \"openid\"", nil
	}
	// Step 6 - Get and validate client
	// - Use s.ClientStore.GetClient()
	// - Return ErrorInvalidClient if not found
	client, err := s.ClientStore.GetClient(ctx, req.ClientID)
	if err != nil {
		return nil, ErrorInvalidClient, "unknown client_id", nil
	}
	if !validateExactRedirectURI(client.RedirectURIs, req.RedirectURI) {
		return nil, ErrorInvalidRedirectURI, "redirect_uri must exactly match a registered URI", nil
	}

	if err := s.validatePKCEForClient(client, req); err != nil {
		return nil, ErrorInvalidRequest, err.Error(), nil
	}

	if req.UserID == "" {
		return nil, ErrorAccessDenied, "authenticated user is required", nil
	}
	if _, err := s.UserinfoProvider.Fetch(ctx, req.UserID); err != nil {
		return nil, ErrorAccessDenied, "authenticated user not found", nil
	}

	// Step 10 - Generate authorization code
	code := uuid.NewString()
	authCode := &entity.AuthCode{
		Code:                code,
		ClientID:            req.ClientID,
		UserID:              req.UserID,
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		State:               req.State,
		Nonce:               req.Nonce,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		ExpiresAt:           time.Now().Add(entity.AuthCodeTTL),
	}
	if err := s.AuthCodeStore.SaveAuthCode(ctx, authCode); err != nil {
		return nil, ErrorServerError, "", err
	}

	redirectURL, err := buildRedirectURL(req.RedirectURI, code, req.State)
	if err != nil {
		return nil, ErrorServerError, "", err
	}

	return &Result{RedirectURL: redirectURL}, ErrorNone, "", nil
}

// validateExactRedirectURI performs RFC 6749 compliant exact redirect URI validation
// SECURITY: Uses string comparison to prevent URL normalization attacks
func validateExactRedirectURI(registeredURIs []string, requestedURI string) bool {
	return slices.Contains(registeredURIs, requestedURI)
}

// validatePKCEForClient validates PKCE requirements based on client type
func (s *Service) validatePKCEForClient(client *entity.Client, req *AuthorizeRequest) error {
	// For public clients, PKCE is mandatory
	if client.IsPublicClient() {
		if req.CodeChallenge == "" {
			return errors.New("PKCE required for public clients")
		}
		if req.CodeChallengeMethod != "S256" {
			return errors.New("only S256 method supported for PKCE")
		}
		return nil
	}

	// For confidential clients, PKCE is optional
	if req.CodeChallenge != "" && req.CodeChallengeMethod != "S256" {
		return errors.New("only S256 method supported for PKCE")
	}

	return nil
}

func buildRedirectURL(redirectURI, code, state string) (string, error) {
	parsed, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Set("code", code)
	query.Set("state", state)
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}
