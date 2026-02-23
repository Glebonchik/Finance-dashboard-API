package dto

// Запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Запрос на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Ответ с токенами
type AuthResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

// Данные пользователя в ответе
type UserDTO struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	GlobalCurrency string `json:"global_currency"`
}

// Запрос на обновление токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
