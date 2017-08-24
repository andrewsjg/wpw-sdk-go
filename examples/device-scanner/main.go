package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin"
)

// Flags
var flagScanTimeout int

// App vars
var wpw wpwithin.WPWithin

func init() {

	flag.IntVar(&flagScanTimeout, "scantimeout", 2000, "Number of milliseconds to scan for. 0 = infinite.")
}

func main() {

	flag.Parse()

	initLog()

	wp, err := wpwithin.Initialise("conorhwp-pi", "Conor Hacketts Raspberry Pi - DEVICE SCANNER")
	wpw = wp

	if err != nil {

		fmt.Println(err.Error())
		return
	}

	done := make(chan bool)
	fnForever := func() {
		for {

			fmt.Println("Scanning network for devices now...")
			fmt.Printf("Will scan for %d milliseconds\n", flagScanTimeout)

			devices, err := wpw.DeviceDiscovery(flagScanTimeout)

			if err != nil {
				fmt.Println("Error finding devices..")
				fmt.Println(err.Error())
			} else if len(devices) > 0 {

				fmt.Println("------------------------------------------------------------")
				fmt.Printf("Found %d devices\n", len(devices))

				for _, dev := range devices {

					fmt.Printf("[%s] [%s] %s @ %s:%d%s\n", dev.DeviceName, dev.ServerID, dev.DeviceDescription, dev.Hostname, dev.PortNumber, dev.URLPrefix)
				}
				fmt.Println("------------------------------------------------------------")
			}
		}
	}

	go fnForever()

	<-done // Block forever
}

func initLog() {

	log.SetFormatter(&log.JSONFormatter{})

	f, err := os.OpenFile("wpwithin.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {

		fmt.Println(err.Error())
	}

	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	log.Debug("Log is initialised.")
}
