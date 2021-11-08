package models

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/pem"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

type Signup struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Nickname   string `json:"nickname"`
}

// Signup is create account, profile and sec key
func (s *Signup) Signup(dbcontext *sql.DB) (Authentication, error) {
	// AES Key   : hashed(hashed password + uuid)
	// Signature : AES Encript (email)
	// Private   : AES Encript (RSA.PrivateKey)
	// Public    : RSA PublicKey

	authentication := Authentication{}

	// Generate UUID
	hash := make([]byte, 64)
	sh := sha3.NewShake256()
	sh.Write([]byte(s.Password))
	sh.Read(hash)
	fmt.Println(hex.EncodeToString(hash))
	u := uuid.NewSHA1(uuid.New(), hash)

	passhash := make([]byte, 64)
	passshake := sha3.NewShake256()
	passshake.Write(hash)
	passshake.Write(u.NodeID())
	passshake.Read(passhash)
	// Passohrase
	// passphrase := hex.EncodeToString(passhash)

	// Generate AES Key
	block, err := aes.NewCipher(passhash)
	if err != nil {
		return authentication, err
	}
	signature := make([]byte, len(s.Email))
	block.Encrypt([]byte(s.Email), signature)

	// Generate Private/public key
	// size of key (bits)
	size := 2048
	// Generate private and public key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return authentication, err
	}
	// Get public key from private key and encrypt
	public := &privateKey.PublicKey
	private := ExportRsaPrivateKeyAsPemStr(privateKey)

	user := User{}
	user.Email = s.Email
	user.FirstName = s.FirstName
	user.MiddleName = s.MiddleName
	user.LastName = s.LastName
	user.Nickname = s.Nickname

	authentication.Email = s.Email
	authentication.Key = u.String()
	authentication.Password = s.Password

	return authentication, nil
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
