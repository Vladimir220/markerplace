package main

import (
	"auth_service/crypto"
	"auth_service/gen"
	"auth_service/models"
	"auth_service/network/auth"
	"context"
	"fmt"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateServer(ctx context.Context, auth auth.IAuthentication, tokenManager crypto.ITokenManager) (server gen.AuthServer) {
	server = Server{
		logger:       logger_lib.CreateLoggerGateway(ctx, "Server"),
		auth:         auth,
		tokenManager: tokenManager,
	}
	return
}

type Server struct {
	gen.UnimplementedAuthServer
	logger       logger_lib.ILogger
	auth         auth.IAuthentication
	tokenManager crypto.ITokenManager
}

func (s Server) ValidateToken(ctx context.Context, req *gen.ValidateTokenRequest) (resp *gen.ValidateTokenResponse, err error) {
	logLabel := fmt.Sprintf("%s:%s:", serviceName, "ValidateToken()")

	if req == nil {
		err = fmt.Errorf("%s %s", logLabel, "ValidateTokenRequest is nil")
		s.logger.WriteError(err.Error())
		return
	}

	user, isValid, isErr := s.tokenManager.ValidateToken(req.Token)
	if isErr {
		err = fmt.Errorf("%s %s", logLabel, "Error")
		return
	}

	resp = &gen.ValidateTokenResponse{
		Login:   user.Login,
		Group:   user.Group,
		IsValid: isValid,
	}

	return
}

func (s Server) Register(ctx context.Context, req *gen.RegisterRequest) (resp *gen.RegisterResponse, err error) {
	logLabel := fmt.Sprintf("%s:%s:", serviceName, "Register()")

	if req == nil {
		err = fmt.Errorf("%s %s", logLabel, "RegisterRequest is nil")
		s.logger.WriteError(err.Error())
		return
	}

	token, err := s.auth.Register(req.Login, req.Password)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	resp = &gen.RegisterResponse{
		Token: token,
	}

	return
}

func (s Server) Login(ctx context.Context, req *gen.LoginRequest) (resp *gen.LoginResponse, err error) {
	logLabel := fmt.Sprintf("%s:%s:", serviceName, "Login()")

	if req == nil {
		err = fmt.Errorf("%s %s", logLabel, "LoginRequest is nil")
		s.logger.WriteError(err.Error())
		return
	}

	token, err := s.auth.Login(req.Login, req.Password)
	if err != nil {
		err = fmt.Errorf("%s %s", logLabel, err)
		return
	}

	resp = &gen.LoginResponse{
		Token: token,
	}

	return

}

func (s Server) GenerateToken(ctx context.Context, req *gen.GenerateTokenRequest) (resp *gen.GenerateTokenResponse, err error) {
	logLabel := fmt.Sprintf("%s:%s:", serviceName, "GenerateToken()")

	if req == nil {
		err = fmt.Errorf("%s %s", logLabel, "GenerateTokenRequest is nil")
		s.logger.WriteError(err.Error())
		return
	}

	token, isErr := s.tokenManager.GenerateToken(models.User{
		Login: req.GetLogin(),
		Group: req.GetGroup(),
	})
	if isErr {
		err = fmt.Errorf("%s %s", logLabel, "Error")
		return
	}

	resp = &gen.GenerateTokenResponse{
		Token: token,
	}

	return

}
