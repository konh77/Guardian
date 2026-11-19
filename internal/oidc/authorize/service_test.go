package authorize

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"Guardian_/internal/store/entity"
	"Guardian_/internal/store/memory"
	"Guardian_/internal/userinfo/source"
)

func TestAuthorizeIssuesCodeForRegisteredUserAndRedirect(t *testing.T) {
	clients := memory.NewClientStore()
	users := source.NewMemoryProvider()
	authCodes := memory.NewAuthCodeStore()
	nonces := memory.NewNonceStore()
	svc := NewService(clients, users, authCodes, nonces)

	ctx := context.Background()
	if err := clients.SaveClient(ctx, &entity.Client{
		ID:           "my-client",
		RedirectURIs: []string{"http://localhost:3000/callback"},
		ClientType:   "confidential",
	}); err != nil {
		t.Fatalf("save client: %v", err)
	}
	if err := users.RegisterUser(ctx, "user123", map[string]any{
		"sub":  "user123",
		"name": "Test User",
	}); err != nil {
		t.Fatalf("register user: %v", err)
	}

	result, errType, description, err := svc.Authorize(ctx, &AuthorizeRequest{
		ClientID:     "my-client",
		RedirectURI:  "http://localhost:3000/callback",
		ResponseType: "code",
		Scope:        "openid profile",
		State:        "abc123",
		UserID:       "user123",
	})
	if err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if errType != ErrorNone {
		t.Fatalf("unexpected error type %q: %s", errType, description)
	}

	redirectURL, err := url.Parse(result.RedirectURL)
	if err != nil {
		t.Fatalf("parse redirect url: %v", err)
	}
	if redirectURL.Query().Get("state") != "abc123" {
		t.Fatalf("expected state abc123, got %q", redirectURL.Query().Get("state"))
	}
	if redirectURL.Query().Get("code") == "" {
		t.Fatal("expected authorization code in redirect")
	}

	if !strings.HasPrefix(result.RedirectURL, "http://localhost:3000/callback") {
		t.Fatalf("expected redirect to callback, got %q", result.RedirectURL)
	}
}
