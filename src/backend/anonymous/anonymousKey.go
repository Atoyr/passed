package anonymous

import (
	"crypto/rsa"
	"time"
)

type AnonymousKeyManager struct {
	Storage map[string]AnonymousKey
}

type AnonymousKey struct {
	PrivateKey     *rsa.PrivateKey
	CreateDatetime time.Time
}
