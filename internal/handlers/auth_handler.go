package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gibbon/finace-dashboard/internal/dto"
	domainService "github.com/gibbon/finace-dashboard/internal/domain/service"
	"github.com/gibbon/finace-dashboard/internal/service"
)

// AuthHandler обрабатывает HTTP запросы для аутентификации
type AuthHandler struct {
	authService domainService.AuthService
}

// NewAuthHandler создаёт новый AuthHandler
func NewAuthHandler(authService domainService.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register
// @Summary Регистрация нового пользователя
// @Description Создаёт новый аккаунт с email и паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 409 {object} map[string]string "Пользователь уже существует"
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "email and password are required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, `{"error": "password must be at least 8 characters"}`, http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			http.Error(w, `{"error": "user already exists"}`, http.StatusConflict)
			return
		}
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(user)
	if err != nil {
		http.Error(w, `{"error": "failed to generate tokens"}`, http.StatusInternalServerError)
		return
	}

	response := dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:             user.ID,
			Email:          user.Email,
			GlobalCurrency: user.GlobalCurrency,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login
// @Summary Вход в систему
// @Description Аутентификация пользователя по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 401 {object} map[string]string "Неверные учётные данные"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, `{"error": "invalid credentials"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := h.authService.GenerateTokens(user)
	if err != nil {
		http.Error(w, `{"error": "failed to generate tokens"}`, http.StatusInternalServerError)
		return
	}

	response := dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:             user.ID,
			Email:          user.Email,
			GlobalCurrency: user.GlobalCurrency,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Refresh
// @Summary Обновление токена
// @Description Получение новой пары токенов по refresh токену
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh токен"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 401 {object} map[string]string "Невалидный токен"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	userID, err := h.authService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, `{"error": "invalid refresh token"}`, http.StatusUnauthorized)
		return
	}

	// Получаем данные пользователя
	// Для этого нам понадобится метод в сервисе
	// Пока заглушка - в реальной реализации нужно получить пользователя из БД
	_ = userID

	// TODO: Получить пользователя и сгенерировать новые токены
	http.Error(w, `{"error": "not implemented"}`, http.StatusNotImplemented)
}

// Logout
// @Summary Выход из системы
// @Description Инвалидация токенов (опционально)
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "Успешный выход"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Реализовать blacklist для токенов в Redis
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out successfully"})
}
