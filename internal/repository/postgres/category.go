package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service/internal/domain/category"
	"strings"

	"github.com/jmoiron/sqlx"

	"payment-service/pkg/store"
)

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (s *CategoryRepository) Select(ctx context.Context) (dest []category.Entity, err error) {
	query := `
		SELECT id, name
		FROM categories`

	err = s.db.SelectContext(ctx, &dest, query)

	return
}

func (s *CategoryRepository) SelectByParentID(ctx context.Context, parentID string) (dest []category.Entity, err error) {
	query := `
		SELECT id, name
		FROM categories
		WHERE parent_id=$1`

	args := []any{parentID}

	err = s.db.SelectContext(ctx, &dest, query, args...)

	return
}

func (s *CategoryRepository) Create(ctx context.Context, data category.Entity) (id string, err error) {
	query := `
		INSERT INTO categories (parent_id, name)
		VALUES ($1, $2)
		RETURNING id`

	args := []any{data.ParentID, data.Name}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *CategoryRepository) Get(ctx context.Context, id string) (dest category.Entity, err error) {
	query := `
		SELECT id, name
		FROM categories
		WHERE id=$1`

	args := []any{id}

	if err = s.db.GetContext(ctx, &dest, query, args...); err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = store.ErrorNotFound
	}

	return
}

func (s *CategoryRepository) Update(ctx context.Context, id string, data category.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE categories SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
		_, err = s.db.ExecContext(ctx, query, args...)
		if err != nil && err != sql.ErrNoRows {
			return
		}

		if err == sql.ErrNoRows {
			err = store.ErrorNotFound
		}
	}

	return
}

func (s *CategoryRepository) prepareArgs(data category.Entity) (sets []string, args []any) {
	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	return
}

func (s *CategoryRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE 
		FROM categories
		WHERE id=$1`

	args := []any{id}

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		err = store.ErrorNotFound
	}

	return
}
