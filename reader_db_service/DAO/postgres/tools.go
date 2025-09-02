package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"reader_db_service/env"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func connect() (db *sql.DB, err error) {
	data, err := env.GetPostgresEnvData()
	if err != nil {
		return
	}

	loginInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", data.User, data.Password, data.Host, data.DbName)

	db, err = sql.Open("postgres", loginInfo)
	if err != nil {
		return
	}

	return
}

func checkDbExistence() (err error) {
	logLabel := "checkDbExistence():"

	data, err := env.GetPostgresEnvData()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		return
	}

	loginInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", data.User, data.Password, data.Host, "postgres")

	tempConn, err := sql.Open("postgres", loginInfo)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	query := "SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"

	var isDbExist bool
	err = tempConn.QueryRow(query, data.DbName).Scan(&isDbExist)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	if !isDbExist {
		_, err = tempConn.Exec(fmt.Sprintf("CREATE DATABASE %s;", data.DbName))
		if err != nil {
			err = fmt.Errorf("%s %v", logLabel, err)
			return
		}
	}
	tempConn.Close()
	return
}

func CheckMigrations(connection *sql.DB) error {
	logLabel := "checkMigrations():"

	if connection == nil {
		return fmt.Errorf("%s no connection", logLabel)
	}

	driver, err := postgres.WithInstance(connection, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s checkMigrations(): %v", logLabel, err)
	}

	dbName := os.Getenv("DB_NAME")
	var path string
	currentDir, _ := os.Getwd()
	path = filepath.ToSlash(currentDir)
	if path != "" {
		path = path + "/"
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path+"db/migrations", dbName, driver)
	if err != nil {
		return fmt.Errorf("%s checkMigrations(): %v", logLabel, err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s checkMigrations(): %v", logLabel, err)
	}

	return nil
}
