package store

import "context"
import "Guardian_/internal/store/entity"

type ClientStore interface {
	GetClient(ctx context.Context, clientID string) (*entity.Client, error)
}

type Stores struct {
	Client ClientStore
}
