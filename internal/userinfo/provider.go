package userinfo

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Provider interface {
	Fetch(ctx context.Context, sub string) (map[string]any, error)
}
