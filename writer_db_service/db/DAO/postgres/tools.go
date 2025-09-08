package postgres

import (
	"database/sql"
	"fmt"
	"time"
	"writer_db_service/env"

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

func checkDbExistence(times uint, waitingTime time.Duration) (err error) {
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
	defer tempConn.Close()

	query := "SELECT EXISTS (SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)"

	var isDbExist bool
	for range times {
		err = tempConn.QueryRow(query, data.DbName).Scan(&isDbExist)
		if err != nil {
			err = fmt.Errorf("%s %v", logLabel, err)
			return
		}
		if isDbExist {
			return
		}
		time.Sleep(waitingTime)
	}

	err = fmt.Errorf("%s %s", logLabel, "db doesn't exist")

	return
}
