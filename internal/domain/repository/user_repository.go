package repository

import (
	"context"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
)

type UserRepository interface {
	// Создаёт нового пользователя
	Create(ctx context.Context, user *model.User) error

	// Находит пользователя по ID
	GetByID(ctx context.Context, id string) (*model.User, error)

	// Находит пользователя по email
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// Находит пользователя по Google ID
	GetByGoogleID(ctx context.Context, googleID string) (*model.User, error)

	// Обновляет данные пользователя
	Update(ctx context.Context, user *model.User) error

	// Удаляет пользователя по ID
	Delete(ctx context.Context, id string) error
}
