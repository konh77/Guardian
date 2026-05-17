package source

import (
	"context"
	"sync"
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

	user := p.users[sub]
	return user, nil
}

func (p *MemoryProvider) RegisterUser(sub string, claims map[string]any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.users[sub] = claims
}
