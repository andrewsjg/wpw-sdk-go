package hte

import (
	"errors"
	"fmt"

	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
)

// OrderManagerImpl Concrete implementation of order manager.. uses an in memory persistence for orders
type OrderManagerImpl struct {
	orders map[string]types.Order
}

// NewOrderManager Create a new instance of OrderManager
func NewOrderManager() (OrderManager, error) {

	result := &OrderManagerImpl{}
	result.orders = make(map[string]types.Order, 0)

	return result, nil
}

// AddOrder add an order
func (om *OrderManagerImpl) AddOrder(order types.Order) error {

	if _, ok := om.orders[order.UUID]; ok {

		return errors.New("Order already exists")
	}

	om.orders[order.UUID] = order

	return nil
}

// GetOrder get an order from the manager by searching for its payment reference
func (om *OrderManagerImpl) GetOrder(paymentReference string) (*types.Order, error) {

	if order, ok := om.orders[paymentReference]; ok {

		return &order, nil
	}

	return nil, errors.New("Order not found")
}

// OrderExists check if an order exists by searching for its payment reference
func (om *OrderManagerImpl) OrderExists(paymentReference string) bool {

	if _, found := om.orders[paymentReference]; found {

		return true
	}

	return false

}

// UpdateOrder update an order in the manager.
func (om *OrderManagerImpl) UpdateOrder(order types.Order) error {

	if om.OrderExists(order.UUID) {

		om.orders[order.UUID] = order

		return nil
	}

	return fmt.Errorf("Cannot update order %s as it does not exist", order.UUID)
}
