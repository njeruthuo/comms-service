package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func DatabaseConnect() (*sql.DB, error) {
	var (
		DatabaseName     = os.Getenv("DatabaseName")
		DatabaseUser     = os.Getenv("DatabaseUser")
		DatabaseHost     = os.Getenv("DatabaseHost")
		DatabasePassword = os.Getenv("DatabasePassword")
	)

	dbInfo := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		DatabaseUser, DatabasePassword, DatabaseHost, DatabaseName,
	)

	var err error
	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
