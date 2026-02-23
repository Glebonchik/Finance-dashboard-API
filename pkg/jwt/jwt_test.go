package jwt

import (
	"testing"
	"time"
)

func TestManager_GenerateAndValidateAccessToken(t *testing.T) {
	manager := NewManager("test-secret-key", 15*time.Minute, 24*time.Hour)

	userID := "test-user-id"
	email := "test@example.com"

	token, err := manager.GenerateAccessToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := manager.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestManager_GenerateAndValidateRefreshToken(t *testing.T) {
	manager := NewManager("test-secret-key", 15*time.Minute, 24*time.Hour)

	userID := "test-user-id"

	token, err := manager.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("Failed to generate refresh token: %v", err)
	}

	validatedUserID, err := manager.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, validatedUserID)
	}
}

func TestManager_InvalidToken(t *testing.T) {
	manager := NewManager("test-secret-key", 15*time.Minute, 24*time.Hour)

	_, err := manager.ValidateAccessToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

func TestManager_ExpiredToken(t *testing.T) {
	// Создаём менеджер с очень коротким временем жизни токена
	manager := NewManager("test-secret-key", 1*time.Second, 24*time.Hour)

	token, _ := manager.GenerateAccessToken("user-id", "test@example.com")

	// Ждём истечения токена
	time.Sleep(2 * time.Second)

	_, err := manager.ValidateAccessToken(token)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}
