package source

import (
	"context"
	"sync"

	"Guardian_/internal/userinfo"
)

type MemoryProvider struct {
	mu    sync.RWMutex
	users map[string]map[string]any
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{
		users: make(map[string]map[string]any),
	}
}

func (p *MemoryProvider) Fetch(ctx context.Context, sub string) (map[string]any, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	user, ok := p.users[sub]
	if !ok {
		return nil, userinfo.ErrUserNotFound
	}
	return user, nil
}

func (p *MemoryProvider) RegisterUser(ctx context.Context, sub string, claims map[string]any) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.users[sub] = claims
	return nil
}
