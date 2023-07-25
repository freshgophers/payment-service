package category

import (
	"payment-service/pkg/store/postgres"
	"time"
)

type Entity struct {
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	ID        string         `db:"id"`
	ParentID  string         `db:"parent_id"`
	Name      *string        `db:"name"`
	Child     postgres.Array `db:"child"`
}
