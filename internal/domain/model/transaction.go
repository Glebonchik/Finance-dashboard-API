package model

import "time"

// Transaction представляет финансовую транзакцию пользователя
type Transaction struct {
	ID          string
	UserID      string
	Amount      float64
	Currency    string
	Description string
	Date        time.Time
	PlaceName   *string
	PlaceLat    *float64
	PlaceLon    *float64
	CategoryID  *int
	IsConfirmed bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TransactionFilter параметры для поиска транзакций
type TransactionFilter struct {
	UserID     string
	CategoryID *int
	FromDate   *time.Time
	ToDate     *time.Time
	Limit      int
	Offset     int
}

// Category представляет категорию транзакции
type Category struct {
	ID        int
	Name      string
	IsDefault bool
	CreatedAt time.Time
}

// UserCategoryRule представляет правило категоризации пользователя
type UserCategoryRule struct {
	ID         string
	UserID     string
	Keyword    string
	CategoryID int
	CreatedAt  time.Time
}
