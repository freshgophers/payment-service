package catalogue

import (
	"payment-service/internal/domain/category"
	"payment-service/internal/domain/product"
)

// Configuration is an alias for a function that will take in a pointer to a Service and modify it
type Configuration func(s *Service) error

// Service is an implementation of the Service
type Service struct {
	categoryRepository category.Repository
	productRepository  product.Repository

	categoryCache category.Cache
	productCache  product.Cache
}

// New takes a variable amount of Configuration functions and returns a new Service
// Each Configuration will be called in the order they are passed in
func New(configs ...Configuration) (s *Service, err error) {
	// Create the service
	s = &Service{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the service into the configuration function
		if err = cfg(s); err != nil {
			return
		}
	}
	return
}

// WithCategoryRepository applies a given category repository to the Service
func WithCategoryRepository(categoryRepository category.Repository) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.categoryRepository = categoryRepository
		return nil
	}
}

// WithProductRepository applies a given product repository to the Service
func WithProductRepository(productRepository product.Repository) Configuration {
	// Create the product repository, if we needed parameters, such as connection strings they could be inputted here
	return func(s *Service) error {
		s.productRepository = productRepository
		return nil
	}
}

// WithCategoryCache applies a given category cache to the Service
func WithCategoryCache(categoryCache category.Cache) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.categoryCache = categoryCache
		return nil
	}
}

// WithProductCache applies a given product cache to the Service
func WithProductCache(productCache product.Cache) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.productCache = productCache
		return nil
	}
}
