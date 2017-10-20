package types

import "fmt"

// Service is a service offered by a producer
type Service struct {
	ID          int           `json:"serviceID"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	ServiceType string        `json:"serviceType"`
	Prices      map[int]Price `json:"prices"`
}

// NewService creates a new instance of Service
func NewService() (*Service, error) {

	result := &Service{}

	result.Prices = make(map[int]Price, 0)

	return result, nil
}

// AddPrice add a price to a service
func (service *Service) AddPrice(price Price) error {

	if _, exists := service.Prices[price.ID]; exists {

		return fmt.Errorf("A price with that ID (%d) already exists for that service.", price.ID)
	}

	service.Prices[price.ID] = price

	return nil
}

// RemovePrice removes a price from a service
func (service *Service) RemovePrice(price Price) error {

	delete(service.Prices, price.ID)

	return nil
}
