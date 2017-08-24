package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/wptechinnovation/wpw-sdk-go/applications/dev-client/types"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
)

var wpw wpwithin.WPWithin

func main() {

	initLog()

	wp, err := wpwithin.Initialise("Robo Car Service", "Car service provided by robot...")
	wpw = wp

	if err != nil {

		fmt.Println(err.Error())
		return
	}

	addService()

	eh := EventHandlerImpl{}
	wp.SetEventHandler(&eh)

	pspConfig := make(map[string]string, 0)
	pspConfig[psp.CfgPSPName] = onlineworldpay.PSPName
	pspConfig[onlineworldpay.CfgMerchantClientKey] = "T_C_97e8cbaa-14e0-4b1c-b2af-469daf8f1356"
	pspConfig[onlineworldpay.CfgMerchantServiceKey] = "T_S_3bdadc9c-54e0-4587-8d91-29813060fecd"
	pspConfig[psp.CfgHTEPrivateKey] = pspConfig[onlineworldpay.CfgMerchantServiceKey]
	pspConfig[psp.CfgHTEPublicKey] = pspConfig[onlineworldpay.CfgMerchantClientKey]
	pspConfig[onlineworldpay.CfgAPIEndpoint] = "https://api.worldpay.com/v1"

	err = wp.InitProducer(pspConfig)

	if err != nil {

		fmt.Println(err.Error())
	} else {

		fmt.Println("Start broadcast")

		// A timeout of 0 means run indefinitely
		wp.StartServiceBroadcast(0)
	}

	done := make(chan bool)
	fnForever := func() {
		for {
			time.Sleep(time.Second * 10)
		}
	}

	go fnForever()

	<-done // Block forever
}

func addService() {

	roboWash, _ := types.NewService()
	roboWash.Name = "RoboWash"
	roboWash.Description = "Car washed by robot"
	roboWash.ID = 1

	washPriceCar := types.Price{

		UnitID:          1,
		ID:              1,
		Description:     "Car wash",
		UnitDescription: "Single wash",
		PricePerUnit: &types.PricePerUnit{
			Amount:       500,
			CurrencyCode: "GBP",
		},
	}

	washPriceSUV := types.Price{

		UnitID:          1,
		ID:              2,
		Description:     "SUV Wash",
		UnitDescription: "Single wash",
		PricePerUnit: &types.PricePerUnit{
			Amount:       650,
			CurrencyCode: "GBP",
		},
	}

	roboWash.AddPrice(washPriceCar)
	roboWash.AddPrice(washPriceSUV)

	if wpw == nil {
		fmt.Println(errors.New(devclienttypes.ErrorDeviceNotInitialised).Error())
	}

	if err := wpw.AddService(roboWash); err != nil {

		fmt.Println(err.Error())
	}
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
