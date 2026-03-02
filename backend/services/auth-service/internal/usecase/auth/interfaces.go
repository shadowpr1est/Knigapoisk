package auth

import "context"

type UseCase interface {
	Register(ctx context.Context, input RegisterInput) (RegisterOutput, error)
	Login(ctx context.Context, input LoginInput) (LoginOutput, error)
	Refresh(ctx context.Context, input RefreshInput) (RefreshOutput, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID int64) error
	Validate(ctx context.Context, accessToken string) (ValidateOutput, error)
}

