package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/shadowpr1est/knigapoisk-auth-service/api/proto"
)

type AuthClient interface {
	Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error)
	Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error)
	Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error)
	Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error)
	Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error)
}

type authClient struct {
	cc authpb.AuthServiceClient
}

func NewAuthClient(addr string) (AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &authClient{cc: authpb.NewAuthServiceClient(conn)}, nil
}

func (c *authClient) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	return c.cc.Register(ctx, req)
}

func (c *authClient) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return c.cc.Login(ctx, req)
}

func (c *authClient) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	return c.cc.Refresh(ctx, req)
}

func (c *authClient) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	return c.cc.Logout(ctx, req)
}

func (c *authClient) Validate(ctx context.Context, req *authpb.ValidateRequest) (*authpb.ValidateResponse, error) {
	return c.cc.Validate(ctx, req)
}

