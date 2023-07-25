package http

import (
	"net/http"
	"payment-service/internal/domain/category"
	"payment-service/internal/service/catalogue"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"payment-service/pkg/server/response"
	"payment-service/pkg/store"
)

type CategoryHandler struct {
	Category *catalogue.Service
}

func NewCategory(s *catalogue.Service) *CategoryHandler {
	return &CategoryHandler{Category: s}
}

func (h *CategoryHandler) Routes() chi.Router {
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

// List of categories from the database
//
//	@Summary	List of categories from the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Success	200				{array}		response.Object
//	@Failure	500				{object}	response.Object
//	@Router		/categories 	[get]
func (h *CategoryHandler) list(w http.ResponseWriter, r *http.Request) {
	res, err := h.Category.ListCategories(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Add a new category to the database
//
//	@Summary	Add a new category to the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		request	body		category.Request	true	"body param"
//	@Success	200		{object}	response.Object
//	@Failure	400		{object}	response.Object
//	@Failure	500		{object}	response.Object
//	@Router		/categories [post]
func (h *CategoryHandler) add(w http.ResponseWriter, r *http.Request) {
	req := category.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	res, err := h.Category.AddCategory(r.Context(), req)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, res)
}

// Read the category from the database
//
//	@Summary	Read the category from the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"path param"
//	@Success	200	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/categories/{id} [get]
func (h *CategoryHandler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	res, err := h.Category.GetCategory(r.Context(), id)
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

// Update the category in the database
//
//	@Summary	Update the category in the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		id		path	string				true	"path param"
//	@Param		request	body	category.Request	true	"body param"
//	@Success	200
//	@Failure	400	{object}	response.Object
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/categories/{id} [put]
func (h *CategoryHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	req := category.Request{}
	if err := render.Bind(r, &req); err != nil {
		response.BadRequest(w, r, err, req)
		return
	}

	err := h.Category.UpdateCategory(r.Context(), id, req)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}
}

// Delete the category from the database
//
//	@Summary	Delete the category from the database
//	@Tags		categories
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"path param"
//	@Success	200
//	@Failure	404	{object}	response.Object
//	@Failure	500	{object}	response.Object
//	@Router		/categories/{id} [delete]
func (h *CategoryHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.Category.DeleteCategory(r.Context(), id)
	if err != nil && err != store.ErrorNotFound {
		response.InternalServerError(w, r, err)
		return
	}

	if err == store.ErrorNotFound {
		response.NotFound(w, r, err)
		return
	}
}
