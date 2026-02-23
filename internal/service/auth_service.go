package service

import (
	"context"
	"errors"
	"time"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/internal/domain/repository"
	domainService "github.com/gibbon/finace-dashboard/internal/domain/service"
	repo "github.com/gibbon/finace-dashboard/internal/repository"
	"github.com/gibbon/finace-dashboard/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var ErrUserNotFound = repo.ErrUserNotFound

// Конфигурация для сервиса
type AuthServiceConfig struct {
	JWTSecret     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type authServiceImpl struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo repository.UserRepository, cfg AuthServiceConfig) domainService.AuthService {
	return &authServiceImpl{
		userRepo:   userRepo,
		jwtManager: jwt.NewManager(cfg.JWTSecret, cfg.AccessExpiry, cfg.RefreshExpiry),
	}
}

func (s *authServiceImpl) Register(ctx context.Context, email, password string) (*model.User, error) {
	// Проверяем, существует ли пользователь
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

	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *authServiceImpl) LoginWithGoogle(ctx context.Context, googleID, email string) (*model.User, error) {
	// Поиск пользователя по Google ID
	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err == nil && user != nil {
		return user, nil
	}

	// Поиск по email (если не нашли по Google ID)
	user, err = s.userRepo.GetByEmail(ctx, email)
	if err == nil && user != nil {
		// Приязка Google ID к существующему пользователю
		user.GoogleID = &googleID
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

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

// Генерация пары аксесс и рефреш токенов
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

func (s *authServiceImpl) ValidateAccessToken(token string) (*jwt.Claims, error) {
	return s.jwtManager.ValidateAccessToken(token)
}

func (s *authServiceImpl) ValidateRefreshToken(token string) (string, error) {
	return s.jwtManager.ValidateRefreshToken(token)
}
