package jwks

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"

	"github.com/google/uuid"
)

type Service struct {
	privateKey *rsa.PrivateKey
	keyID      string
}

func NewService() *Service {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return &Service{
		privateKey: privateKey,
		keyID:      uuid.NewString(),
	}
}

func (s *Service) GetKeySet() *KeySet {
	return &KeySet{
		Keys: []Key{s.publicKey()},
	}
}

func (s *Service) SignJWT(claims map[string]any) (string, error) {
	header := map[string]any{
		"alg": "RS256",
		"kid": s.keyID,
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := encodedHeader + "." + encodedClaims

	digest := crypto.SHA256.New()
	if _, err := digest.Write([]byte(signingInput)); err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, digest.Sum(nil))
	if err != nil {
		return "", err
	}

	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func (s *Service) KeyID() string {
	return s.keyID
}

func (s *Service) PublicKey() *rsa.PublicKey {
	return &s.privateKey.PublicKey
}

func (s *Service) publicKey() Key {
	pub := &s.privateKey.PublicKey
	return Key{
		Kty: "RSA",
		Use: "sig",
		Kid: s.keyID,
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}
}
