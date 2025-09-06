package main

import (
	"auth_service/crypto"
	"auth_service/db/DAO/redis"
	"auth_service/env"
	"auth_service/gen"
	"auth_service/network"
	"auth_service/network/auth"
	"context"
	"fmt"
	"net"

	"github.com/Vladimir220/markerplace/logger_lib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	logLabel := "main():"

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	logsConfig := logger_lib.LoggerGatewayConfig{
		PrintErrorsToStdOut:   true,
		PrintWarningsToStdOut: true,
		PrintInfoToStdOut:     true,
	}

	logger := logger_lib.CreateLoggerGateway(ctx, "main()")

	host, _, err := env.GetServiceData()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	go network.HealthListener()

	tokensDao, err := redis.CreateTokensDAO()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		logger.WriteError(err.Error())
		panic(err)
	}

	tokenManager := crypto.CreateTokenManager(ctx, tokensDao, false)
	auth, err := auth.CreateAuthentication(ctx, tokenManager, true)
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		logger.WriteError(err.Error())
		panic(err)
	}

	server := CreateServer(ctx, auth, tokenManager)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	var opts []grpc.ServerOption
	serverGRPC := grpc.NewServer(opts...)

	gen.RegisterAuthServer(serverGRPC, server)

	err = serverGRPC.Serve(lis)
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}
}
