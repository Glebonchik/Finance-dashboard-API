package dto

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse представляет ответ с токенами
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

// UserDTO представляет данные пользователя в ответе
type UserDTO struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	GlobalCurrency string `json:"global_currency"`
}

// RefreshRequest представляет запрос на обновление токена
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
