package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/shadowpr1est/knigapoisk-auth-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/domain/repository"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/jwt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
)

type AuthUseCase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	jwt       jwt.Manager
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	jwtManager jwt.Manager,
) UseCase {
	return &AuthUseCase{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwt:       jwtManager,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	_, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return RegisterOutput{}, ErrUserAlreadyExists
	}

	_, err = u.userRepo.GetByUsername(ctx, input.Username)
	if err == nil {
		return RegisterOutput{}, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterOutput{}, err
	}
	user := &entity.User{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: string(passwordHash),
		Role:         entity.RoleUser,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return RegisterOutput{}, err
	}
	return u.generateTokenPair(ctx, user)
}

func (u *AuthUseCase) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}

	output, err := u.generateTokenPair(ctx, user)
	if err != nil {
		return LoginOutput{}, err
	}

	return LoginOutput{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (u *AuthUseCase) Refresh(ctx context.Context, input RefreshInput) (RefreshOutput, error) {
	tokenHash := hashToken(input.RefreshToken)
	savedToken, err := u.tokenRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		return RefreshOutput{}, ErrInvalidToken
	}

	if time.Now().After(savedToken.ExpiresAt) {
		return RefreshOutput{}, ErrTokenExpired
	}

	if savedToken.Revoked {
		return RefreshOutput{}, ErrInvalidToken
	}

	user, err := u.userRepo.GetByID(ctx, savedToken.UserID)
	if err != nil {
		return RefreshOutput{}, err
	}

	if err := u.tokenRepo.Revoke(ctx, tokenHash); err != nil {
		return RefreshOutput{}, err
	}

	output, err := u.generateTokenPair(ctx, user)
	if err != nil {
		return RefreshOutput{}, err
	}

	return RefreshOutput{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (u *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return u.tokenRepo.Revoke(ctx, tokenHash)
}

func (u *AuthUseCase) LogoutAll(ctx context.Context, userID int64) error {
	return u.tokenRepo.RevokeAllByUserID(ctx, userID)
}

func (u *AuthUseCase) Validate(ctx context.Context, accessToken string) (ValidateOutput, error) {
	claims, err := u.jwt.Validate(accessToken)
	if err != nil {
		return ValidateOutput{Valid: false}, nil
	}
	return ValidateOutput{
		Valid:  true,
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}

func (u *AuthUseCase) generateTokenPair(ctx context.Context, user *entity.User) (RegisterOutput, error) {
	accessToken, err := u.jwt.Generate(user.ID, string(user.Role))
	if err != nil {
		return RegisterOutput{}, err
	}

	refreshToken, err := generateRandomToken()
	if err != nil {
		return RegisterOutput{}, err
	}
	token := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: hashToken(refreshToken),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	if err := u.tokenRepo.Create(ctx, token); err != nil {
		return RegisterOutput{}, err
	}
	return RegisterOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

