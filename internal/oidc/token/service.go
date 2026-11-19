package token

import (
	"context"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"errors"

	"Guardian_/internal/oidc/jwks"
	"Guardian_/internal/store"
	"Guardian_/internal/store/entity"
	"time"
)

var ErrInvalidGrant = errors.New("invalid_grant")

type Service struct {
	AuthCodeStore store.AuthCodeStore
	JWKS          *jwks.Service
	Issuer        string
}

func NewService(authCodeStore store.AuthCodeStore, jwksService *jwks.Service, issuer string) *Service {
	return &Service{
		AuthCodeStore: authCodeStore,
		JWKS:          jwksService,
		Issuer:        issuer,
	}
}

func (s *Service) ExchangeToken(ctx context.Context, code string) (*Result, error) {
	authCode, err := s.AuthCodeStore.GetAuthCode(ctx, code)
	if err != nil {
		if errors.Is(err, entity.ErrAuthCodeNotFound) {
			return nil, ErrInvalidGrant
		}
		return nil, err
	}

	// Single-use: consume the code immediately so it can't be replayed.
	if err := s.AuthCodeStore.DeleteAuthCode(ctx, code); err != nil {
		return nil, err
	}

	issuedAt := time.Now().UTC()
	accessToken, err := newOpaqueToken()
	if err != nil {
		return nil, err
	}

	idTokenClaims := map[string]any{
		"iss": s.Issuer,
		"sub": authCode.UserID,
		"aud": authCode.ClientID,
		"exp": issuedAt.Add(time.Hour).Unix(),
		"iat": issuedAt.Unix(),
	}
	if authCode.Nonce != "" {
		idTokenClaims["nonce"] = authCode.Nonce
	}

	idToken, err := s.JWKS.SignJWT(idTokenClaims)
	if err != nil {
		return nil, err
	}

	return &Result{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		IDToken:     idToken,
		Scope:       authCode.Scope,
	}, nil
}

func newOpaqueToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := cryptoRand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
