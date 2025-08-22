package tools

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func connect() (db *sql.DB, err error) {
	user, password, dbName, host, err := getEnvLoginData()
	if err != nil {
		return
	}

	loginInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)

	db, err = sql.Open("postgres", loginInfo)
	if err != nil {
		return
	}

	return
}

func checkDbExistence() (err error) {
	user, password, dbName, host, err := getEnvLoginData()
	if err != nil {
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
		_, err = tempConn.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
		if err != nil {
			return
		}
	}
	tempConn.Close()
	return
}
func getEnvLoginData() (user, password, dbName, host string, err error) {
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	host = os.Getenv("DB_HOST")
	if user == "" || password == "" || dbName == "" || host == "" {
		err = errors.New("one of the following variables is not specified in env: DB_USER, DB_PASSWORD, DB_NAME, DB_HOST")
		return
	}
	return
}
