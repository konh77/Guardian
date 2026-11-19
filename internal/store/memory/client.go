package memory

import (
	"context"
	"sync"

	"Guardian_/internal/store/entity"
)

type ClientStore struct {
	mu      sync.RWMutex
	clients map[string]*entity.Client
}

func NewClientStore() *ClientStore {
	return &ClientStore{clients: make(map[string]*entity.Client)}
}

func (s *ClientStore) GetClient(ctx context.Context, clientID string) (*entity.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	client, ok := s.clients[clientID]
	if !ok {
		return nil, entity.ErrClientNotFound
	}
	return client, nil
}

func (s *ClientStore) SaveClient(ctx context.Context, client *entity.Client) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[client.ID] = client
	return nil
}
