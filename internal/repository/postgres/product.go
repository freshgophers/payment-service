package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"payment-service/internal/domain/product"
	"strings"

	"github.com/jmoiron/sqlx"

	"payment-service/pkg/store"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (s *ProductRepository) Select(ctx context.Context) (dest []product.Entity, err error) {
	query := `
		SELECT id,category_id, name,description,measure,image_url,country,barcode,brand
		FROM products`

	err = s.db.SelectContext(ctx, &dest, query)

	return
}

func (s *ProductRepository) Create(ctx context.Context, data product.Entity) (id string, err error) {
	query := `
		INSERT INTO products (id,category_id, name,description,measure )
		VALUES ($1, $2, $3,$4,$5)
		RETURNING id`

	args := []any{data.ID, data.CategoryID, data.Name, data.Description, data.Measure}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&id)

	return
}

func (s *ProductRepository) Get(ctx context.Context, id string) (dest product.Entity, err error) {
	query := `
		SELECT id,category_id, name,description,measure,image_url,country,barcode,brand
		FROM products
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

func (s *ProductRepository) Update(ctx context.Context, id string, data product.Entity) (err error) {
	sets, args := s.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE products SET %s WHERE id=$%d", strings.Join(sets, ", "), len(args))
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

func (s *ProductRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE 
		FROM products
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

func (s *ProductRepository) prepareArgs(data product.Entity) (sets []string, args []any) {
	if data.Description != nil {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("description=$%d", len(args)))
	}

	if data.Measure != nil {
		args = append(args, data.Measure)
		sets = append(sets, fmt.Sprintf("measure=$%d", len(args)))
	}

	return
}
