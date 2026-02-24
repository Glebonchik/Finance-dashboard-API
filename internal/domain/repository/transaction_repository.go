package repository

import (
	"context"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
)

// TransactionRepository определяет интерфейс для работы с транзакциями
type TransactionRepository interface {
	// Create создаёт новую транзакцию
	Create(ctx context.Context, tx *model.Transaction) error

	// GetByID находит транзакцию по ID
	GetByID(ctx context.Context, id string) (*model.Transaction, error)

	// GetByUserID находит транзакции пользователя
	GetByUserID(ctx context.Context, filter model.TransactionFilter) ([]*model.Transaction, error)

	// Update обновляет транзакцию
	Update(ctx context.Context, tx *model.Transaction) error

	// Delete удаляет транзакцию по ID
	Delete(ctx context.Context, id string) error

	// GetTotalCount возвращает общее количество транзакций пользователя
	GetTotalCount(ctx context.Context, userID string) (int64, error)
}

// CategoryRepository определяет интерфейс для работы с категориями
type CategoryRepository interface {
	// GetAll возвращает все категории
	GetAll(ctx context.Context) ([]*model.Category, error)

	// GetByID находит категорию по ID
	GetByID(ctx context.Context, id int) (*model.Category, error)

	// GetDefault возвращает системные категории
	GetDefault(ctx context.Context) ([]*model.Category, error)
}

// UserCategoryRuleRepository определяет интерфейс для работы с правилами категоризации
type UserCategoryRuleRepository interface {
	// Create создаёт новое правило
	Create(ctx context.Context, rule *model.UserCategoryRule) error

	// GetByUserID возвращает правила пользователя
	GetByUserID(ctx context.Context, userID string) ([]*model.UserCategoryRule, error)

	// GetByKeyword находит правило по ключевому слову
	GetByKeyword(ctx context.Context, userID, keyword string) (*model.UserCategoryRule, error)

	// Delete удаляет правило по ID
	Delete(ctx context.Context, id string) error

	// Update обновляет правило
	Update(ctx context.Context, rule *model.UserCategoryRule) error
}
