package main

import (
	"net/http"

	"github.com/urfave/cli/v2"
)

const APPNAME = "passed"

func main() {
	app := new(cli.App)
	app.Name = APPNAME
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/signup", signupHandler)
	http.ListenAndServe(":8080", nil)
}
