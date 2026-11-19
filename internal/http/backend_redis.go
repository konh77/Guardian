//go:build redis
// +build redis

package http

import (
	"os"

	"Guardian_/internal/store/redis"
	"Guardian_/internal/userinfo/source"

	goredis "github.com/redis/go-redis/v9"
)

func newBackend() backend {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	rc := goredis.NewClient(&goredis.Options{Addr: addr})

	return backend{
		clients:   redis.NewClientStore(rc),
		users:     source.NewRedisProvider(rc),
		authCodes: redis.NewAuthCodeStore(rc),
		nonces:    redis.NewNonceStore(rc),
	}
}
