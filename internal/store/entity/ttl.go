package entity

import "time"

const AuthCodeTTL = 5 * time.Minute

// NonceTTL must stay >= AuthCodeTTL so a nonce can't be replayed while its auth code is still valid.
const NonceTTL = AuthCodeTTL + 2*time.Minute
