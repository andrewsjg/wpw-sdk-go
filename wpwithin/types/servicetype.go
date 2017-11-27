package types

// Price represents price
type ServiceType struct {
	ServiceType string `json:"serviceTypes"`
}

// NewPrice create a new instance of Price
func NewServiceType() (*ServiceType, error) {

	result := &ServiceType{}

	return result, nil
}
