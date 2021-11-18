package models

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"fmt"

	"github.com/atoyr/passed/database"
	"github.com/google/uuid"
)

type Signup struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Nickname   string `json:"nickname"`
}

const SIGNUP_SYSTEMID string = "D35DFEBF-8443-4EB5-95CF-F31DAF27D0DD"

// Signup is create account, profile and sec key
func (s *Signup) Signup(dbcontext *sql.DB) (Authentication, error) {
	// AES Key   : hashed(hashed password + uuid)
	// Signature : AES Encript (email)
	// Private   : AES Encript (RSA.PrivateKey)
	// Public    : RSA PublicKey

	authentication := Authentication{}

	// Generate UUID
	hash := GetSha3Hash([]byte(s.Password))
	u := uuid.NewSHA1(uuid.New(), hash)
	passhash := GetSha3Hash(hash, u.NodeID())

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

	// Generate Signature
	signature, err := AesEncript(passhash, []byte(s.Email))
	if err != nil {
		return authentication, err
	}
	// Generate private
	private, err := AesEncript(passhash, x509.MarshalPKCS1PrivateKey(privateKey))
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
		profile.InsertSystemID = SIGNUP_SYSTEMID
		profile.UpdateSystemID = SIGNUP_SYSTEMID

		account := database.Account{}
		account.Signature = signature
		account.Private = private
		account.Public = x509.MarshalPKCS1PublicKey(publicKey)
		account.ValidFlg = true
		account.InsertSystemID = SIGNUP_SYSTEMID
		account.UpdateSystemID = SIGNUP_SYSTEMID

		tx, err := dbcontext.Begin()
		if err != nil {
			return authentication, err
		}
		err = profile.Insert(tx)
		if err != nil {
			tx.Rollback()
			return authentication, err
		}
		account.ProfileID = profile.ID
		account.InsertProfileID = account.ProfileID
		account.UpdateProfileID = account.ProfileID
		err = account.Insert(tx)
		if err != nil {
			tx.Rollback()
			return authentication, err
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
