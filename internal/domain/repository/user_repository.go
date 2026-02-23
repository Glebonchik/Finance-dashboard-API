package repository

import (
	"context"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	// Create создаёт нового пользователя
	Create(ctx context.Context, user *model.User) error
	
	// GetByID находит пользователя по ID
	GetByID(ctx context.Context, id string) (*model.User, error)
	
	// GetByEmail находит пользователя по email
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	
	// GetByGoogleID находит пользователя по Google ID
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)
	
	// Update обновляет данные пользователя
	Update(ctx context.Context, user *model.User) error
	
	// Delete удаляет пользователя по ID
	Delete(ctx context.Context, id string) error
}
