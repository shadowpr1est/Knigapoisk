package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "github.com/shadowpr1est/knigapoisk-auth-service/api/proto"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/usecase/auth"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	authUseCase auth.UseCase
	logger      *zap.Logger
}

func NewAuthServer(authUseCase auth.UseCase, logger *zap.Logger) *AuthServer {
	return &AuthServer{
		authUseCase: authUseCase,
		logger:      logger,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	output, err := s.authUseCase.Register(ctx, auth.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case auth.ErrUserAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			s.logger.Error("register error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &authpb.RegisterResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	output, err := s.authUseCase.Login(ctx, auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			s.logger.Error("login error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &authpb.LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (s *AuthServer) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	output, err := s.authUseCase.Refresh(ctx, auth.RefreshInput{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		switch err {
		case auth.ErrInvalidToken:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case auth.ErrTokenExpired:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			s.logger.Error("refresh error", zap.Error(err))
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &authpb.RefreshResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if err := s.authUseCase.Logout(ctx, req.RefreshToken); err != nil {
		s.logger.Error("logout error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authpb.LogoutResponse{Success: true}, nil
}

func (s *AuthServer) Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error) {
	out, err := s.authUseCase.Validate(ctx, req.AccessToken)
	if err != nil {
		s.logger.Error("validate error", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &authpb.ValidateResponse{
		Valid:  out.Valid,
		UserId: out.UserID,
		Role:   out.Role,
	}, nil
}

func (s *AuthServer) Health(ctx context.Context, _ *authpb.HealthRequest) (*authpb.HealthResponse, error) {
	return &authpb.HealthResponse{Status: "ok"}, nil
}

