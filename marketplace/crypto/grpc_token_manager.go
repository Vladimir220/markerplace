package crypto

import (
	"context"
	"fmt"
	"marketplace/env"
	"marketplace/gen/auth_service"
	"marketplace/models"

	"github.com/Vladimir220/markerplace/logger_lib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
type ITokenManager interface {
	GenerateToken(user models.User) (token string, isErr bool)
	ValidateToken(token string) (user models.User, isValid, isErr bool)
}
*/

func CreateGrpcTokenManager(ctx context.Context) (tokenManager ITokenManager, err error) {
	data, err := env.GetServicesData()
	if err != nil {
		err = fmt.Errorf("CreateGrpcTokenManager(): %v", err)
		return
	}

	tokenManager = &GrpcTokenManager{
		host:   data.AuthGRPCHost,
		logger: logger_lib.CreateLoggerGateway(ctx, "GrpcTokenManager"),
		ctx:    ctx,
	}

	return
}

type GrpcTokenManager struct {
	host   string
	ctx    context.Context
	logger logger_lib.ILogger
}

func (tm GrpcTokenManager) ValidateToken(token string) (user models.User, isValid, isErr bool) {
	logLabel := "ValidateToken():"

	conn, err := grpc.NewClient(tm.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		tm.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		isErr = true
		return
	}
	defer conn.Close()

	client := auth_service.NewAuthClient(conn)

	params := &auth_service.ValidateTokenRequest{
		Token: token,
	}

	resp, err := client.ValidateToken(tm.ctx, params)
	if err != nil {
		tm.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		isErr = true
		return
	}

	user.Login = resp.GetLogin()
	user.Group = resp.GetGroup()
	isValid = resp.GetIsValid()
	return
}

func (tm GrpcTokenManager) GenerateToken(user models.User) (token string, isErr bool) {
	logLabel := "GenerateToken():"

	conn, err := grpc.NewClient(tm.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		tm.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		isErr = true
		return
	}
	defer conn.Close()

	client := auth_service.NewAuthClient(conn)

	params := &auth_service.GenerateTokenRequest{
		Login: user.Login,
		Group: user.Group,
	}

	resp, err := client.GenerateToken(tm.ctx, params)
	if err != nil {
		tm.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		isErr = true
		return
	}

	token = resp.GetToken()
	return
}
