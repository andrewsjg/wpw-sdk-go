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
