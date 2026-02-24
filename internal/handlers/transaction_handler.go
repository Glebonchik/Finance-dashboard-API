package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/gibbon/finace-dashboard/internal/dto"
	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/internal/middleware"
	"github.com/gibbon/finace-dashboard/internal/service"
)

// TransactionHandler обрабатывает HTTP запросы для транзакций
type TransactionHandler struct {
	txService service.TransactionService
}

// NewTransactionHandler создаёт новый TransactionHandler
func NewTransactionHandler(txService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		txService: txService,
	}
}

// Create
// @Summary Создать новую транзакцию
// @Description Создание новой транзакции с автоматической категоризацией
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body dto.CreateTransactionRequest true "Данные транзакции"
// @Success 201 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Router /api/v1/transactions [post]
func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Amount <= 0 {
		http.Error(w, `{"error": "amount must be positive"}`, http.StatusBadRequest)
		return
	}

	if req.Description == "" {
		http.Error(w, `{"error": "description is required"}`, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, use RFC3339"}`, http.StatusBadRequest)
		return
	}

	tx := &model.Transaction{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		Date:        date,
		PlaceName:   req.PlaceName,
		PlaceLat:    req.PlaceLat,
		PlaceLon:    req.PlaceLon,
	}

	created, err := h.txService.Create(r.Context(), userID, tx)
	if err != nil {
		http.Error(w, `{"error": "failed to create transaction"}`, http.StatusInternalServerError)
		return
	}

	response := h.toTransactionResponse(created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetByID
// @Summary Получить транзакцию по ID
// @Description Получение данных транзакции по идентификатору
// @Tags transactions
// @Produce json
// @Param id path string true "ID транзакции"
// @Success 200 {object} dto.TransactionResponse
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 404 {object} map[string]string "Не найдено"
// @Router /api/v1/transactions/{id} [get]
func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	tx, err := h.txService.GetByID(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, `{"error": "transaction not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error": "failed to get transaction"}`, http.StatusInternalServerError)
		return
	}

	response := h.toTransactionResponse(tx)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAll
// @Summary Получить список транзакций
// @Description Получение списка транзакций пользователя с фильтрацией и пагинацией
// @Tags transactions
// @Produce json
// @Param category_id query int false "ID категории"
// @Param from_date query string false "Дата от (RFC3339)"
// @Param to_date query string false "Дата до (RFC3339)"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} dto.TransactionsListResponse
// @Failure 401 {object} map[string]string "Неавторизован"
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	filter := model.TransactionFilter{
		UserID: userID,
		Limit:  20,
		Offset: 0,
	}

	// Парсинг параметров
	if categoryID := r.URL.Query().Get("category_id"); categoryID != "" {
		if id, err := strconv.Atoi(categoryID); err == nil {
			filter.CategoryID = &id
		}
	}

	if fromDate := r.URL.Query().Get("from_date"); fromDate != "" {
		if date, err := time.Parse(time.RFC3339, fromDate); err == nil {
			filter.FromDate = &date
		}
	}

	if toDate := r.URL.Query().Get("to_date"); toDate != "" {
		if date, err := time.Parse(time.RFC3339, toDate); err == nil {
			filter.ToDate = &date
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil && o >= 0 {
			filter.Offset = o
		}
	}

	transactions, err := h.txService.GetByUserID(r.Context(), filter)
	if err != nil {
		http.Error(w, `{"error": "failed to get transactions"}`, http.StatusInternalServerError)
		return
	}

	total, err := h.txService.(interface{ GetTotalCount(context.Context, string) (int64, error) }).GetTotalCount(r.Context(), userID)
	if err != nil {
		total = int64(len(transactions))
	}

	response := dto.TransactionsListResponse{
		Transactions: h.toTransactionResponses(transactions),
		Total:        total,
		Limit:        filter.Limit,
		Offset:       filter.Offset,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update
// @Summary Обновить транзакцию
// @Description Обновление данных транзакции
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "ID транзакции"
// @Param request body dto.UpdateTransactionRequest true "Данные транзакции"
// @Success 200 {object} dto.TransactionResponse
// @Failure 400 {object} map[string]string "Некорректные данные"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 404 {object} map[string]string "Не найдено"
// @Router /api/v1/transactions/{id} [put]
func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UpdateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		http.Error(w, `{"error": "invalid date format, use RFC3339"}`, http.StatusBadRequest)
		return
	}

	tx := &model.Transaction{
		ID:          id,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		Date:        date,
		PlaceName:   req.PlaceName,
		PlaceLat:    req.PlaceLat,
		PlaceLon:    req.PlaceLon,
		CategoryID:  req.CategoryID,
		IsConfirmed: req.IsConfirmed,
	}

	updated, err := h.txService.Update(r.Context(), userID, tx)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, `{"error": "transaction not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error": "failed to update transaction"}`, http.StatusInternalServerError)
		return
	}

	response := h.toTransactionResponse(updated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delete
// @Summary Удалить транзакцию
// @Description Удаление транзакции по идентификатору
// @Tags transactions
// @Produce json
// @Param id path string true "ID транзакции"
// @Success 200 {object} map[string]string "Успешно"
// @Failure 401 {object} map[string]string "Неавторизован"
// @Failure 404 {object} map[string]string "Не найдено"
// @Router /api/v1/transactions/{id} [delete]
func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	if err := h.txService.Delete(r.Context(), userID, id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, `{"error": "transaction not found"}`, http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error": "failed to delete transaction"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "transaction deleted successfully"})
}

func (h *TransactionHandler) toTransactionResponse(tx *model.Transaction) *dto.TransactionResponse {
	response := &dto.TransactionResponse{
		ID:          tx.ID,
		Amount:      tx.Amount,
		Currency:    tx.Currency,
		Description: tx.Description,
		Date:        tx.Date,
		PlaceName:   tx.PlaceName,
		PlaceLat:    tx.PlaceLat,
		PlaceLon:    tx.PlaceLon,
		CategoryID:  tx.CategoryID,
		IsConfirmed: tx.IsConfirmed,
		CreatedAt:   tx.CreatedAt,
		UpdatedAt:   tx.UpdatedAt,
	}
	return response
}

func (h *TransactionHandler) toTransactionResponses(txs []*model.Transaction) []*dto.TransactionResponse {
	responses := make([]*dto.TransactionResponse, len(txs))
	for i, tx := range txs {
		responses[i] = h.toTransactionResponse(tx)
	}
	return responses
}
