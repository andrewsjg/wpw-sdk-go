package utils

import (
	"errors"
	"fmt"
	"net"
	"reflect"
)

// FirstExternalIPv4 Return the IPv4 external address of this device.
// Note external does not necessarily mean WAN IP. On most networks it will be the LAN IP of device as opposed
// to internal localhost address (127.0.0.1)
func FirstExternalIPv4() (net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if iface.Flags&net.FlagBroadcast == 0 {
			fmt.Printf("Interface %s does not support broadcast - skipping\n", iface.Name)
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ip, nil
		}
	}
	return nil, errors.New("Device does not appear to be network connected.")
}

// NetMask get current netmask
func NetMask() (net.IPMask, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {

			var ip net.IPMask

			switch v := addr.(type) {

			case *net.IPNet:

				if v.IP.To4() != nil {

					ones, _ := v.Mask.Size()

					i := ones / 8

					for j := 1; j <= i; j++ {

						fmt.Print("255")

						if j < i {

							fmt.Print(".")
						} else {

							fmt.Println(".0")
						}
					}

					fmt.Printf("V %s\n", v)
					fmt.Printf("Type %s\n", reflect.TypeOf(v))
					fmt.Println("******************")
				}
			}
			fmt.Println("--------------------")
			if ip == nil {
				fmt.Print("Not IPV4\n")
				continue // not an ipv4 address
			}

			return ip, nil
		}
	}
	return nil, errors.New("Unable to calculate netmask")
}
