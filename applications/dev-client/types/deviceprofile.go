package devclienttypes

import (
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
)

type DeviceProfile struct {
	DeviceEntity *DeviceEntity `json:"device"`
}

type DeviceEntity struct {
	//Uid         string    `json:"uid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Producer    *Producer `json:"producer"`
	Consumer    *Consumer `json:"consumer"`
}

type Producer struct {
	Services       []*ServiceProfile `json:"services"`
	ProducerConfig *ProducerConfig   `json:"config"`
}

type ServiceProfile struct {
	Id          int             `json:"serviceID"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Prices      []*PriceProfile `json:"prices"`
}

type PriceProfile struct {
	ID              int                 `json:"priceID"`
	Description     string              `json:"priceDescription"`
	PricePerUnit    *types.PricePerUnit `json:"pricePerUnit"`
	UnitID          int                 `json:"unitID"`
	UnitDescription string              `json:"unitDescription"`
}

type Consumer struct {
	ConsumerConfig *ConsumerConfig `json:"config"`
	HCECard        *types.HCECard  `json:"hceCard"`
	AutoConsume    *AutoConsume    `json:"autoConsume"`
}

type AutoConsume struct {
	DeviceUid     string `json:"deviceUid"`
	ServiceID     int    `json:"serviceID"`
	MaximumAmount int    `json:"maximumAmount"`
	//CurrencyCode  string `json:"currencyCode"`
	UnitID int `json:"unitID"`
}

type ProducerConfig struct {
	//BroadcastTimeout      int    `json:"broadcastTimeout"`
	PspMerchantServiceKey string `json:"pspMerchantServiceKey"`
	PspMerchantClientKey  string `json:"pspMerchantClientKey"`
	// No way to set these yet, see https://github.com/wptechinnovation/worldpay-within-sdk/issues/18
	//HTEPort               int    `json:"htePort"`
	//HTEURLPrefix          string `json:"hteURLPrefix"`
	//HTEScheme             string `json:"hteScheme"`
}

type ConsumerConfig struct {
	DeviceDiscoveryTimeout int `json:"deviceDiscoveryTimeout"`
}
