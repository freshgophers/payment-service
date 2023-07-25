package payment

import (
	"payment-service/internal/domain/billing"
)

// Configuration is an alias for a function that will take in a pointer to a Service and modify it
type Configuration func(s *Service) error

// Service is an implementation of the Service
type Service struct {
	billingRepository billing.Repository
	billingCache      billing.Cache
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

// WithBillingRepository applies a given category repository to the Service
func WithBillingRepository(billingRepository billing.Repository) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.billingRepository = billingRepository
		return nil
	}
}

// WithBillingCache applies a given category cache to the Service
func WithBillingCache(billingCache billing.Cache) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.billingCache = billingCache
		return nil
	}
}
