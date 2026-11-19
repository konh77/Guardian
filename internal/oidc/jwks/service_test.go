package jwks

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestServiceGetKeySetAndSignJWT(t *testing.T) {
	svc := NewService()

	keySet := svc.GetKeySet()
	if len(keySet.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keySet.Keys))
	}
	if keySet.Keys[0].Kid != svc.KeyID() {
		t.Fatalf("expected kid %q, got %q", svc.KeyID(), keySet.Keys[0].Kid)
	}

	token, err := svc.SignJWT(map[string]any{
		"iss": "http://localhost:8080",
		"sub": "user123",
		"aud": "my-client",
		"exp": 1234567890,
		"iat": 1234560000,
	})
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}

	parts := strings.Split(token, ".")
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
	if header["alg"] != "RS256" {
		t.Fatalf("expected alg RS256, got %v", header["alg"])
	}
	if header["kid"] != svc.KeyID() {
		t.Fatalf("expected kid %q, got %v", svc.KeyID(), header["kid"])
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
	if err := rsa.VerifyPKCS1v15(svc.PublicKey(), crypto.SHA256, digest.Sum(nil), sig); err != nil {
		t.Fatalf("verify signature: %v", err)
	}
}
