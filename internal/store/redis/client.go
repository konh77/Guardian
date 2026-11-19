//go:build redis
// +build redis

package redis

import (
	"context"
	"encoding/json"

	"Guardian_/internal/store/entity"

	"github.com/redis/go-redis/v9"
)

type ClientStore struct {
	client *redis.Client
	prefix string
}

func NewClientStore(client *redis.Client) *ClientStore {
	return &ClientStore{
		client: client,
		prefix: "client:",
	}
}

func (s *ClientStore) GetClient(ctx context.Context, id string) (*entity.Client, error) {
	data, err := s.client.Get(ctx, s.prefix+id).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, entity.ErrClientNotFound
		}
		return nil, err
	}

	var client entity.Client
	if err := json.Unmarshal([]byte(data), &client); err != nil {
		return nil, err
	}
	return &client, nil
}

func (s *ClientStore) SaveClient(ctx context.Context, client *entity.Client) error {
	data, err := json.Marshal(client)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, s.prefix+client.ID, data, 0).Err()
}
