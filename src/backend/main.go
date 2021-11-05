package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/atoyr/passed/models"
	"github.com/urfave/cli/v2"
)

const APPNAME = "passed"

// SQL DB Information
var sqlserver string
var instance string
var user string
var password string
var database string

func main() {
	app := new(cli.App)
	app.Name = APPNAME
	database = "passed"
	instance = "..."
	sqlserver = "..."
	user = "..."
	password = "..."

	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/signup", signupHandler)
	http.ListenAndServe(":8080", nil)

	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	user := models.User{}
	user.Signup(d)
}

func connectionstring() string {
	var ret = make([]byte, 0, 1024)
	ret = append(ret, "server="...)
	ret = append(ret, sqlserver...)
	if instance != "" {
		ret = append(ret, "\\"...)
		ret = append(ret, instance...)
	}
	ret = append(ret, ";user id="...)
	ret = append(ret, user...)
	ret = append(ret, ";password="...)
	ret = append(ret, password...)
	ret = append(ret, ";database="...)
	ret = append(ret, database...)
	return string(ret)
}
