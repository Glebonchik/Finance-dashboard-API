package dto

import "time"

// Запрос на создание транзакции
type CreateTransactionRequest struct {
	Amount      float64  `json:"amount"`
	Currency    string   `json:"currency"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	PlaceName   *string  `json:"place_name,omitempty"`
	PlaceLat    *float64 `json:"place_lat,omitempty"`
	PlaceLon    *float64 `json:"place_lon,omitempty"`
}

// Запрос на обновление транзакции
type UpdateTransactionRequest struct {
	Amount      float64  `json:"amount"`
	Currency    string   `json:"currency"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	PlaceName   *string  `json:"place_name,omitempty"`
	PlaceLat    *float64 `json:"place_lat,omitempty"`
	PlaceLon    *float64 `json:"place_lon,omitempty"`
	CategoryID  *int     `json:"category_id,omitempty"`
	IsConfirmed bool     `json:"is_confirmed"`
}

// Jтвет с данными транзакции
type TransactionResponse struct {
	ID          string     `json:"id"`
	Amount      float64    `json:"amount"`
	Currency    string     `json:"currency"`
	Description string     `json:"description"`
	Date        time.Time  `json:"date"`
	PlaceName   *string    `json:"place_name,omitempty"`
	PlaceLat    *float64   `json:"place_lat,omitempty"`
	PlaceLon    *float64   `json:"place_lon,omitempty"`
	CategoryID  *int       `json:"category_id,omitempty"`
	Category    *string    `json:"category,omitempty"`
	IsConfirmed bool       `json:"is_confirmed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Ответ с данными категории
type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Запрос на создание правила категоризации
type CreateCategoryRuleRequest struct {
	Keyword    string `json:"keyword"`
	CategoryID int    `json:"category_id"`
}

// Ответ с данными правила
type CategoryRuleResponse struct {
	ID         string `json:"id"`
	Keyword    string `json:"keyword"`
	CategoryID int    `json:"category_id"`
	Category   string `json:"category"`
}

// Список транзакций с пагинацией
type TransactionsListResponse struct {
	Transactions []*TransactionResponse `json:"transactions"`
	Total        int64                  `json:"total"`
	Limit        int                    `json:"limit"`
	Offset       int                    `json:"offset"`
}
