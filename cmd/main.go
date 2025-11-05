package main

import (
	"database/sql"
	"log"

	"github.com/eugenius-watchman/ecom_go_rest_api/cmd/api"
	"github.com/eugenius-watchman/ecom_go_rest_api/config"
	"github.com/eugenius-watchman/ecom_go_rest_api/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	// get DB
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	initStorage(db)

	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

// initialise DB connection
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}