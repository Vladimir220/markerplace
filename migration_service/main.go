package main

import (
	"context"

	"github.com/Vladimir220/markerplace/logger_lib"
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

	logger := logger_lib.CreateLoggerGateway(ctx, "main()")

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
