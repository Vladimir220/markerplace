package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

func connect() (db *sql.DB, err error) {
	var (
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbName   = os.Getenv("DB_NAME")
		host     = os.Getenv("DB_HOST")
	)
	if user == "" || password == "" || dbName == "" || host == "" {
		err = errors.New("one of the following variables is not specified in env: DB_USER, DB_PASSWORD, DB_NAME, DB_HOST")
		return
	}

	loginInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, "postgres")

	tempConn, err := sql.Open("postgres", loginInfo)
	if err != nil {
		return
	}

	query := "SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"

	var isDbExist bool
	err = tempConn.QueryRow(query, dbName).Scan(&isDbExist)
	if err != nil {
		return
	}

	if !isDbExist {
		_, err = tempConn.Exec("CREATE DATABASE $1;", dbName)
		if err != nil {
			return
		}
	}
	tempConn.Close()

	loginInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)

	db, err = sql.Open("postgres", loginInfo)
	if err != nil {
		return
	}

	return
}
