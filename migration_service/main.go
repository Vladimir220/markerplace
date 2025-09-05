package main

import (
	"context"
	"migration_service/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.CreateLoggerAdapter(ctx, "main()")

	connection, err := Connect()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}
	defer connection.Close()

	err = CheckDbExistence()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	err = CheckMigrations(connection)
	if err != nil && err != migrate.ErrNoChange {
		logger.WriteError(err.Error())
		panic(err)
	}
}
