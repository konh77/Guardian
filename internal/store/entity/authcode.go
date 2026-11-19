package entity

import (
	"errors"
	"time"
)

var ErrAuthCodeNotFound = errors.New("auth code not found")

type AuthCode struct {
	Code                string
	ClientID            string
	UserID              string
	RedirectURI         string
	Scope               string
	State               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
	ExpiresAt           time.Time
}
