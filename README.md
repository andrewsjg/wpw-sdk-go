# Worldpay Within - The Go SDK

The Go implementation for the Worldpay Within IoT payment SDK. This SDK, or Software Development Kit, enables smart devices to discover each other, negogiate a price for services, make a payment (through the Worldpay payments gateway) for services, and then consume services via a 'trusted trigger'. For more information see our documentation website here: http://www.worldpaywithin.com

![The Worldpay Within puzzle piece](http://wptechinnovation.github.io/worldpay-within-sdk/images/architecture/worldpayWithinFig1.png)

## Get started
1. Download this repo
3. Run the examples...

## Run the examples
* build examles using 'go build'
* Run the consumer project
* Simultaneously run the producer or producer-callbacks project
* The two smart devices should communicate with each other and make a payment

## Compatibility and pre-requisites
* Assumption here is you have go already installed and setup okay

## See the payments:
1. Sign up to https://online.worldpay.com if you haven't already done so
2. Got to settings > API keys and get your test keys
3. Replace the keys in the consumer and producer java example source files
4. Re-run the examples and you should see the payments coming through on the WPOP (Worldpay Online) payments dashboard
  
## So what does it do:

![The Worldpay Within Flows sequence diagram](http://wptechinnovation.github.io/worldpay-within-sdk/images/architecture/serviceOverview.png)

You can see there are four phases; discover, negotiation, payment and then service delivery, for more information visit our website at http://www.worldpaywithin.com.

[The flows and API can be found here](http://wptechinnovation.github.io/worldpay-within-sdk/the-flows.html)

## Want to contribute:

Want to contribute, then please clone the repo and create a branch, once you've made your changes create a pull request, we will review your code, and if accepted it will be merged into the code base. You can also raise issues here on github, or contact us direct at innovation@worldpay.com or alternatively join our slack channel at iotpay.slack.com.

Go is slightly different, as it is what the innards of the 'IoT Core Componet' or 'rpc-agent' is written in, continue to learn more...

# SDK Core package

This is the source for the core of Worldpay Within and is separated into the following paths:

* `wpwithin` - Worldpay Within SDK core implementation
* `wpwithin_test` - contains tests to various parts of the SDK (under development)
* `examples` - Some sample code showing how to develop a producer and consumer

## wpwithin

Implementation of the [Worldpay Within architecture](http://wptechinnovation.github.io/worldpay-within-sdk/architecture.html).

### configuration

Used to load configuration files

### core

Contains a `core` structure for holding state within the SDK along with a factory for creating SDK components such as the HTE Service, RPC Layer, Device broadcast/discover etc. The core holds references to all the critical objects. Or to be more precise, the SDK Core acts as a container for dependencies of the SDK.

### hte

Host Terminal Emulation - Point of interaction between consumers and producers. As per the architecture, HTE exposes a REST HTTP service allowing consumers to discover services and make payments.

Service - (via service, servicehandler and serviceimpl) A REST service, exposing an interact for consumers
Client - A HTTP rest client to interact with the HTE service
OrderManager - Manages state of orders and processed payments

The HTE client allows interaction with the HTE service. There is a credential store for the 'terminal' or payments gateway credentials. The Order Manager coordinates during negotiation, payment and delivery flows. There are also some help http request object(s).

* [More detail about the flows and associate diagrams can be found in the detailed documentation here](http://wptechinnovation.github.io/worldpay-within-sdk/architecture.html)

### psp

This code enables communication with the online.worldpay.com payments gatesway to make payments. The results of these payment can be viewed by having an associated account (associated with the credentials) on online.worldpay.com.

### rpc

Implementation of a Thrift server to allow non Go languages call into the SDK. Ad-hoc callbacks from the SDK are also supported.

### servicediscovery

This contains the broadcaster, scanner and communicator. The broadcaster allows for a devices presence to be seen on a network while the scanner is used to detect the broadcast messages.

Communicator contains various functions for abstracting communication, with UDP Broadcast currently supported.

### types

Data types to model the Worldpay Within architecture.

### utils

Various utilities for networking, text processing, UUID processing, etc.
