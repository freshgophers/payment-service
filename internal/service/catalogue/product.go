package catalogue

import (
	"context"
	"payment-service/internal/domain/product"
)

func (s *Service) ListProducts(ctx context.Context) (res []product.Response, err error) {
	data, err := s.productRepository.Select(ctx)
	if err != nil {
		return
	}
	res = product.ParseFromEntities(data)

	return
}

func (s *Service) AddProduct(ctx context.Context, req product.Request) (res product.Response, err error) {
	data := product.Entity{
		ID:          req.ID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: &req.Description,
		Measure:     &req.Measure,
		ImageURL:    &req.ImageURL,
		Country:     &req.Country,
		Barcode:     &req.Barcode,
		Brand:       &req.Brand,
	}

	data.ID, err = s.productRepository.Create(ctx, data)
	if err != nil {
		return
	}
	res = product.ParseFromEntity(data)

	return
}

func (s *Service) GetProduct(ctx context.Context, id string) (res product.Response, err error) {
	data, err := s.productRepository.Get(ctx, id)
	if err != nil {
		return
	}
	res = product.ParseFromEntity(data)

	return
}

func (s *Service) UpdateProduct(ctx context.Context, id string, req product.Request) (err error) {
	data := product.Entity{
		Description: &req.Description,
		Measure:     &req.Measure,
		ImageURL:    &req.ImageURL,
		Country:     &req.Country,
		Barcode:     &req.Barcode,
		Brand:       &req.Brand,
	}
	return s.productRepository.Update(ctx, id, data)
}

func (s *Service) DeleteProduct(ctx context.Context, id string) (err error) {
	return s.productRepository.Delete(ctx, id)
}
