package anonymous

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"
)

type AnonymousKey struct {
	PrivateKey     *rsa.PrivateKey
	CreateDatetime time.Time
}

type AnonymousKeyManager struct {
	storage     map[string]AnonymousKey
	RefreshRate int
}

func NewAnonymousKey() AnonymousKey {
	a := AnonymousKey{}
	size := 2048
	privateKey, _ := rsa.GenerateKey(rand.Reader, size)
	a.PrivateKey = privateKey
	a.CreateDatetime = time.Now()
	return a
}

func NewAnonymousKeyManager() AnonymousKeyManager {
	am := AnonymousKeyManager{}
	am.storage = make(map[string]AnonymousKey)
	am.RefreshRate = 300
	return am
}

func (anonymousKeyManager *AnonymousKeyManager) NewAnonymousKey(key string) AnonymousKey {
	anonymousKey := NewAnonymousKey()
	anonymousKeyManager.storage[key] = anonymousKey
	return anonymousKey
}

func (anonymousKeyManager *AnonymousKeyManager) Get(key string) (AnonymousKey, error) {
	if anonymousKey, ok := anonymousKeyManager.storage[key]; ok {
		now := time.Now()
		if anonymousKey.CreateDatetime.After(now.Add(time.Duration(-anonymousKeyManager.RefreshRate) * time.Second)) {
			return AnonymousKey{}, fmt.Errorf("AnonymousKey not validate")
		}
		return anonymousKey, nil
	} else {
		return anonymousKey, fmt.Errorf("AnonymousKey not found")
	}
}
