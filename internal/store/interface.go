package store

import "context"
import "Guardian_/internal/store/entity"

type ClientStore interface {
	GetClient(ctx context.Context, clientID string) (*entity.Client, error)
}

type AuthCodeStore interface {
	SaveAuthCode(ctx context.Context, code *entity.AuthCode) error
	GetAuthCode(ctx context.Context, code string) (*entity.AuthCode, error)
	DeleteAuthCode(ctx context.Context, code string) error
}

type NonceStore interface {
	MarkNonceSeen(ctx context.Context, nonce, clientID string) error
}

type Stores struct {
	Client   ClientStore
	AuthCode AuthCodeStore
	Nonce    NonceStore
}
