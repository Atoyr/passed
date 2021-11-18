package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"

	"golang.org/x/crypto/sha3"
)

func AesEncript(key []byte, src []byte) ([]byte, error) {
	// Generate AES Key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create IV
	chiperBytes := make([]byte, aes.BlockSize+len(src))
	iv := chiperBytes[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	copy(chiperBytes, iv)

	// Encrypt
	encryptStream := cipher.NewCTR(block, iv)
	encryptStream.XORKeyStream(chiperBytes[aes.BlockSize:], src)

	return chiperBytes, nil
}

func AesDecript(key []byte, src []byte) ([]byte, error) {
	// Generate AES Key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Decrpt
	descripted := make([]byte, len(src[aes.BlockSize:]))
	// src[:aes.BlockSize] is iv
	decryptStream := cipher.NewCTR(block, src[:aes.BlockSize])
	decryptStream.XORKeyStream(descripted, src[aes.BlockSize:])
	return descripted, nil
}

func GetSha3Hash(u ...[]byte) []byte {
	hash := make([]byte, 32)
	shake := sha3.NewShake128()
	for _, v := range u {
		shake.Write(v)
	}
	shake.Read(hash)
	return hash
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToString(priv *rsa.PrivateKey) string {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	return string(privBytes)
}

// PublicKeyToBytes public key to bytes
func PublicKeyToString(pub *rsa.PublicKey) string {
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pub),
	})
	return string(pubBytes)
}
