# Worldpay Within (SecureNet integration)

This example project contains both a consumer and producer configured to use the SecureNet payment service for payments.

## Usage

* Inside each directory { consumer, producer } run `go get ./...` to install required dependencies
* Inside each directory { consumer, producer } run `go build` to build each component
* Run `./producer` from inside the producer directory - This will start a producer and start broadcasting
* Run `./consumer` from inside the consumer directory - This will setup a consumer and start scanning for producer devices. Once a device is discovered it will run through service and price discovery etc, picking the first item returned in each list. Once a service has been selected a payment is made to the SecureNet service.

**Note:** You may want to examine the code in `producer/main.go` that populates the `pspConfig` variable. For the demo to work you will most likely have to populate the details as found in your SecureNet developer account.
