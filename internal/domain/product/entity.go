package product

import "time"

type Entity struct {
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	ID          string    `db:"id"`
	CategoryID  string    `db:"category_id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Measure     *string   `db:"measure"`
	ImageURL    *string   `db:"image_url"`
	Country     *string   `db:"country"`
	Barcode     *string   `db:"barcode"`
	Brand       *string   `db:"brand"`
}
