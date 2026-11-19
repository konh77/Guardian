//go:build !redis
// +build !redis

package http

import (
	"Guardian_/internal/store/memory"
	"Guardian_/internal/userinfo/source"
)

func newBackend() backend {
	return backend{
		clients:   memory.NewClientStore(),
		users:     source.NewMemoryProvider(),
		authCodes: memory.NewAuthCodeStore(),
		nonces:    memory.NewNonceStore(),
	}
}
