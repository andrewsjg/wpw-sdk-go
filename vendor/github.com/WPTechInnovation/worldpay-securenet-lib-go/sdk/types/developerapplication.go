package types

// DeveloperApplication represents information about an application integration
type DeveloperApplication struct {
	// DeveloperID is the ID of the application developer
	DeveloperID int32 `json:"developerId"`
	// Version of integrators application
	Version string `json:"version"`
}
