package models

import (
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

	// Generate private and public key pair
	privateKey, err := GetPrivateKey()
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
		account := database.Account{}
		account.Email = s.Email
		account.Signature = signature
		account.Private = private
		account.Public = x509.MarshalPKCS1PublicKey(publicKey)
		account.ValidFlg = true
		account.InsertSystemID = SIGNUP_SYSTEMID
		account.UpdateSystemID = SIGNUP_SYSTEMID

		profile := database.Profile{}
		profile.FirstName = s.FirstName
		profile.MiddleName = s.MiddleName
		profile.LastName = s.LastName
		profile.Nickname = s.Nickname
		profile.InsertSystemID = SIGNUP_SYSTEMID
		profile.UpdateSystemID = SIGNUP_SYSTEMID
		tx, err := dbcontext.Begin()
		if err != nil {
			return authentication, err
		}

		accountwps := database.WherePhrases{}
		accountwps.Append(database.Equal, "Email", s.Email)
		accountwps.Append(database.Equal, "ValidFlg", true)
		accounts, err := database.GetAccounts(tx, accountwps)
		if err != nil {
			tx.Rollback()
			return authentication, err
		}
		if len(accounts) > 0 {
			tx.Rollback()
			return authentication, fmt.Errorf("Account is exists")
		}

		err = account.Insert(tx)
		if err != nil {
			tx.Rollback()
			return authentication, err
		}

		profile.AccountID = account.ID
		profile.InsertAccountID = account.ID
		profile.UpdateAccountID = account.ID

		err = profile.Insert(tx)
		if err != nil {
			tx.Rollback()
			return authentication, err
		}

		tx.Commit()
	}

	authentication.Email = s.Email
	authentication.Key = u.String()
	authentication.Password = s.Password

	return authentication, nil
}
