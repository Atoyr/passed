package Models

import (
	"encoding/hex"
	"encoding/pem"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"golang.org/x/crypto/sha3"
	"github.com/google/uuid"
)

type Signup struct {
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Nickname   string    `json:"nickname"`
}

// Signup is create account, profile and sec key
func (s *Signup) Signup() (Signin, error) {
	// AES Key   : hashed password + uuid
	// Signature : AES Encript (email)
	// Private   : AES Encript (RSA.PrivateKey)
	// Public    : RSA PublicKey

	// Generate UUID
	hash := make([]byte, 64)
	sh := sha3.NewShake256()
	sh.Write([]byte(s.Password))
	sh.Read(hash)
	fmt.Println(hex.EncodeToString(hash))
	u := uuid.NewSHA1(uuid.New(), hash)
	fmt.Println(u.String())

	passhash := make([]byte, 64)
	passshake := sha3.NewShake256()
	passshake.Write(hash)
	passshake.Write(u)
	passshake.Read(passhash)
	passphrase := hex.EncodeToString(passhash)

	// Generate AES Key
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}


	// Generate Private/public key
	// size of key (bits)
	size := 2048
	// Generate private and public key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}
	// Get public key from private key and encrypt
	publicKey := &privateKey.PublicKey




}
