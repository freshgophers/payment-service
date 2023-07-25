package billing

import (
	"payment-service/pkg/store/postgres"
	"time"
)

type Entity struct {
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
	ID              string         `db:"id"`
	Child           postgres.Array `db:"child"`
	CorrelationID   string         `db:"correlation_id"`
	Source          string         `db:"source"`
	Amount          string         `db:"amount"`
	Currency        string         `db:"currency"`
	Name            string         `db:"name"`
	TerminalID      string         `db:"terminal_id"`
	InvoiceID       string         `db:"invoice_id"`
	Description     string         `db:"description"`
	AccountID       string         `db:"account_id"`
	Email           string         `db:"email"`
	Phone           string         `db:"phone"`
	Backlink        string         `db:"backlink"`
	FailureBacklink string         `db:"failure_backlink"`
	PostLink        string         `db:"post_link"`
	FailurePostLink string         `db:"failure_post_link"`
	Language        string         `db:"language"`
	PaymentType     string         `db:"payment_type"`
}
