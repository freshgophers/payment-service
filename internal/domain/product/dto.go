package product

import (
	"errors"
	"net/http"
)

type Request struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Measure     string `json:"measure"`
	ImageURL    string `json:"image_url"`
	Country     string `json:"country"`
	Barcode     string `json:"barcode"`
	Brand       string `json:"brand"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.CategoryID == "" {
		return errors.New("CategoryID: cannot be blank")
	}

	if s.Barcode == "" {
		return errors.New("Barcode: cannot be blank")
	}

	if s.Name == "" {
		return errors.New("Name: cannot be blank")
	}

	if s.Measure == "" {
		return errors.New("Measure: cannot be blank")
	}

	return nil
}

type Response struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Measure     string `json:"measure"`
	ImageURL    string `json:"image_url"`
	Country     string `json:"country"`
	Barcode     string `json:"barcode"`
	Brand       string `json:"brand"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		ID:          data.ID,
		CategoryID:  data.CategoryID,
		Name:        data.Name,
		Description: *data.Description,
		Measure:     *data.Measure,
		ImageURL:    *data.ImageURL,
		Country:     *data.Country,
		Barcode:     *data.Barcode,
		Brand:       *data.Brand,
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
