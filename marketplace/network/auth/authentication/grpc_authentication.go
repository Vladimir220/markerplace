package authentication

import (
	"context"
	"fmt"
	"marketplace/env"
	"marketplace/gen/auth_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateGrpcAuthentication(ctx context.Context) (auth IAuthentication, err error) {
	data, err := env.GetServicesData()
	if err != nil {
		err = fmt.Errorf("CreateGrpcAuthentication(): %v", err)
		return
	}

	auth = &GrpcAuthentication{
		host: data.AuthGRPCHost,
		ctx:  ctx,
	}

	return
}

type GrpcAuthentication struct {
	host string
	ctx  context.Context
}

func (auth GrpcAuthentication) Register(login, password string) (token string, err error) {
	logLabel := "GrpcAuthentication:Register():"

	conn, err := grpc.NewClient(auth.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	defer conn.Close()

	client := auth_service.NewAuthClient(conn)

	params := &auth_service.RegisterRequest{
		Login:    login,
		Password: password,
	}

	resp, err := client.Register(auth.ctx, params)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	token = resp.Token
	return
}
func (auth GrpcAuthentication) Login(login, password string) (token string, err error) {
	logLabel := "GrpcAuthentication:Login():"

	conn, err := grpc.NewClient(auth.host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	defer conn.Close()

	client := auth_service.NewAuthClient(conn)

	params := &auth_service.LoginRequest{
		Login:    login,
		Password: password,
	}

	resp, err := client.Login(auth.ctx, params)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	token = resp.Token
	return
}
