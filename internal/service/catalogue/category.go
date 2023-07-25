package catalogue

import (
	"context"
	"payment-service/internal/domain/category"
)

func (s *Service) ListCategories(ctx context.Context) (res []category.Response, err error) {
	data, err := s.categoryRepository.Select(ctx)
	if err != nil {
		return
	}
	res = category.ParseFromEntities(data)

	return
}

func (s *Service) AddCategory(ctx context.Context, req category.Request) (res category.Response, err error) {
	data := category.Entity{
		ParentID: req.ParentID,
		Name:     &req.Name,
	}

	data.ID, err = s.categoryRepository.Create(ctx, data)
	if err != nil {
		return
	}
	res = category.ParseFromEntity(data)

	return
}

func (s *Service) GetCategory(ctx context.Context, id string) (res category.Response, err error) {
	parent, err := s.categoryRepository.Get(ctx, id)
	if err != nil {
		return
	}
	res = category.ParseFromEntity(parent)

	children, err := s.categoryRepository.SelectByParentID(ctx, parent.ID)
	if err != nil {
		return
	}
	res.Children = category.ParseFromEntities(children)

	return
}

func (s *Service) UpdateCategory(ctx context.Context, id string, req category.Request) (err error) {
	data := category.Entity{
		Name: &req.Name,
	}
	return s.categoryRepository.Update(ctx, id, data)
}

func (s *Service) DeleteCategory(ctx context.Context, id string) (err error) {
	return s.categoryRepository.Delete(ctx, id)
}
