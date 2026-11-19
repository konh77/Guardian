package token

import (
	"context"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"Guardian_/internal/oidc/jwks"
	"Guardian_/internal/store/entity"
	"Guardian_/internal/store/memory"
)

func TestExchangeTokenReturnsSignedIDToken(t *testing.T) {
	authCodes := memory.NewAuthCodeStore()
	keyService := jwks.NewService()
	svc := NewService(authCodes, keyService, "http://localhost:8080")

	ctx := context.Background()
	if err := authCodes.SaveAuthCode(ctx, &entity.AuthCode{
		Code:        "auth-code-1",
		ClientID:    "my-client",
		UserID:      "user123",
		Scope:       "openid profile",
		Nonce:       "nonce-1",
		ExpiresAt:   time.Now().Add(time.Minute),
		RedirectURI: "http://localhost:3000/callback",
	}); err != nil {
		t.Fatalf("save auth code: %v", err)
	}

	result, err := svc.ExchangeToken(ctx, "auth-code-1")
	if err != nil {
		t.Fatalf("exchange token: %v", err)
	}

	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if result.IDToken == "" {
		t.Fatal("expected id token")
	}
	if result.TokenType != "Bearer" {
		t.Fatalf("expected token type Bearer, got %q", result.TokenType)
	}

	parts := strings.Split(result.IDToken, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3 token parts, got %d", len(parts))
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("decode header: %v", err)
	}
	var header map[string]any
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		t.Fatalf("unmarshal header: %v", err)
	}
	if header["kid"] != keyService.KeyID() {
		t.Fatalf("expected kid %q, got %v", keyService.KeyID(), header["kid"])
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("decode claims: %v", err)
	}
	var claims map[string]any
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		t.Fatalf("unmarshal claims: %v", err)
	}
	if claims["iss"] != "http://localhost:8080" {
		t.Fatalf("expected issuer claim, got %v", claims["iss"])
	}
	if claims["sub"] != "user123" {
		t.Fatalf("expected sub user123, got %v", claims["sub"])
	}
	if claims["aud"] != "my-client" {
		t.Fatalf("expected aud my-client, got %v", claims["aud"])
	}
	if claims["nonce"] != "nonce-1" {
		t.Fatalf("expected nonce nonce-1, got %v", claims["nonce"])
	}

	signingInput := parts[0] + "." + parts[1]
	digest := crypto.SHA256.New()
	if _, err := digest.Write([]byte(signingInput)); err != nil {
		t.Fatalf("hash token: %v", err)
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatalf("decode signature: %v", err)
	}
	if err := rsa.VerifyPKCS1v15(keyService.PublicKey(), crypto.SHA256, digest.Sum(nil), sig); err != nil {
		t.Fatalf("verify signature: %v", err)
	}
}
