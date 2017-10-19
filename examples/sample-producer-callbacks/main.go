package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/WPTechInnovation/wpw-sdk-go/applications/dev-client/types"
	"github.com/WPTechInnovation/wpw-sdk-go/examples/exutils"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
	log "github.com/sirupsen/logrus"
)

var wpw wpwithin.WPWithin

func main() {

	initLog()

	cfgFileName := "sample-producer-callbacks.json"
	cfg, err := exutils.LoadConfiguration(cfgFileName)
	if err != nil {
		fmt.Println("error, failed to read config file", cfgFileName, ":", err)
		os.Exit(1)
	}

	wp, err := wpwithin.Initialise(cfg.DeviceName, "Car service provided by robot...")
	wpw = wp

	if err != nil {

		fmt.Println(err.Error())
		return
	}

	addService()

	eh := EventHandlerImpl{}
	wp.SetEventHandler(&eh)

	pspConfig := make(map[string]string, 0)

	pspConfig[psp.CfgPSPName] = cfg.PspConfig.PspName
	pspConfig[onlineworldpay.CfgMerchantClientKey] = cfg.PspConfig.MerchantClientKey
	pspConfig[onlineworldpay.CfgMerchantServiceKey] = cfg.PspConfig.MerchantServiceKey
	pspConfig[psp.CfgHTEPrivateKey] = cfg.PspConfig.HtePrivateKey
	pspConfig[psp.CfgHTEPublicKey] = cfg.PspConfig.HtePublicKey
	pspConfig[onlineworldpay.CfgAPIEndpoint] = cfg.PspConfig.ApiEndpoint

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
	roboWash.ServiceType = "wash"

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
