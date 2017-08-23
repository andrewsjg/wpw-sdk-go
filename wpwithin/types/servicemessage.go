package types

// BroadcastMessage is a message sent out during broadcasts to specify details of a producer on the network
type BroadcastMessage struct {
	DeviceDescription string `json:"deviceDescription"`
	Hostname          string `json:"hostname"`
	PortNumber        int    `json:"portNumber"`
	ServerID          string `json:"serverID"`
	URLPrefix         string `json:"urlPrefix"`
	Scheme            string `json:"scheme"`
	DeviceName		  string `json:"deviceName"`
}
