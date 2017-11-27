package types

// Device details of a device
type Device struct {
	UID         string
	Name        string
	Description string
	Services    map[int]*Service
	IPv4Address string
}

// NewDevice create a new device
func NewDevice(name, description, uid, ipv4Address, currencyCode string) (*Device, error) {

	result := &Device{
		Name:        name,
		Description: description,
		UID:         uid,
		IPv4Address: ipv4Address,
	}

	return result, nil
}
