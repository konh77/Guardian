package memory

import (
	"context"
	"sync"
	"time"

	"Guardian_/internal/store/entity"
)

type AuthCodeStore struct {
	mu    sync.Mutex
	codes map[string]*entity.AuthCode
}

func NewAuthCodeStore() *AuthCodeStore {
	return &AuthCodeStore{codes: make(map[string]*entity.AuthCode)}
}

func (s *AuthCodeStore) SaveAuthCode(ctx context.Context, code *entity.AuthCode) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.purgeExpiredLocked()
	s.codes[code.Code] = code
	return nil
}

func (s *AuthCodeStore) GetAuthCode(ctx context.Context, code string) (*entity.AuthCode, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	authCode, exists := s.codes[code]
	if !exists {
		return nil, entity.ErrAuthCodeNotFound
	}
	if time.Now().After(authCode.ExpiresAt) {
		delete(s.codes, code)
		return nil, entity.ErrAuthCodeNotFound
	}
	return authCode, nil
}

func (s *AuthCodeStore) DeleteAuthCode(ctx context.Context, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.codes, code)
	return nil
}

// purgeExpiredLocked drops expired codes; there's no background TTL sweeper here,
// so cleanup piggybacks on writes instead. Caller must hold s.mu.
func (s *AuthCodeStore) purgeExpiredLocked() {
	now := time.Now()
	for code, authCode := range s.codes {
		if now.After(authCode.ExpiresAt) {
			delete(s.codes, code)
		}
	}
}
