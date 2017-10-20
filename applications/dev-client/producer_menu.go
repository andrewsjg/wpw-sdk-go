package main

import (
	"errors"
	"fmt"

	devclienttypes "github.com/WPTechInnovation/wpw-sdk-go/applications/dev-client/types"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
)

func mNewProducer() error {

	if err := mInitNewDevice(); err != nil {
		return err
	}

	fmt.Println("Add new HTE credentials")

	fmt.Print("Merchant Client Key: ")
	var merchantClientKey string
	if err := getUserInput(&merchantClientKey); err != nil {
		return err
	}

	fmt.Print("Merchant Service Key: ")
	var merchantServiceKey string
	if err := getUserInput(&merchantServiceKey); err != nil {
		return err
	}

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	defer fmt.Println("Initialised new producer")

	// Could come from a config file..
	var pspConfig = make(map[string]string, 0)
	return sdk.InitProducer(pspConfig)
}

func mAddRoboWashService() error {

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

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	if err := sdk.AddService(roboWash); err != nil {

		return err
	}

	fmt.Println("Added robowash service")

	return nil
}

func mAddRoboAirService() error {

	roboAir, _ := types.NewService()
	roboAir.Name = "RoboAir"
	roboAir.Description = "Car tyre pressure checked and topped up by robot"
	roboAir.ID = 2

	airSinglePrice := types.Price{

		UnitID:          1,
		ID:              1,
		Description:     "Measure and adjust pressure",
		UnitDescription: "Tyre",
		PricePerUnit: &types.PricePerUnit{
			Amount:       25,
			CurrencyCode: "GBP",
		},
	}

	airFourPrice := types.Price{

		UnitID:          2,
		ID:              2,
		Description:     "Measure and adjust pressure - four tyres for the price of three",
		UnitDescription: "4 Tyre",
		PricePerUnit: &types.PricePerUnit{
			Amount:       airSinglePrice.PricePerUnit.Amount * 3,
			CurrencyCode: "GBP",
		},
	}

	roboAir.AddPrice(airSinglePrice)
	roboAir.AddPrice(airFourPrice)

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	if err := sdk.AddService(roboAir); err != nil {

		return err
	}

	fmt.Println("Added roboair service")

	return nil
}

func mStartBroadcast() error {

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	fmt.Print("Broadcast timeout in milliseconds: ")
	var timeout int
	if err := getUserInput(&timeout); err != nil {
		return err
	}

	if err := sdk.StartServiceBroadcast(timeout); err != nil {
		return err
	}

	fmt.Println("Broadcast started...")
	return nil
}

func mStopBroadcast() error {

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	sdk.StopServiceBroadcast()

	fmt.Println("Broadcast stopped")
	return nil
}
