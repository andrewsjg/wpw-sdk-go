package main

import (
	"fmt"
	"os"
	"time"

	"github.com/WPTechInnovation/wpw-sdk-go/examples/exutils"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay"
	log "github.com/sirupsen/logrus"
)

var wpw wpwithin.WPWithin

func convertToMoneyString(amt int, currencyCode string) string {
	currencyString := ""
	currencyStringAfter := ""

	if "GBP" == currencyCode {
		currencyString = "£ "
	} else if "USD" == currencyCode {
		currencyString = "$ "
	} else if "RON" == currencyCode {
		currencyStringAfter = " LEU"
	}

	return fmt.Sprintf(currencyString+"%.2f"+currencyStringAfter, float64(amt/100))
}

func main() {

	initLog()

	cfgFileName := "sample-consumer.json"
	cfg, err := exutils.LoadConfiguration(cfgFileName)
	if err != nil {
		fmt.Println("error, failed to read config file", cfgFileName, ":", err)
		os.Exit(1)
	}

	hceCard := cfg.HceCard

	wpw, err := wpwithin.Initialise(cfg.DeviceName, "A WPWithin client written in Go")

	errCheckExit(err)

	serviceMessages, err := wpw.DeviceDiscovery(20000)

	errCheckExit(err)

	if len(serviceMessages) == 0 {

		fmt.Printf("Found %d devices\n", len(serviceMessages))

		if len(serviceMessages) == 0 {

			fmt.Println("Quitting...")
			os.Exit(2)
		}
	}

	for _, sm := range serviceMessages {

		fmt.Printf("--------------%s--------------\n", sm.DeviceDescription)
		fmt.Printf("Description: %s\n", sm.ServerID)
		fmt.Printf("Service HTTP string: %s%s:%d%s\n", sm.Scheme, sm.Hostname, sm.PortNumber, sm.URLPrefix)

		if len(sm.ServiceTypes) > 0 {
			fmt.Printf("Supported service type(s): ")
			for _, st := range sm.ServiceTypes {
				fmt.Printf("%s ", st)
			}
			fmt.Println()
		}
		fmt.Println("-----------------------------------------------")
	}
	fmt.Println()
	fmt.Println()

	sm := serviceMessages[0]

	fmt.Printf("Will select device: [%s] %s\n", sm.ServerID, sm.DeviceDescription)

	pspConfig := make(map[string]string, 0)

	pspConfig[psp.CfgPSPName] = cfg.PspConfig.PspName
	pspConfig[onlineworldpay.CfgAPIEndpoint] = cfg.PspConfig.ApiEndpoint

	err = wpw.InitConsumer(sm.Scheme, sm.Hostname, sm.PortNumber, sm.URLPrefix, wpw.GetDevice().UID, &hceCard, pspConfig)

	errCheckExit(err)

	services, err := wpw.RequestServices()

	errCheckExit(err)

	if len(services) == 0 {

		fmt.Printf("Found %d services\n", len(services))

		if len(services) == 0 {

			fmt.Println("Quitting...")
			os.Exit(2)
		}
	}

	for _, service := range services {

		fmt.Printf("Id: %d Name: %s Description: %s\n", service.ServiceID, service.ServiceName, service.ServiceDescription)
	}
	fmt.Println()
	fmt.Println()

	service := services[0]

	fmt.Printf("Will request prices for %d - %s\n", service.ServiceID, service.ServiceName)

	prices, err := wpw.GetServicePrices(service.ServiceID)

	errCheckExit(err)

	if len(prices) == 0 {

		fmt.Printf("Found %d prices\n", len(prices))

		if len(prices) == 0 {

			fmt.Println("Quitting...")
			os.Exit(2)
		}
	}
	fmt.Println()
	fmt.Println()

	for _, price := range prices {

		fmt.Printf("[price] Id: %d Description: %s\n", price.ID, price.Description)
		fmt.Printf("[unit] Id: %d Description: %s\n", price.UnitID, price.UnitDescription)
		fmt.Printf("[pricePerUnit] Amount: %s CurrencyCode: %s\n", convertToMoneyString(price.PricePerUnit.Amount, price.PricePerUnit.CurrencyCode), price.PricePerUnit.CurrencyCode)
	}
	fmt.Println()
	fmt.Println()

	price := prices[0]
	numUnits := 10
	fmt.Printf("Will select %d units of price %s - %s\n", numUnits, convertToMoneyString(price.PricePerUnit.Amount, price.PricePerUnit.CurrencyCode), price.Description)

	tpr, err := wpw.SelectService(service.ServiceID, numUnits, price.ID)

	errCheckExit(err)

	fmt.Println("Did receive total price response:")
	fmt.Printf("ClientID: %s\n", tpr.ClientID)
	fmt.Printf("CurrencyCode: %s\n", tpr.CurrencyCode)
	fmt.Printf("MerchantClientKey: %s\n", tpr.MerchantClientKey)
	fmt.Printf("PaymentReferenceID: %s\n", tpr.PaymentReferenceID)
	fmt.Printf("PriceID: %d\n", tpr.PriceID)
	fmt.Printf("ServerID: %s\n", tpr.ServerID)
	fmt.Printf("TotalPrice: %s\n", convertToMoneyString(tpr.TotalPrice, tpr.CurrencyCode))
	fmt.Printf("UnitsToSupply: %d\n", tpr.UnitsToSupply)
	fmt.Println()
	fmt.Println()

	fmt.Println("Making payment..")
	fmt.Println()
	fmt.Println()

	paymentResponse, err := wpw.MakePayment(tpr)

	errCheckExit(err)

	fmt.Println("Did payment response:")
	fmt.Printf("ClientID: %s\n", paymentResponse.ClientID)
	fmt.Printf("ServerID: %s\n", paymentResponse.ServerID)
	fmt.Printf("TotalPaid: %s\n", convertToMoneyString(paymentResponse.TotalPaid, tpr.CurrencyCode)) // Payment response doesn't have currency
	fmt.Println("ServiceDeliveryToken:")
	fmt.Printf("\tIssued: %s\n", paymentResponse.ServiceDeliveryToken.Issued)
	fmt.Printf("\tExpiry: %s\n", paymentResponse.ServiceDeliveryToken.Expiry)
	fmt.Printf("\tKey: %s\n", paymentResponse.ServiceDeliveryToken.Key)
	fmt.Printf("\tRefundOnExpiry: %t\n", paymentResponse.ServiceDeliveryToken.RefundOnExpiry)
	fmt.Printf("\tSignature: %s\n", paymentResponse.ServiceDeliveryToken.Signature)
	fmt.Println()
	fmt.Println()

	fmt.Println("Calling BeginServiceDelivery")
	_, err = wpw.BeginServiceDelivery(service.ServiceID, *paymentResponse.ServiceDeliveryToken, numUnits)
	errCheckExit(err)

	fmt.Println()
	fmt.Println()
	fmt.Println("Waiting 5 seconds before calling end service delivery..")
	time.Sleep(time.Second * 5)
	fmt.Println()
	fmt.Println()

	fmt.Println("Calling endServiceDelivery")
	_, err = wpw.EndServiceDelivery(service.ServiceID, *paymentResponse.ServiceDeliveryToken, numUnits)
	errCheckExit(err)

	fmt.Println("Program end, quitting...")
	os.Exit(0)
}

func initLog() {

	log.SetFormatter(&log.TextFormatter{})

	f, err := os.OpenFile("wpwithin.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {

		fmt.Println(err.Error())
	}

	log.SetOutput(f)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	log.Debug("Log is initialised.")
}

func errCheckExit(err error) {

	if err != nil {

		fmt.Println(err.Error())
		os.Exit(1)
	}
}
