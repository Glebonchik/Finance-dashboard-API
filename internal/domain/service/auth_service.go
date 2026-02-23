package service

import (
	"context"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/pkg/jwt"
)

type AuthService interface {
	// Регистрирует нового пользователя
	Register(ctx context.Context, email, password string) (*model.User, error)

	// Ввыполняет вход пользователя
	Login(ctx context.Context, email, password string) (*model.User, error)

	// Выполняет вход через Google
	LoginWithGoogle(ctx context.Context, googleID, email string) (*model.User, error)

	// Генерирует пару токенов для пользователя
	GenerateTokens(user *model.User) (accessToken, refreshToken string, err error)

	// Валидирует access токен
	ValidateAccessToken(token string) (*jwt.Claims, error)

	// Валидирует refresh токен
	ValidateRefreshToken(token string) (string, error)
}
