package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

const APPNAME = "passed"
const VERSION = "0.0.1"

// SQL DB Information
var sqlserver string
var instance string
var user string
var password string
var database string
var webport int

func main() {
	app := new(cli.App)
	app.Name = APPNAME
	app.Version = VERSION
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "server",
			Aliases:     []string{"s"},
			Value:       "",
			Usage:       "SQLServer Server Name",
			EnvVars:     []string{"DBSERVER"},
			Destination: &sqlserver,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "instance",
			Aliases:     []string{"i"},
			Value:       "",
			Usage:       "SQLServer Server Instance Name",
			EnvVars:     []string{"DBINSTANCE"},
			Destination: &instance,
		},
		&cli.StringFlag{
			Name:        "user",
			Aliases:     []string{"u"},
			Value:       "sa",
			Usage:       "SQLServer Server User",
			EnvVars:     []string{"DBUSER"},
			Destination: &user,
		},
		&cli.StringFlag{
			Name:        "password",
			Aliases:     []string{"p"},
			Value:       "",
			Usage:       "SQLServer Server Password",
			EnvVars:     []string{"DBPASS"},
			Destination: &password,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "database",
			Aliases:     []string{"d"},
			Value:       "master",
			Usage:       "SQLServer Server using database",
			Destination: &database,
		},
		&cli.IntFlag{
			Name:        "httpport",
			Aliases:     []string{"hp"},
			Value:       8080,
			Usage:       "http access port no",
			Destination: &webport,
		},
	}

	app.Action = action
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

	// d, err := sql.Open("sqlserver", connectionstring())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer d.Close()

}

func action(c *cli.Context) error {
	// DBチェック
	d, err := sql.Open("sqlserver", connectionstring())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer d.Close()
	if err := d.Ping(); err != nil {
		log.Fatal(err)
		return err
	}

	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/singin", signinHandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", webport), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
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
