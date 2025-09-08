package main

import (
	"context"
	"net"
	"os"
	"reader_db_service/db/DAO/postgres"
	"reader_db_service/env"
	"reader_db_service/gen"
	"reader_db_service/log/proxies"

	"github.com/Vladimir220/markerplace/logger_lib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var serviceName string

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	serviceName = os.Getenv("SERVICE_NAME")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := logger_lib.CreateLoggerGateway(ctx, "main()")

	host, _, err := env.GetServiceData()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	go HealthListener()

	dao, err := postgres.CreateReaderMarketplaceDAO()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}
	daoWithLogs := proxies.CreateDAOWithLog(ctx, dao)
	defer daoWithLogs.Close()

	server := CreateServer(ctx, daoWithLogs)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}
	defer lis.Close()

	var opts []grpc.ServerOption
	serverGRPC := grpc.NewServer(opts...)

	gen.RegisterReaderServer(serverGRPC, server)

	err = serverGRPC.Serve(lis)
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}
}
