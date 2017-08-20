# Dev Client
The dev-client example console app gives some examples on how to use the sdk and ultimately how to enact a payment between a producer and a consumer. The app runs as a console app, and you select items by entering the item number:

```
----------------------------- Worldpay Within SDK Client ----------------------------
0 - -------------------- GENERAL  --------------------
1 - Init new device
2 - Start RPC Service
3 - Get device info
4 - Load device profile
5 - Reset session
6 - -------------------- PRODUCER --------------------
7 - Init new producer
8 - Add RoboWash service
9 - Add RoboAir service
10 - Start service broadcast
11 - Stop broadcast
12 - -------------------- CONSUMER --------------------
13 - Prepare new consumer
14 - Scan services
15 - Auto consume from profile info
16 - --------------------------------------------------
17 - Exit
-------------------------------------------------------------------------------------
Please select choice:
```

There are example profiles in the root of the dev-client application directory: test-consumer.profile and test-producer.profile.
To load a profile:
* Enter 4 followed by the name of the profile you wish to load, for example 'test-producer.profile'.
* Then press 'y' to continue
* Enter 10 to start the service broadcast
* Enter a timeout in milliseconds for the broadcast e.g 60000.

Then in another terminal (or on another device connected to the same wifi network):
* Enter 4 to load another profile, for example 'test-consumer.profile'.
* Enter 15 to auto consume from the profile info - after a few seconds a payment should now occur based on the selection criteria in the "autoConsume" section of the device profile.

As shown below:

```
Please select choice: 15
Starting auto consume...
Found Service:: (192.168.43.143:8080/) - Car services offered by robot
Selecting service: 1 - Car tyre pressure checked and topped up by robot
Selecting price: (2) Measure and adjust pressure - four tyres for the price of three @ 75GBP, 4 Tyre (Unit id = 2)
Payment of 75 made successfully
```

This happens because the "autoConsume" section of 'test-consumer.profile' has the following data:

```
"autoConsume": {
	"deviceUid": "3b0f50f5-c4e6-4d0f-5cbf-52466e858ea3",
	"serviceID": 1,
	"unitID": 2
}
```
This means that the device with Uid "3b0f50f5-c4e6-4d0f-5cbf-52466e858ea3" will be selected. (The uuid can be found in a text file 'uuid.txt' that gets written to the device upon initialisation, you can use the dev-client app to initialise the device, or simply place a 'uuid.txt' file with just the uuid in, at the same location as the dev-client app).
When the device Uid has been selected, it then selects the service with serviceID 1, and finally it will select the unit with unitID 2.

## Detailed steps on setting up and running the dev client app

### Setting up the app to run

1. copy in the test-consumer.profile and the test-producer.profile
2. Run the dev client and choose option 1, to init new device, this will generate the uuid. The uuid is of the device, and the service uses it as it's id, so configuring this, allows the particular service to be auto consumed.
3. This uuid needs to be copied and pasted into the Test consumer profile document, under the autoconsume field
4. Open two terminal windows
5. Run the dev client in each
6. In one choose 4. Load device profile and type 'test-consumer.profile'
7. In the other choose 4. Load device profile and type 'test-producer.profile'
8. In the producer window choose 10 to start broadcast and put in timeout e.g. 68000
9. In the consumer window choose 15 and it will go through a test flow e.g.

````
Please select choice: 15
Starting auto consume...
Found Service:: (<ip.address.removed.for.illustration.purposes>:8080/) - Car services offered by robot
Selecting service: 1 - Car tyre pressure checked and topped up by robot
Selecting price: (2) Measure and adjust pressure - four tyres for the price of three @ 75GBP, 4 Tyre (Unit id = 2)
Payment of 75 made successfully
Continue (y/n): y
```

### Manually broadcasting and scanning for services
1. Setup as above the producer and consumer
2. For the producer choose '10' to start the broadcast e.g. for 68000
3. For the consumer choose '14' scan services for e.g. 58000
4. Nothing at this point will be output to the console if successful however if you tail -f output.log, it will show you that the consumer is successfully reading in the messages that are broacast


### Running the app as producer, and running the consumer in one of the example apps e.g. Java
1. Setup the producer app as above
2. Setup the consumer app in Java - this defaults to 9090
3. So start an instance of the rpc-agent on 9090
4. Run the consumer app in Java
5. Choose option 10 of in the producer window, and that will start the app broadcasting.
6. The Java example app, will then consume one of the services at a arbitrarily chosen price point, and make the payment. For example;

```
Starting Consumer Example Written in Java.
1 services found:
Device Description: Car services offered by robot
Hostname: <removed-for-illustration>
Port: 64521
URL Prefix:
ServerId: <removed-for-illustration>
--------
2 services found
Service:
Id: 2
Description: Car washed by robot
------
Service:
Id: 1
Description: Car tyre pressure checked and topped up by robot
------
2 prices found for service id 2
Price:
Id: 2
Description: SUV Wash
UnitId: 1
UnitDescription: Single wash
Unit Price Amount: 650
Unit Price CurrencyCode: GBP
------
Price:
Id: 1
Description: Car wash
UnitId: 1
UnitDescription: Single wash
Unit Price Amount: 500
Unit Price CurrencyCode: GBP
------
Did retrieve price quote:
Merchant client key: <removed-for-illustration>
Payment reference id: <removed-for-illustration>
Units to supply: 1
Total price: 0
Payment response: Client UUID: <removed-for-illustration>
Client ServiceId: <removed-for-illustration>
Total paid: 650
ServiceDeliveryToken.issued: 2016-08-16T13:32:03Z
ServiceDeliveryToken.expiry: 2016-08-17T13:32:03Z
ServiceDeliveryToken.key: <removed-for-illustration>
ServiceDeliveryToken.signature: <removed-for-illustration>
ServiceDeliveryToken.refundOnExpiry: false
````


### Running the app as consumer, and running the producer in one of the example apps e.g. Java
1.Setup the consumer app as above
2. Setup the producer app in Java - this defaults to 9091
3. So start an instance of the rpc-agent on 9091
4. Run the producer app in Java
5. Choose option 14 of the consumer window, and that will start the app scanning for services.
6. The consumer app, will find the services for the producer app service message that is broadcast. For example;
