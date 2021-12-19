package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/atoyr/passed/models"
	"github.com/atoyr/passed/anonymous"
)

var AnonymousKeyManager *anonymous.AnonymousKeyManager

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var retJson []byte
	var err error
	switch r.Method {
	case "GET":
		signup := models.Signup{}
		signup.Email = "example@example.com"
		signup.Password = "encript password"
		signup.FirstName = "firstName"
		signup.MiddleName = "middleName"
		signup.LastName = "lastName"
		signup.Nickname = "nickname"
		retJson, err = json.Marshal(signup)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case "POST":
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//Read body data to parse json
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//parse json
		var signup models.Signup
		if err := json.Unmarshal(body, &signup); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		d, err := sql.Open("sqlserver", connectionstring())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer d.Close()

		auth, err := signup.Signup(d)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		retJson, err = json.Marshal(auth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(retJson)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func AnonymousKeyHandler(w http.ResponseWriter, r *http.Request) {
	if (AnonymousKeyManager == nil) {
		AnonymousKeyManager = new(anonymous.AnonymousKeyManager)
	}
	var retJson []byte
	switch r.Method {
	case "GET":
		ip, err := GetIP(r)
		if (err != nil) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		anonymousKey, err := AnonymousKeyManager.Get(ip)
		if (err != nil) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		anonymousPublicKey := anonymousKey.CreateAnonymousPublicKey()

		retJson, err = json.Marshal(anonymousPublicKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(retJson)
		return
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
	}
}
