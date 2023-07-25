package http

import (
	"net/http"
	"payment-service/internal/service/catalogue"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"payment-service/internal/domain/product"
	"payment-service/pkg/server/response"
	"payment-service/pkg/store"
)

type ProductHandler struct {
	productService *catalogue.Service
}

func NewProductHandler(s *catalogue.Service) *ProductHandler {
	return &ProductHandler{productService: s}
}

func (h *ProductHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.list)
	r.Post("/", h.add)

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.delete)
	})

	return r
}

// List of products from the database
//
//	@Summary	List of products from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Success	200			{array}		response.Object
//	@Failure	500			{object}	response.Object
//	@Router		/products 	[get]
func (h *ProductHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.productService.ListProducts(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Add a new product to the database
//
//	@Summary	Add a new product to the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		request	body		product.Request	true	"body param"
//	@Success	200		{object}	response.Object
//	@Failure	400		{object}	response.Object
//	@Failure	500		{object}	response.Object
//	@Router		/products [post]
func (h *ProductHandler) add(w http.ResponseWriter, r *http.Request) {
	req := product.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	res, err := h.productService.AddProduct(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Read the product from the database
//
//	@Summary	Read the product from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"path param"
//	@Success	200	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/products/{id} [get]
func (h *ProductHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.productService.GetProduct(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Update the product in the database
//
//	@Summary	Update the product in the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string			true	"path param"
//	@Param		request	body	product.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/products/{id} [put]
func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := product.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	err := h.productService.UpdateProduct(r.Context(), id, req)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}
}

// Delete the product from the database
//
//	@Summary	Delete the product from the database
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"path param"
//	@Success	200
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/products/{id} [delete]
func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.productService.DeleteProduct(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}
}
