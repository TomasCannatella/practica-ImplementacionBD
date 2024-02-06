package main

import (
	"app/internal/handler/application"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// env
	// ...

	// application
	// - config
	cfg := &application.ConfigDefault{
		Database: mysql.Config{
			User:      os.Getenv("DB_USER"),
			Passwd:    os.Getenv("DB_PASSWORD"),
			Net:       "tcp",
			Addr:      os.Getenv("DB_HOST"),
			DBName:    os.Getenv("DB_NAME"),
			ParseTime: true,
		},
		Address: "127.0.0.1:8080",
	}
	app := application.NewDefault(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
