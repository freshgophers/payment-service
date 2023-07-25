package billing

import (
	"errors"
	"net/http"
)

type Request struct {
	CorrelationID   string `json:"correlation_id"`
	Source          string `json:"source"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
	Name            string `json:"name"`
	TerminalID      string `json:"terminal_id"`
	InvoiceID       string `json:"invoice_id"`
	Description     string `json:"description"`
	AccountID       string `json:"account_id"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Backlink        string `json:"backlink"`
	FailureBacklink string `json:"failure_backlink"`
	PostLink        string `json:"post_link"`
	FailurePostLink string `json:"failure_post_link"`
	Language        string `json:"language"`
	PaymentType     string `json:"payment_type"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.Name == "" {
		return errors.New("name: cannot be blank")
	}

	return nil
}

type Response struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:   data.ID,
		Link: "https://freshgopher-account-service.onrender.com/api/v1/invoices/" + data.ID + "/pay",
	}

	return
}

func ParseFromEntities(data []Entity) (res []Response) {
	res = make([]Response, 0)
	for _, object := range data {
		res = append(res, ParseFromEntity(object))
	}
	return
}
