package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func sample() {
	fmt.Println("Hello, playground")
	plainText := []byte("Bob loves Alice.")

	// size of key (bits)
	size := 2048

	// Generate private and public key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}

	fmt.Printf("privateKey: %s", ExportRsaPrivateKeyAsPemStr(privateKey))

	// Get public key from private key and encrypt
	publicKey := &privateKey.PublicKey

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}
	fmt.Printf("Cipher text: %x\n", cipherText)

	// Decrypt with private key
	decryptedText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return
	}

	fmt.Printf("Decrypted text: %s\n", decryptedText)
}

func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) string {
    privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
    privkey_pem := pem.EncodeToMemory(
            &pem.Block{
                    Type:  "RSA PRIVATE KEY",
                    Bytes: privkey_bytes,
            },
    )
    return string(privkey_pem)
}
