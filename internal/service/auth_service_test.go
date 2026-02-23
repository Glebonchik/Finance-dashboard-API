// - Регистрация нового пользователя 
// - Обработка дубликатов при регистрации 
// - Вход по email и паролю 
// - Обработка неверных учётных данных 
// - Генерация и валидация JWT токенов

package service

import (
	"context"
	"testing"
	"time"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/internal/repository"
)

type mockUserRepository struct {
	users      map[string]*model.User
	emailIndex map[string]*model.User
}

var errUserNotFound = repository.ErrUserNotFound

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:      make(map[string]*model.User),
		emailIndex: make(map[string]*model.User),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, errUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, ok := m.emailIndex[email]
	if !ok {
		return nil, errUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	for _, user := range m.users {
		if user.GoogleID != nil && *user.GoogleID == googleID {
			return user, nil
		}
	}
	return nil, errUserNotFound
}

func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error {
	m.users[user.ID] = user
	m.emailIndex[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

func TestAuthService_Register(t *testing.T) {
	repo := newMockUserRepository()
	authService := NewAuthService(repo, AuthServiceConfig{
		JWTSecret:     "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	})

	user, err := authService.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}

	if user.PasswordHash == nil {
		t.Error("PasswordHash should not be nil")
	}
}

func TestAuthService_Register_Duplicate(t *testing.T) {
	repo := newMockUserRepository()
	authService := NewAuthService(repo, AuthServiceConfig{
		JWTSecret:     "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	})

	_, err := authService.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register first user: %v", err)
	}

	_, err = authService.Register(context.Background(), "test@example.com", "password456")
	if err != ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestAuthService_Login(t *testing.T) {
	repo := newMockUserRepository()
	authService := NewAuthService(repo, AuthServiceConfig{
		JWTSecret:     "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	})

	_, err := authService.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	user, err := authService.Login(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	repo := newMockUserRepository()
	authService := NewAuthService(repo, AuthServiceConfig{
		JWTSecret:     "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	})

	_, err := authService.Login(context.Background(), "nonexistent@example.com", "password123")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_GenerateTokens(t *testing.T) {
	repo := newMockUserRepository()
	authService := NewAuthService(repo, AuthServiceConfig{
		JWTSecret:     "test-secret",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	})

	user := &model.User{
		ID:    "test-user-id",
		Email: "test@example.com",
	}

	accessToken, refreshToken, err := authService.GenerateTokens(user)
	if err != nil {
		t.Fatalf("Failed to generate tokens: %v", err)
	}

	if accessToken == "" {
		t.Error("AccessToken should not be empty")
	}

	if refreshToken == "" {
		t.Error("RefreshToken should not be empty")
	}

	claims, err := authService.ValidateAccessToken(accessToken)
	if err != nil {
		t.Errorf("Failed to validate access token: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID %s, got %s", user.ID, claims.UserID)
	}

	validatedUserID, err := authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("Failed to validate refresh token: %v", err)
	}

	if validatedUserID != user.ID {
		t.Errorf("Expected UserID %s, got %s", user.ID, validatedUserID)
	}
}
