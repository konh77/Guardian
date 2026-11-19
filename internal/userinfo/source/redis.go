package source

import (
	"context"
	"encoding/json"
	"errors"

	"Guardian_/internal/userinfo"

	goredis "github.com/redis/go-redis/v9"
)

const userKeyPrefix = "user:"

type RedisProvider struct {
	client *goredis.Client
}

func NewRedisProvider(client *goredis.Client) *RedisProvider {
	return &RedisProvider{client: client}
}

func (p *RedisProvider) Fetch(ctx context.Context, sub string) (map[string]any, error) {
	data, err := p.client.Get(ctx, userKeyPrefix+sub).Bytes()
	if err != nil {
		if errors.Is(err, goredis.Nil) {
			return nil, userinfo.ErrUserNotFound
		}
		return nil, err
	}

	var claims map[string]any
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, err
	}
	return claims, nil
}

func (p *RedisProvider) RegisterUser(ctx context.Context, sub string, claims map[string]any) error {
	data, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	return p.client.Set(ctx, userKeyPrefix+sub, data, 0).Err()
}
