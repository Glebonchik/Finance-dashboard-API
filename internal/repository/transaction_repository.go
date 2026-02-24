package repository

import (
	"context"
	"time"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresTransactionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTransactionRepository(pool *pgxpool.Pool) repository.TransactionRepository {
	return &postgresTransactionRepository{pool: pool}
}

func (r *postgresTransactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, user_id, amount, currency, description, date,
			place_name, place_lat, place_lon, category_id, is_confirmed,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.pool.Exec(ctx, query,
		tx.ID,
		tx.UserID,
		tx.Amount,
		tx.Currency,
		tx.Description,
		tx.Date,
		tx.PlaceName,
		tx.PlaceLat,
		tx.PlaceLon,
		tx.CategoryID,
		tx.IsConfirmed,
		tx.CreatedAt,
		tx.UpdatedAt,
	)

	return err
}

func (r *postgresTransactionRepository) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, amount, currency, description, date,
		       place_name, place_lat, place_lon, category_id, is_confirmed,
		       created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	tx := &model.Transaction{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tx.ID,
		&tx.UserID,
		&tx.Amount,
		&tx.Currency,
		&tx.Description,
		&tx.Date,
		&tx.PlaceName,
		&tx.PlaceLat,
		&tx.PlaceLon,
		&tx.CategoryID,
		&tx.IsConfirmed,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return tx, nil
}

func (r *postgresTransactionRepository) GetByUserID(ctx context.Context, filter model.TransactionFilter) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, amount, currency, description, date,
		       place_name, place_lat, place_lon, category_id, is_confirmed,
		       created_at, updated_at
		FROM transactions
		WHERE user_id = $1
	`

	args := []interface{}{filter.UserID}
	argIndex := 2

	if filter.CategoryID != nil {
		query += " AND category_id = $" + string(rune('0'+argIndex))
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	if filter.FromDate != nil {
		query += " AND date >= $" + string(rune('0'+argIndex))
		args = append(args, *filter.FromDate)
		argIndex++
	}

	if filter.ToDate != nil {
		query += " AND date <= $" + string(rune('0'+argIndex))
		args = append(args, *filter.ToDate)
		argIndex++
	}

	query += " ORDER BY date DESC"

	if filter.Limit > 0 {
		query += " LIMIT $" + string(rune('0'+argIndex))
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += " OFFSET $" + string(rune('0'+argIndex))
		args = append(args, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		tx := &model.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Amount,
			&tx.Currency,
			&tx.Description,
			&tx.Date,
			&tx.PlaceName,
			&tx.PlaceLat,
			&tx.PlaceLon,
			&tx.CategoryID,
			&tx.IsConfirmed,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *postgresTransactionRepository) Update(ctx context.Context, tx *model.Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $2, currency = $3, description = $4, date = $5,
		    place_name = $6, place_lat = $7, place_lon = $8,
		    category_id = $9, is_confirmed = $10, updated_at = $11
		WHERE id = $1
	`

	tx.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx, query,
		tx.ID,
		tx.Amount,
		tx.Currency,
		tx.Description,
		tx.Date,
		tx.PlaceName,
		tx.PlaceLat,
		tx.PlaceLon,
		tx.CategoryID,
		tx.IsConfirmed,
		tx.UpdatedAt,
	)

	return err
}

func (r *postgresTransactionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM transactions WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *postgresTransactionRepository) GetTotalCount(ctx context.Context, userID string) (int64, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`
	var count int64
	err := r.pool.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}
