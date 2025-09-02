package main

import (
	"context"
	"net"
	"reader_db_service/DAO/postgres"
	"reader_db_service/env"
	"reader_db_service/gen"
	"reader_db_service/log"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.CreateLoggerAdapter(ctx, "main()")

	host, _, err := env.GetServiceData()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	dao, err := postgres.CreateMarketplaceDAO()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	server := CreateServer(dao)

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
