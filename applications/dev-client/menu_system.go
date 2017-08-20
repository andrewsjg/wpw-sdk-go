package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	devclienttypes "github.com/wptechinnovation/worldpay-within-sdk/applications/dev-client/types"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin"
)

var sdk wpwithin.WPWithin
var menuItems []MenuItem
var deviceProfile devclienttypes.DeviceProfile

type MenuItem struct {
	Label  string
	Action func() error
}

func doUI() {

	menuItems = make([]MenuItem, 0)

	menuItems = append(menuItems, MenuItem{"-------------------- GENERAL  --------------------", mInvalidSelection})
	menuItems = append(menuItems, MenuItem{"Init new device", mInitNewDevice})
	menuItems = append(menuItems, MenuItem{"Start RPC Service", mStartRPCService})
	menuItems = append(menuItems, MenuItem{"Get device info", mGetDeviceInfo})
	menuItems = append(menuItems, MenuItem{"Load device profile", mLoadDeviceProfile})
	menuItems = append(menuItems, MenuItem{"Reset session", mResetSessionState})
	menuItems = append(menuItems, MenuItem{"-------------------- PRODUCER --------------------", mInvalidSelection})
	menuItems = append(menuItems, MenuItem{"Init new producer", mNewProducer})
	menuItems = append(menuItems, MenuItem{"Add RoboWash service", mAddRoboWashService})
	menuItems = append(menuItems, MenuItem{"Add RoboAir service", mAddRoboAirService})
	menuItems = append(menuItems, MenuItem{"Start service broadcast", mStartBroadcast})
	menuItems = append(menuItems, MenuItem{"Stop broadcast", mStopBroadcast})
	menuItems = append(menuItems, MenuItem{"-------------------- CONSUMER --------------------", mInvalidSelection})
	menuItems = append(menuItems, MenuItem{"Prepare new consumer", mPrepareNewConsumer})
	menuItems = append(menuItems, MenuItem{"Scan services", mScanService})
	menuItems = append(menuItems, MenuItem{"Auto consume from profile info", mAutoConsume})
	menuItems = append(menuItems, MenuItem{"--------------------------------------------------", mInvalidSelection})
	menuItems = append(menuItems, MenuItem{"Exit", mQuit})

	renderMenu()
}
func mInvalidSelection() error {

	return errors.New(devclienttypes.ErrorInvalidMenuSelection)
}

func promptContinue() bool {

	fmt.Print("Continue (y/n): ")
	var cont string
	_, err := fmt.Scanf("%s\n", &cont)

	if err != nil {

		fmt.Println("Please enter \"y\" or \"n\"")
		return promptContinue()
	} else if strings.EqualFold(cont, "y") {

		return true
	} else {

		return false
	}
}

func renderMenu() {

	fmt.Println("----------------------------- Worldpay Within SDK Client ----------------------------")

	for i, item := range menuItems {

		fmt.Printf("%d - %s\n", i, item.Label)
	}

	fmt.Println("-------------------------------------------------------------------------------------")

	fmt.Print("Please select choice: ")
	var input string

	if _, err := fmt.Scanln(&input); err != nil {

		fmt.Printf("Selection error: %q\n", err.Error())
		renderMenu()
		return
	}

	inputInt, err := strconv.Atoi(input)

	if err != nil {
		fmt.Println("Please type an integer choice!")
		renderMenu()
		return
	}

	if inputInt >= len(menuItems) {
		fmt.Println("Index out of bounds!")
		renderMenu()
		return
	}

	if err = menuItems[inputInt].Action(); err != nil {

		fmt.Println(err.Error())
	}

	if promptContinue() {
		renderMenu()
	} else {
		fmt.Println("")
		fmt.Println("Goodbye...")
		os.Exit(1)
	}

}

func setUp() {

	fmt.Println("Setup")
}

func mQuit() error {

	fmt.Println("")
	fmt.Println("Goodbye...")

	os.Exit(1)

	return nil
}

func getUserInput(input interface{}) error {

	var err error

	switch t := input.(type) {
	case *int:
		_, err = fmt.Scanf("%d\n", input)
	case *int32:
		_, err = fmt.Scanf("%d\n", input)
	case *string:
		_, err = fmt.Scanf("%s\n", input)
	default:
		fmt.Printf("unexpected type %T", t)
	}

	if err != nil {
		return err
	}

	return nil
}
