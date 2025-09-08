package main

import (
	"auth_service/crypto"
	"auth_service/db/DAO/postgres"
	"auth_service/db/DAO/redis"
	"auth_service/env"
	"auth_service/gen"
	"auth_service/log/proxies"
	"auth_service/network"
	"auth_service/network/auth"
	"context"
	"fmt"
	"net"
	"os"

	"github.com/Vladimir220/markerplace/logger_lib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var serviceName string

func main() {
	logLabel := "main():"

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

	go network.HealthListener()

	tokensDao, err := redis.CreateTokensDAO()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		logger.WriteError(err.Error())
		panic(err)
	}
	tokensDaoWithLog := proxies.CreateTokensDAOWithLog(ctx, tokensDao)
	defer tokensDaoWithLog.Close()

	tokenManager := crypto.CreateTokenManager(ctx, tokensDaoWithLog)

	dao, err := postgres.CreateMarketplaceDAO()
	if err != nil {
		err = fmt.Errorf("%s%v", logLabel, err)
		logger.WriteError(err.Error())
		panic(err)
	}

	daoWithLogs := proxies.CreateDAOWithLog(ctx, dao)
	defer daoWithLogs.Close()

	auth := auth.CreateAuthentication(ctx, tokenManager, daoWithLogs)

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
