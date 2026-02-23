package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/internal/domain/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &postgresUserRepository{pool: pool}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, google_id, global_currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.GoogleID,
		user.GlobalCurrency,
		now,
		now,
	)

	if err != nil {
		return err
	}

	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, google_id, global_currency, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.GoogleID,
		&user.GlobalCurrency,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, google_id, global_currency, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.GoogleID,
		&user.GlobalCurrency,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresUserRepository) GetByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, google_id, global_currency, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.GoogleID,
		&user.GlobalCurrency,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, google_id = $4, global_currency = $5, updated_at = $6
		WHERE id = $1
	`

	user.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.GoogleID,
		user.GlobalCurrency,
		user.UpdatedAt,
	)

	return err
}

func (r *postgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}
