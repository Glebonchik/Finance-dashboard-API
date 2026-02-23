package service

import (
	"context"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/pkg/jwt"
)

// AuthService определяет интерфейс для аутентификации
type AuthService interface {
	// Register регистрирует нового пользователя
	Register(ctx context.Context, email, password string) (*model.User, error)
	
	// Login выполняет вход пользователя
	Login(ctx context.Context, email, password string) (*model.User, error)
	
	// LoginWithGoogle выполняет вход через Google
	LoginWithGoogle(ctx context.Context, googleID, email string) (*model.User, error)
	
	// GenerateTokens генерирует пару токенов для пользователя
	GenerateTokens(user *model.User) (accessToken, refreshToken string, err error)
	
	// ValidateAccessToken валидирует access токен
	ValidateAccessToken(token string) (*jwt.Claims, error)
	
	// ValidateRefreshToken валидирует refresh токен
	ValidateRefreshToken(token string) (string, error)
}
