package payment

import (
	"context"
	"payment-service/internal/domain/billing"
)

func (s *Service) AddBilling(ctx context.Context, req billing.Request) (res billing.Response, err error) {
	data := billing.Entity{
		CorrelationID:   req.CorrelationID,
		Source:          req.Source,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Name:            req.Name,
		TerminalID:      req.TerminalID,
		InvoiceID:       req.InvoiceID,
		Description:     req.Description,
		AccountID:       req.AccountID,
		Email:           req.Email,
		Phone:           req.Phone,
		Backlink:        req.Backlink,
		FailureBacklink: req.FailureBacklink,
		PostLink:        req.PostLink,
		FailurePostLink: req.FailurePostLink,
		Language:        req.Language,
		PaymentType:     req.PaymentType,
	}
	data.ID, err = s.billingRepository.Create(ctx, data)
	if err != nil {
		return
	}
	res = billing.ParseFromEntity(data)

	return
}
