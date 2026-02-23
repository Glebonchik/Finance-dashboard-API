package model

import "time"

// User представляет пользователя системы
type User struct {
	ID            string
	Email         string
	PasswordHash  *string
	GoogleID      *string
	GlobalCurrency string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Currency представляет код валюты (ISO 4217)
type Currency string

const (
	CurrencyRUB Currency = "RUB"
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

// IsValid проверяет валидность валюты
func (c Currency) IsValid() bool {
	switch c {
	case CurrencyRUB, CurrencyUSD, CurrencyEUR:
		return true
	}
	return false
}
