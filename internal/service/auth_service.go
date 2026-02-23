package service

import (
	"context"
	"errors"
	"time"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	domainService "github.com/gibbon/finace-dashboard/internal/domain/service"
	"github.com/gibbon/finace-dashboard/internal/domain/repository"
	repo "github.com/gibbon/finace-dashboard/internal/repository"
	"github.com/gibbon/finace-dashboard/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Ошибки репозитория
var ErrUserNotFound = repo.ErrUserNotFound

// AuthServiceConfig содержит конфигурацию для сервиса
type AuthServiceConfig struct {
	JWTSecret        string
	AccessExpiry     time.Duration
	RefreshExpiry    time.Duration
}

// authServiceImpl реализует AuthService
type authServiceImpl struct {
	userRepo  repository.UserRepository
	jwtManager *jwt.Manager
}

// NewAuthService создаёт новый экземпляр сервиса аутентификации
func NewAuthService(userRepo repository.UserRepository, cfg AuthServiceConfig) domainService.AuthService {
	return &authServiceImpl{
		userRepo:   userRepo,
		jwtManager: jwt.NewManager(cfg.JWTSecret, cfg.AccessExpiry, cfg.RefreshExpiry),
	}
}

func (s *authServiceImpl) Register(ctx context.Context, email, password string) (*model.User, error) {
	// Проверяем, существует ли уже пользователь
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	passwordStr := string(hashedPassword)

	// Создаём нового пользователя
	user := &model.User{
		ID:             uuid.New().String(),
		Email:          email,
		PasswordHash:   &passwordStr,
		GlobalCurrency: string(model.CurrencyRUB),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authServiceImpl) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Проверяем пароль
	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *authServiceImpl) LoginWithGoogle(ctx context.Context, googleID, email string) (*model.User, error) {
	// Пытаемся найти существующего пользователя по Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err == nil && user != nil {
		return user, nil
	}

	// Если не найден, пытаемся найти по email
	user, err = s.userRepo.GetByEmail(ctx, email)
	if err == nil && user != nil {
		// Привязываем Google ID к существующему аккаунту
		user.GoogleID = &googleID
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// Создаём нового пользователя
	user = &model.User{
		ID:             uuid.New().String(),
		Email:          email,
		GoogleID:       &googleID,
		GlobalCurrency: string(model.CurrencyRUB),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateTokens генерирует пару токенов для пользователя
func (s *authServiceImpl) GenerateTokens(user *model.User) (accessToken, refreshToken string, err error) {
	accessToken, err = s.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken валидирует access токен
func (s *authServiceImpl) ValidateAccessToken(token string) (*jwt.Claims, error) {
	return s.jwtManager.ValidateAccessToken(token)
}

// ValidateRefreshToken валидирует refresh токен
func (s *authServiceImpl) ValidateRefreshToken(token string) (string, error) {
	return s.jwtManager.ValidateRefreshToken(token)
}
