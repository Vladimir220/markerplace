package main

import (
	"context"
	"fmt"
	"net"
	"reader_db_service/DAO/postgres"
	"reader_db_service/env"
	"reader_db_service/gen"

	"github.com/Vladimir220/markerplace/logger_lib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger_lib.CreateLoggerAdapter(ctx, "main()")

	host, _, err := env.GetServiceData()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	go HealthListener()

	dao, err := postgres.CreateMarketplaceDAO()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	server := CreateServer(dao)
	fmt.Println("я умею писать")

	lis, err := net.Listen("tcp", host)
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	var opts []grpc.ServerOption
	serverGRPC := grpc.NewServer(opts...)

	gen.RegisterReaderServer(serverGRPC, server)

	err = serverGRPC.Serve(lis)
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}
}
