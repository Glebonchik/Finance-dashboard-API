package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/gibbon/finace-dashboard/internal/dto"
	"github.com/gibbon/finace-dashboard/internal/middleware"
	"github.com/gibbon/finace-dashboard/internal/service"
)

// CategoryHandler обрабатывает HTTP запросы для категорий
type CategoryHandler struct {
	txService service.TransactionService
}

// NewCategoryHandler создаёт новый CategoryHandler
func NewCategoryHandler(txService service.TransactionService) *CategoryHandler {
	return &CategoryHandler{
		txService: txService,
	}
}

// GetAll
// @Summary Получить все категории
// @Description Получение списка всех категорий
// @Tags categories
// @Produce json
// @Success 200 {array} dto.CategoryResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.txService.GetCategories(r.Context())
	if err != nil {
		http.Error(w, `{"error": "failed to get categories"}`, http.StatusInternalServerError)
		return
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = dto.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CategoryRuleHandler обрабатывает HTTP запросы для правил категоризации
type CategoryRuleHandler struct {
	txService service.TransactionService
}

// NewCategoryRuleHandler создаёт новый CategoryRuleHandler
func NewCategoryRuleHandler(txService service.TransactionService) *CategoryRuleHandler {
	return &CategoryRuleHandler{
		txService: txService,
	}
}

// Create
// @Summary Создать правило категоризации
// @Description Создание нового правила для автоматической категоризации
// @Tags category-rules
// @Accept json
// @Produce json
// @Param request body dto.CreateCategoryRuleRequest true "Данные правила"
// @Success 201 {object} dto.CategoryRuleResponse
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Router /api/v1/category-rules [post]
func (h *CategoryRuleHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req dto.CreateCategoryRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Keyword == "" {
		http.Error(w, `{"error": "keyword is required"}`, http.StatusBadRequest)
		return
	}

	if req.CategoryID <= 0 {
		http.Error(w, `{"error": "category_id must be positive"}`, http.StatusBadRequest)
		return
	}

	rule, err := h.txService.CreateRule(r.Context(), userID, req.Keyword, req.CategoryID)
	if err != nil {
		http.Error(w, `{"error": "failed to create rule"}`, http.StatusInternalServerError)
		return
	}

	response := dto.CategoryRuleResponse{
		ID:         rule.ID,
		Keyword:    rule.Keyword,
		CategoryID: rule.CategoryID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetAll
// @Summary Получить правила пользователя
// @Description Получение списка правил категоризации пользователя
// @Tags category-rules
// @Produce json
// @Success 200 {array} dto.CategoryRuleResponse
// @Failure 401 {object} map[string]string "Неавторизован"
// @Router /api/v1/category-rules [get]
func (h *CategoryRuleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rules, err := h.txService.GetRules(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error": "failed to get rules"}`, http.StatusInternalServerError)
		return
	}

	// Получаем категории для названий
	categories, _ := h.txService.GetCategories(r.Context())
	categoryMap := make(map[int]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	response := make([]dto.CategoryRuleResponse, len(rules))
	for i, rule := range rules {
		response[i] = dto.CategoryRuleResponse{
			ID:         rule.ID,
			Keyword:    rule.Keyword,
			CategoryID: rule.CategoryID,
			Category:   categoryMap[rule.CategoryID],
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delete
// @Summary Удалить правило
// @Description Удаление правила категоризации по ID
// @Tags category-rules
// @Produce json
// @Param id path string true "ID правила"
// @Success 200 {object} map[string]string "Успешно"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 404 {object} map[string]string "Не найдено"
// @Router /api/v1/category-rules/{id} [delete]
func (h *CategoryRuleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error": "id is required"}`, http.StatusBadRequest)
		return
	}

	if err := h.txService.DeleteRule(r.Context(), userID, id); err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, `{"error": "rule not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error": "failed to delete rule"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "rule deleted successfully"})
}

// GetCategories
// @Summary Получить все категории
// @Description Получение списка всех категорий для использования в правилах
// @Tags categories
// @Produce json
// @Success 200 {array} dto.CategoryResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.txService.GetCategories(r.Context())
	if err != nil {
		http.Error(w, `{"error": "failed to get categories"}`, http.StatusInternalServerError)
		return
	}

	response := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = dto.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
