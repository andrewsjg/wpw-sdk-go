{
	"device": {
		"name": "RoboCar Station",
		"description": "Car services offered by robot",
		"producer": {
			"services": [{
				"serviceID": 1,
				"name": "RoboAir",
				"description": "Car tyre pressure checked and topped up by robot",
				"prices": [{
						"priceID": 1,
						"priceDescription": "Measure and adjust pressure",
						"pricePerUnit": {
							"amount": 25,
							"currencyCode": "GBP"
						},
						"unitID": 1,
						"unitDescription": "Tyre"
					}, {
						"priceID": 2,
						"priceDescription": "Measure and adjust pressure - four tyres for the price of three",
						"pricePerUnit": {
							"amount": 75,
							"currencyCode": "GBP"
						},
						"unitID": 2,
						"unitDescription": "4 Tyre"
					}

				]
			}, {
				"serviceID": 2,
				"name": "RoboWash",
				"description": "Car washed by robot",
				"prices": [{
						"priceID": 1,
						"priceDescription": "Car wash",
						"pricePerUnit": {
							"amount": 500,
							"currencyCode": "GBP"
						},
						"unitID": 1,
						"unitDescription": "Single wash"
					}, {
						"priceID": 2,
						"priceDescription": "SUV Wash",
						"pricePerUnit": {
							"amount": 650,
							"currencyCode": "GBP"
						},
						"unitID": 1,
						"unitDescription": "Single wash"
					}

				]
			}],
			"config": {
				"pspMerchantServiceKey": "99999",
				"pspMerchantClientKey": "8888"
			}
		},
		"consumer": null,
		"config": null
	}
}
