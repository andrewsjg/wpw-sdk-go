package main

import (
	"errors"
	"fmt"

	devclienttypes "github.com/WPTechInnovation/wpw-sdk-go/applications/dev-client/types"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
	log "github.com/sirupsen/logrus"
)

func mPrepareNewConsumer() error {

	fmt.Println("Preparing new consumer")

	if err := mInitNewDevice(); err != nil {
		return err
	}

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	if err := newHCECredential(); err != nil {
		return err
	}

	return nil
}

func newHCECredential() error {

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	fmt.Print("First Name: ")
	var firstName string
	if err := getUserInput(&firstName); err != nil {
		return err
	}

	fmt.Print("Last Name: ")
	var lastName string
	if err := getUserInput(&lastName); err != nil {
		return err
	}

	fmt.Print("Expiry month: ")
	var expMonth int32
	if err := getUserInput(&expMonth); err != nil {
		return err
	}

	fmt.Print("Expiry year: ")
	var expYear int32
	if err := getUserInput(&expYear); err != nil {
		return err
	}

	fmt.Print("CardNumber: ")
	var cardNumber string
	if err := getUserInput(&cardNumber); err != nil {
		return err
	}

	fmt.Print("Type: ")
	var cardType string
	if err := getUserInput(&cardType); err != nil {
		return err
	}

	fmt.Print("CVC: ")
	var cvc string
	if err := getUserInput(&cvc); err != nil {
		return err
	}

	card := types.HCECard{

		FirstName:  firstName,
		LastName:   lastName,
		ExpMonth:   expMonth,
		ExpYear:    expYear,
		CardNumber: cardNumber,
		Type:       cardType,
		Cvc:        cvc,
	}

	hceCard = card

	return nil
}

func mScanService() error {

	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}

	fmt.Print("Scan timeout in milliseconds: ")
	var timeout int
	if err := getUserInput(&timeout); err != nil {
		return err
	}

	log.Debug("pre scan for services")
	services, err := sdk.DeviceDiscovery(timeout)
	log.Debug("end scan for services")

	if err != nil {
		return err
	}

	for _, svc := range services {
		log.Debug("(%s:%d/%s) - %s", svc.Hostname, svc.PortNumber, svc.URLPrefix, svc.DeviceDescription)
	}
	return nil
}

func mAutoConsume() error {

	fmt.Println("Starting auto consume...")

	log.Debug("pre scan for services")
	services, err := sdk.DeviceDiscovery(deviceProfile.DeviceEntity.Consumer.ConsumerConfig.DeviceDiscoveryTimeout)
	log.Debug("end scan for services")

	if err != nil {

		return err
	}

	if len(services) >= 1 {

		var foundServiceIdx int = -1
		for i, service := range services {
			if service.ServerID == deviceProfile.DeviceEntity.Consumer.AutoConsume.DeviceUid {
				foundServiceIdx = i
				break
			}
		}

		if foundServiceIdx == -1 {
			fmt.Println("Could not find service - is the device id in the autoconsume section correct?")
		} else {

			fmt.Printf("Found Service:: (%s:%d/%s) - %s\n", services[foundServiceIdx].Hostname, services[foundServiceIdx].PortNumber, services[foundServiceIdx].URLPrefix, services[foundServiceIdx].DeviceDescription)

			log.Debug("Init consumer")

			// Could come from a config file..
			var pspConfig = make(map[string]string, 0)
			pspConfig[psp.CfgPSPName] = onlineworldpay.PSPName
			pspConfig[onlineworldpay.CfgAPIEndpoint] = "https://api.worldpay.com/v1"

			err := sdk.InitConsumer("http://", services[foundServiceIdx].Hostname, services[foundServiceIdx].PortNumber, services[foundServiceIdx].URLPrefix, services[foundServiceIdx].ServerID, deviceProfile.DeviceEntity.Consumer.HCECard, pspConfig)
			if err != nil {

				return err
			}

			log.Debug("Client created..")

			serviceDetails, err := sdk.RequestServices()

			if err != nil {

				return err
			}

			if len(serviceDetails) >= 1 {

				var foundDetailsIdx int = -1
				for i, serviceDetail := range serviceDetails {
					//fmt.Printf("%d - %s\n", serviceDetail.ServiceID, serviceDetail.ServiceDescription)
					if serviceDetail.ServiceID == deviceProfile.DeviceEntity.Consumer.AutoConsume.ServiceID {
						foundDetailsIdx = i
						break
					}
				}

				if foundDetailsIdx == -1 {
					fmt.Println("Could not find service details - is the service id in the autoconsume section correct?")
				} else {
					fmt.Printf("Selecting service: %d - %s\n", serviceDetails[foundDetailsIdx].ServiceID, serviceDetails[foundDetailsIdx].ServiceDescription)

					prices, err := sdk.GetServicePrices(serviceDetails[foundDetailsIdx].ServiceID)

					if err != nil {

						return err
					}

					if len(prices) >= 1 {

						var foundUnitIdIdx int = -1
						for i, price := range prices {
							//fmt.Printf("(%d) %s @ %d%s, %s (Unit id = %d)\n", price.ID, price.Description, price.PricePerUnit.Amount, price.PricePerUnit.CurrencyCode, price.UnitDescription, price.UnitID)
							if price.UnitID == deviceProfile.DeviceEntity.Consumer.AutoConsume.UnitID {
								foundUnitIdIdx = i
								break
							}
						}

						if foundUnitIdIdx == -1 {
							fmt.Println("Could not find unit id - is the unit id in the autoconsume section correct?")
						} else {

							fmt.Printf("Selecting price: (%d) %s @ %d%s, %s (Unit id = %d)\n", prices[foundUnitIdIdx].ID, prices[foundUnitIdIdx].Description, prices[foundUnitIdIdx].PricePerUnit.Amount, prices[foundUnitIdIdx].PricePerUnit.CurrencyCode, prices[foundUnitIdIdx].UnitDescription, prices[foundUnitIdIdx].UnitID)

							tpr, err := sdk.SelectService(serviceDetails[foundDetailsIdx].ServiceID, 1, prices[foundUnitIdIdx].ID)

							if err != nil {

								return err
							}

							log.Debug("Making payment of %d\n", tpr.TotalPrice)

							payResp, err := sdk.MakePayment(tpr)

							if err != nil {

								return err
							}

							_, err = sdk.BeginServiceDelivery(serviceDetails[foundDetailsIdx].ServiceID, *payResp.ServiceDeliveryToken, 1)

							if err != nil {

								fmt.Println("Error calling sdk.BeginServiceDelivery()")
								return err
							}

							fmt.Printf("Payment of %d made successfully\n", payResp.TotalPaid)
						}
					}
				}
			}
		}
	}
	return nil

}
