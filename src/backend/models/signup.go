package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/atoyr/passed/database"
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

	passhash := make([]byte, 32)
	passshake := sha3.NewShake128()
	passshake.Write(hash)
	passshake.Write(u.NodeID())
	passshake.Read(passhash)
	// Passohrase
	// passphrase := hex.EncodeToString(passhash)
	fmt.Printf("%x\n", passhash)

	// Generate Private/public key
	// size of key (bits)
	size := 2048
	// Generate private and public key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return authentication, err
	}
	// Get public key from private key and encrypt
	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)

	signature, err := AesEncript(passhash, []byte(s.Email))
	if err != nil {
		return authentication, err
	}

	private, err := AesEncript(passhash, privateKeyBytes)
	if err != nil {
		return authentication, err
	}

	if dbcontext != nil {
		// Insert Database
		profile := database.Profile{}
		profile.Email = s.Email
		profile.FirstName = s.FirstName
		profile.MiddleName = s.MiddleName
		profile.LastName = s.LastName
		profile.Nickname = s.Nickname
		profile.InsertSystemID = "models/signup"
		profile.ModifiedSystemID = "models/signup"

		account := database.Account{}
		account.Signature = signature
		account.Private = private
		account.Public = publicKeyBytes
		account.InsertSystemID = "models/signup"
		account.ModifiedSystemID = "models/signup"

		tx, err := dbcontext.Begin()
		if err != nil {
			return authentication, err
		}
		err = profile.Insert(tx)
		if err != nil {
			tx.Rollback()
		}
		account.ProfileID = profile.ID
		err = account.Insert(tx)
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	} else {
		signatureStr, err := AesDecript(passhash, signature)
		if err != nil {
			fmt.Println(err)
		}
		privateStr, err := AesDecript(passhash, private)
		if err != nil {
			fmt.Println(err)
		}
		pr, err := x509.ParsePKCS1PrivateKey(privateStr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s.Email)
		fmt.Println(u.String())
		fmt.Println()
		fmt.Println("--- sign ---")
		fmt.Printf("%x\n", signature)
		fmt.Printf("%s|\n", string(signatureStr))

		fmt.Println("--- private ---")
		fmt.Printf("%x\n", private)
		fmt.Printf("%s\n", PrivateKeyToString(pr))
		fmt.Println("--- result ---")
		fmt.Printf(PrivateKeyToString(privateKey))
		fmt.Printf(PublicKeyToString(publicKey))
	}

	authentication.Email = s.Email
	authentication.Key = u.String()
	authentication.Password = s.Password

	return authentication, nil
}
