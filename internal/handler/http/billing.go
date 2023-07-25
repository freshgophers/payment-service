package http

import (
	"net/http"
	"payment-service/internal/domain/billing"
	"payment-service/internal/service/payment"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"payment-service/pkg/server/response"
)

type BillingHandler struct {
	Billing *payment.Service
}

func NewBilling(s *payment.Service) *BillingHandler {
	return &BillingHandler{Billing: s}
}

func (h *BillingHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.add)

	return r
}

// Add a new billing to the database
//
//	@Summary	Add a new billing to the database
//	@Tags		billings
//	@Accept		json
//	@Produce	json
//	@Param		request	body		billing.Request	true	"body param"
//	@Success	200		{object}	response.Object
//	@Failure	400		{object}	response.Object
//	@Failure	500		{object}	response.Object
//	@Router		/billings [post]
func (h *BillingHandler) add(w http.ResponseWriter, r *http.Request) {
	req := billing.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	res, err := h.Billing.AddBilling(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}
