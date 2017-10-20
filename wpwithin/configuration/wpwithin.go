package configuration

import "fmt"

// WPWithin WorldpayWithin specific configuration
type WPWithin struct {
	WSLogEnable    bool
	WSLogPort      int
	WSLogLevel     string
	WSLogHost      string
}

// ParseConfig load in a Configuration and read it into WorldpayWithin specific config
func (wpw *WPWithin) ParseConfig(cfg Configuration) {

	enable, err := cfg.GetValue("wsLogEnable").ReadBool()
	if err != nil {
		fmt.Printf("Error parsing wsLogEnable as boolean: %s\n", err.Error())
		wpw.WSLogEnable = false
	} else {
		wpw.WSLogEnable = enable
	}

	if wpw.WSLogEnable {
		port, err := cfg.GetValue("wsLogPort").ReadInt()
		if err != nil {
			fmt.Printf("Error parsing wsLogPort as int: %s\n", err.Error())
		} else {
			wpw.WSLogPort = port
		}

		wpw.WSLogHost = cfg.GetValue("wsLogHost").Value
	}

	wpw.WSLogLevel = cfg.GetValue("wsLogLevel").Value
}
