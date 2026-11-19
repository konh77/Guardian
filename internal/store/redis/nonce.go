//go:build redis
// +build redis

package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"Guardian_/internal/store/entity"
)

type NonceStore struct {
	client *redis.Client
}

func NewNonceStore(client *redis.Client) *NonceStore {
	return &NonceStore{
		client: client,
	}
}

func (s *NonceStore) MarkNonceSeen(ctx context.Context, nonce, clientID string) error {
	key := fmt.Sprintf("nonce:%s:%s", clientID, nonce)

	// SETNX with TTL (atomic operation)
	result := s.client.SetNX(ctx, key, "seen", entity.NonceTTL)
	if result.Err() != nil {
		return result.Err()
	}

	if !result.Val() {
		return entity.ErrNonceAlreadySeen // Key already exists = replay attack
	}

	return nil // Success: nonce marked as seen for first time
}
