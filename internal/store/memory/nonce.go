package memory

import (
	"context"
	"sync"
	"time"

	"Guardian_/internal/store/entity"
)

type NonceStore struct {
	mu     sync.Mutex
	nonces map[string]time.Time // value = expiry time
}

func NewNonceStore() *NonceStore {
	return &NonceStore{nonces: make(map[string]time.Time)}
}

func (s *NonceStore) MarkNonceSeen(ctx context.Context, nonce, clientID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for k, expiresAt := range s.nonces {
		if now.After(expiresAt) {
			delete(s.nonces, k)
		}
	}

	key := clientID + ":" + nonce
	if expiresAt, seen := s.nonces[key]; seen && now.Before(expiresAt) {
		return entity.ErrNonceAlreadySeen
	}
	s.nonces[key] = now.Add(entity.NonceTTL)
	return nil
}
