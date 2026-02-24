package repository

import (
	"context"
	"time"

	"github.com/gibbon/finace-dashboard/internal/domain/model"
	"github.com/gibbon/finace-dashboard/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresCategoryRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresCategoryRepository(pool *pgxpool.Pool) repository.CategoryRepository {
	return &postgresCategoryRepository{pool: pool}
}

func (r *postgresCategoryRepository) GetAll(ctx context.Context) ([]*model.Category, error) {
	query := `SELECT id, name, is_default, created_at FROM categories ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		cat := &model.Category{}
		err := rows.Scan(&cat.ID, &cat.Name, &cat.IsDefault, &cat.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (r *postgresCategoryRepository) GetByID(ctx context.Context, id int) (*model.Category, error) {
	query := `SELECT id, name, is_default, created_at FROM categories WHERE id = $1`

	cat := &model.Category{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&cat.ID, &cat.Name, &cat.IsDefault, &cat.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return cat, nil
}

func (r *postgresCategoryRepository) GetDefault(ctx context.Context) ([]*model.Category, error) {
	query := `SELECT id, name, is_default, created_at FROM categories WHERE is_default = true ORDER BY name`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		cat := &model.Category{}
		err := rows.Scan(&cat.ID, &cat.Name, &cat.IsDefault, &cat.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

type postgresUserCategoryRuleRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserCategoryRuleRepository(pool *pgxpool.Pool) repository.UserCategoryRuleRepository {
	return &postgresUserCategoryRuleRepository{pool: pool}
}

func (r *postgresUserCategoryRuleRepository) Create(ctx context.Context, rule *model.UserCategoryRule) error {
	query := `
		INSERT INTO user_category_rules (id, user_id, keyword, category_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()

	_, err := r.pool.Exec(ctx, query,
		rule.ID,
		rule.UserID,
		rule.Keyword,
		rule.CategoryID,
		rule.CreatedAt,
	)

	return err
}

func (r *postgresUserCategoryRuleRepository) GetByUserID(ctx context.Context, userID string) ([]*model.UserCategoryRule, error) {
	query := `SELECT id, user_id, keyword, category_id, created_at FROM user_category_rules WHERE user_id = $1`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*model.UserCategoryRule
	for rows.Next() {
		rule := &model.UserCategoryRule{}
		err := rows.Scan(&rule.ID, &rule.UserID, &rule.Keyword, &rule.CategoryID, &rule.CreatedAt)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

func (r *postgresUserCategoryRuleRepository) GetByKeyword(ctx context.Context, userID, keyword string) (*model.UserCategoryRule, error) {
	query := `SELECT id, user_id, keyword, category_id, created_at FROM user_category_rules WHERE user_id = $1 AND keyword = $2`

	rule := &model.UserCategoryRule{}
	err := r.pool.QueryRow(ctx, query, userID, keyword).Scan(&rule.ID, &rule.UserID, &rule.Keyword, &rule.CategoryID, &rule.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return rule, nil
}

func (r *postgresUserCategoryRuleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_category_rules WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *postgresUserCategoryRuleRepository) Update(ctx context.Context, rule *model.UserCategoryRule) error {
	query := `
		UPDATE user_category_rules
		SET keyword = $2, category_id = $3
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, rule.ID, rule.Keyword, rule.CategoryID)
	return err
}
