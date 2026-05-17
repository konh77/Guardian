package memory

import (
	"Guardian_/internal/store/entity"
)

type ClientStore struct{}

func NewClientStore() *ClientStore {
	return nil
}

func (s *ClientStore) RegisterClient(client *entity.Client) (*entity.Client, error) {
	return nil, nil
}
