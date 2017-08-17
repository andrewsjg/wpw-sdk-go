package types

// ServiceListResponse HTTP Message - List of available services from a server
type ServiceListResponse struct {
	ServerID string           `json:"serverID"`
	Services []ServiceDetails `json:"services"`
}
