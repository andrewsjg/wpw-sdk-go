package hte

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types/event"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/utils"
)

// ServiceHandler Coordinate requests between RPC interface and internal SDK interface
type ServiceHandler struct {
	device       *types.Device
	psp          psp.PSP
	credential   *Credential
	orderManager OrderManager
	eventHandler event.Handler
}

// NewServiceHandler Create a new Service Handler
func NewServiceHandler(device *types.Device, psp psp.PSP, credential *Credential, orderManager OrderManager, eventHandler event.Handler) *ServiceHandler {

	result := &ServiceHandler{
		device:       device,
		psp:          psp,
		credential:   credential,
		orderManager: orderManager,
		eventHandler: eventHandler,
	}

	return result
}

// ServiceDiscovery List all the services available on the current device
func (srv *ServiceHandler) ServiceDiscovery(w http.ResponseWriter, r *http.Request) {

	// GET

	defer func() {
		if err := recover(); err != nil {

			returnMessage(w, http.StatusInternalServerError, err)
		}
	}()

	var responseServices []types.ServiceDetails

	for _, srv := range srv.device.Services {

		responseServices = append(responseServices, types.ServiceDetails{
			ServiceID:          srv.ID,
			ServiceDescription: srv.Description,
			ServiceName:        srv.Name,
		})
	}

	response := types.ServiceListResponse{
		Services: responseServices,
		ServerID: srv.device.UID,
	}

	returnMessage(w, http.StatusOK, response)
}

// ServicePrices List all the price variants for a specified service
func (srv *ServiceHandler) ServicePrices(w http.ResponseWriter, r *http.Request) {

	// GET

	defer func() {
		if err := recover(); err != nil {

			returnMessage(w, http.StatusInternalServerError, err)
		}
	}()

	// Parse variables from request
	reqVars := mux.Vars(r)
	svcID, err := strconv.Atoi(reqVars["service_id"])

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse input service id",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if svc, ok := srv.device.Services[svcID]; ok {

		response := types.ServicePriceResponse{}
		response.ServerID = srv.device.UID

		for _, price := range svc.Prices {

			response.Prices = append(response.Prices, price)
		}

		returnMessage(w, http.StatusOK, response)

	} else {

		errorResponse := types.ErrorResponse{
			Message: fmt.Sprintf("Service not found for id %d", svcID),
		}

		returnMessage(w, http.StatusNotFound, errorResponse)
	}
}

// ServiceTotalPrice Get the total price for a current service selection
func (srv *ServiceHandler) ServiceTotalPrice(w http.ResponseWriter, r *http.Request) {

	// POST

	defer func() {
		if err := recover(); err != nil {

			returnMessage(w, http.StatusInternalServerError, err)
		}
	}()

	// Parse variables from URI
	reqVars := mux.Vars(r)
	svcID, err := strconv.Atoi(reqVars["service_id"])

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse input service id",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Parse message body (POST)
	var totalPriceRequest types.TotalPriceRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to read POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := r.Body.Close(); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to close POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := json.Unmarshal(body, &totalPriceRequest); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse POST body",
		}

		returnMessage(w, 422 /*Unprocessable Entity*/, errorResponse)
		return
	}

	if svc, ok := srv.device.Services[svcID]; ok {

		if price, ok := svc.Prices[totalPriceRequest.SelectedPriceID]; ok {

			response := types.TotalPriceResponse{}
			response.ServerID = srv.device.UID
			response.ClientID = totalPriceRequest.ClientID
			response.PriceID = totalPriceRequest.SelectedPriceID
			response.UnitsToSupply = totalPriceRequest.SelectedNumberOfUnits
			response.TotalPrice = price.PricePerUnit.Amount * totalPriceRequest.SelectedNumberOfUnits
			response.CurrencyCode = price.PricePerUnit.CurrencyCode
			response.MerchantClientKey = srv.credential.MerchantClientKey

			orderUUID, err := utils.NewUUID()

			if err != nil {

				errorResponse := types.ErrorResponse{
					Message: "Internal error [generate order UUID]",
				}

				returnMessage(w, http.StatusInternalServerError, errorResponse)
			}
			response.PaymentReferenceID = orderUUID

			order := types.Order{
				UUID:                  orderUUID,
				ClientID:              response.ClientID,
				SelectedNumberOfUnits: response.UnitsToSupply,
				SelectedPriceID:       response.PriceID,
				Service:               *svc,
			}

			err = srv.orderManager.AddOrder(order)

			if err != nil {

				errorResponse := types.ErrorResponse{
					Message: "Unable to add order to local store",
				}

				returnMessage(w, http.StatusInternalServerError, errorResponse)

			} else {

				returnMessage(w, http.StatusOK, response)
			}

		} else {

			errorResponse := types.ErrorResponse{
				Message: fmt.Sprintf("Price not found for id %d", totalPriceRequest.SelectedPriceID),
			}

			returnMessage(w, http.StatusNotFound, errorResponse)
		}

	} else {

		errorResponse := types.ErrorResponse{
			Message: fmt.Sprintf("Service not found for id %d", svcID),
		}

		returnMessage(w, http.StatusNotFound, errorResponse)
	}
}

// Payment Make a payment for a service
func (srv *ServiceHandler) Payment(w http.ResponseWriter, r *http.Request) {

	// POST

	defer func() {
		if err := recover(); err != nil {

			returnMessage(w, http.StatusInternalServerError, err)
		}
	}()

	// Parse message body (POST)
	var paymentRequest types.PaymentRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to read POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := r.Body.Close(); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to close POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := json.Unmarshal(body, &paymentRequest); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse POST body",
		}

		returnMessage(w, 422 /*HTTP Status Code: Unprocessable Entity*/, errorResponse)
		return
	}

	_order, err := srv.orderManager.GetOrder(paymentRequest.PaymentReferenceID)

	orderCurrency := srv.device.Services[_order.Service.ID].Prices[_order.SelectedPriceID].PricePerUnit.CurrencyCode

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: fmt.Sprintf("Unable to find order for payment ref %s", paymentRequest.PaymentReferenceID),
		}

		returnMessage(w, http.StatusNotFound, errorResponse)
	} else {

		// Some quick checks to compare validity of incoming request
		if strings.Compare(_order.ClientID, paymentRequest.ClientID) != 0 {

			errorResponse := types.ErrorResponse{
				Message: "Client ID does not match Order Client ID",
			}

			returnMessage(w, http.StatusBadRequest, errorResponse)
		} else {

			totalPrice := _order.Service.Prices[_order.SelectedPriceID].PricePerUnit.Amount * _order.SelectedNumberOfUnits
			orderPPU := _order.Service.Prices[_order.SelectedPriceID].PricePerUnit

			// Convert order price per unit from minor units (e.g. Pence) to major units (e.g. Pounds) (e.g. 100p -> £1.00)
			strUnitPrice := utils.DoUnitConvertFormat(orderPPU.Amount, 2, fmt.Sprintf("%s <amount>", orderPPU.CurrencyCode))

			orderDescription := fmt.Sprintf("%s - %d units @ %s per unit.", _order.Service.Name, _order.SelectedNumberOfUnits, strUnitPrice)
			pspReference, err := srv.psp.MakePayment(totalPrice, orderCurrency, paymentRequest.ClientToken, orderDescription, _order.UUID)

			if err != nil {

				log.WithFields(log.Fields{"Error": err, "totalPrice": totalPrice, "orderCurrency": orderCurrency,
					"paymentRequest.ClientToken": paymentRequest.ClientToken, "Order descripion": orderDescription,
					"Order UUID": _order.UUID}).Error("Unable to process payment with gateway at this time")

				errorResponse := types.ErrorResponse{
					Message: "Unable to process payment with gateway at this time",
				}

				returnMessage(w, http.StatusInternalServerError, errorResponse)

			} else {

				deliveryToken := &types.ServiceDeliveryToken{
					Key:            _order.UUID, /* UUID is a hack, allows us to find order later :/ */ // TODO cryptographically secure, random key generation.
					Issued:         time.Now(),
					Expiry:         time.Now().Add(time.Duration(168 * time.Hour)),
					RefundOnExpiry: false, // TODO Map this into the Service Struct
					Signature:      nil,   // TODO implement HMAC generation scheme.
				}

				paymentResponse := types.PaymentResponse{
					ClientID:             _order.ClientID,
					ServerID:             srv.device.UID,
					TotalPaid:            totalPrice,
					ServiceDeliveryToken: deliveryToken,
				}

				_order.PSPReference = pspReference
				_order.DeliveryToken = *deliveryToken

				if err := srv.orderManager.UpdateOrder(*_order); err != nil {

					errorResponse := types.ErrorResponse{
						Message: "Unable to update order internally",
					}

					returnMessage(w, http.StatusInternalServerError, errorResponse)
					return
				}

				returnMessage(w, http.StatusOK, paymentResponse)
			}
		}
	}
}

// ServiceDeliveryBegin Begin delivery of a purchased service
func (srv *ServiceHandler) ServiceDeliveryBegin(w http.ResponseWriter, r *http.Request) {

	// POST

	log.Debug("begin hte.ServiceHandlerImpl.ServiceDeliveryBegin()")

	defer log.Debug("end hte.ServiceHandlerImpl.ServiceDeliveryBegin()")

	defer func() {
		if err := recover(); err != nil {

			returnMessage(w, http.StatusInternalServerError, err)
		}
	}()

	// Parse variables from request
	reqVars := mux.Vars(r)
	svcID, err := strconv.Atoi(reqVars["service_id"])

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse input service id",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Parse message body (POST)
	var deliveryRequest types.BeginServiceDeliveryRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to read POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := r.Body.Close(); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to close POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := json.Unmarshal(body, &deliveryRequest); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse POST body",
		}

		returnMessage(w, 422 /*Unprocessable Entity*/, errorResponse)
		return
	}

	// Validate the delivery request units against what is currenly available
	strValidation, err := validateDeliveryToken(deliveryRequest, srv.orderManager)

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Error validating ServiceDeliveryToken",
		}

		returnMessage(w, 500, errorResponse)
		return
	} else if !strings.EqualFold(strValidation, "") {

		errorResponse := types.ErrorResponse{
			Message: strValidation,
		}

		returnMessage(w, 400, errorResponse)
		return
	}

	if err != nil {

		log.Errorf("Error retrieving order for delivery token: %s", err.Error())

		errorResponse := types.ErrorResponse{
			Message: fmt.Sprintf("Error retrieving order using delivery token.."),
		}

		returnMessage(w, 400, errorResponse)
		return
	}

	// Order has been found for delivery token - ensure the service id matches the service id passed in URL.
	_order, err := srv.orderManager.GetOrder(deliveryRequest.ServiceDeliveryToken.Key)

	if svcID != _order.Service.ID {

		log.Info("Requested service does not match the service id presented in request ServiceDeliveryToken")

		errorResponse := types.ErrorResponse{
			Message: fmt.Sprintf("Could not deliver service. Please contact service provider."),
		}

		returnMessage(w, 404, errorResponse)
		return
	}

	// Need to validate that the number of units requested is legal
	// i.e. It does not exceed the number of units in the order or the number already consumed by consumer

	if _, ok := srv.device.Services[svcID]; ok {

		response := types.BeginServiceDeliveryResponse{}
		response.ClientID = deliveryRequest.ClientID
		response.ServiceDeliveryToken = deliveryRequest.ServiceDeliveryToken
		response.ServerID = srv.device.UID
		response.UnitsToSupply = deliveryRequest.UnitsToSupply

		if srv.eventHandler != nil {

			log.Debug("Core event handler is set, calling event in core EventHandler")

			go srv.eventHandler.BeginServiceDelivery(svcID, deliveryRequest.ServiceDeliveryToken, deliveryRequest.UnitsToSupply)
		}

		returnMessage(w, http.StatusOK, response)

	} else {

		errorResponse := types.ErrorResponse{
			Message: fmt.Sprintf("Service not found for id %d", svcID),
		}

		returnMessage(w, http.StatusNotFound, errorResponse)
	}
}

// ServiceDeliveryEnd End delivery of a purchased service
func (srv *ServiceHandler) ServiceDeliveryEnd(w http.ResponseWriter, r *http.Request) {

	// POST

	log.Debug("begin hte.ServiceHandlerImpl.ServiceDeliveryEnd()")

	defer log.Debug("end hte.ServiceHandlerImpl.ServiceDeliveryEnd()")

	defer func() {
		if err := recover(); err != nil {

			returnMessage(w, http.StatusInternalServerError, err)
		}
	}()

	// Parse variables from request
	reqVars := mux.Vars(r)
	svcID, err := strconv.Atoi(reqVars["service_id"])

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse input service id",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	// Parse message body (POST)
	var deliveryRequest types.EndServiceDeliveryRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to read POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := r.Body.Close(); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to close POST body",
		}

		returnMessage(w, http.StatusBadRequest, errorResponse)
		return
	}

	if err := json.Unmarshal(body, &deliveryRequest); err != nil {

		errorResponse := types.ErrorResponse{
			Message: "Unable to parse POST body",
		}

		returnMessage(w, 422 /*Unprocessable Entity*/, errorResponse)
		return
	}

	var _order types.Order
	if srv.orderManager.OrderExists(deliveryRequest.ServiceDeliveryToken.Key) {

		tmp, err := srv.orderManager.GetOrder(deliveryRequest.ServiceDeliveryToken.Key)

		if err != nil {

			errorResponse := types.ErrorResponse{
				Message: "Unable to retrieve order internally",
			}

			returnMessage(w, 500, errorResponse)
			return
		}
		_order = *tmp

	} else {

		errorResponse := types.ErrorResponse{
			Message: "Unable to find order for ServiceDeliveryToken",
		}

		returnMessage(w, 404, errorResponse)
		return
	}

	if _, ok := srv.device.Services[svcID]; ok {

		response := types.EndServiceDeliveryResponse{}
		response.ClientID = deliveryRequest.ClientID
		response.ServiceDeliveryToken = deliveryRequest.ServiceDeliveryToken
		response.ServerID = srv.device.UID
		response.UnitsJustSupplied = deliveryRequest.UnitsReceived

		// TODO NASTY!! - We trust what clients are telling us, need to validate with producer application!
		_order.ConsumedUnits = deliveryRequest.UnitsReceived + _order.ConsumedUnits
		remain := _order.SelectedNumberOfUnits - _order.ConsumedUnits
		response.UnitsRemaining = remain

		if srv.eventHandler != nil {

			log.Debug("Core event handler is set, calling event in core EventHandler")

			go srv.eventHandler.EndServiceDelivery(svcID, deliveryRequest.ServiceDeliveryToken, deliveryRequest.UnitsReceived)
		}

		returnMessage(w, http.StatusOK, response)

	} else {

		errorResponse := types.ErrorResponse{
			Message: fmt.Sprintf("Service not found for id %d", svcID),
		}

		returnMessage(w, http.StatusNotFound, errorResponse)
	}
}

// Helper function for returning HTTP responses
func returnMessage(w http.ResponseWriter, statusCode int, message interface{}) {

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(message); err != nil {

		panic(err)
	}
}

// validateDeliveryToken will validate the various parameters of a ServiceDeliveryToken
// If the returned string is empty then validation passed, if not then the contents of the string
// is the reason for validation failure
func validateDeliveryToken(request types.BeginServiceDeliveryRequest, orderManager OrderManager) (string, error) {

	sdt := request.ServiceDeliveryToken

	if !orderManager.OrderExists(sdt.Key) {

		return "Order not found for ServiceDeliveryToken", nil
	}

	_order, err := orderManager.GetOrder(sdt.Key)

	if err != nil {

		return "", err
	}

	if strings.EqualFold("", sdt.Key) {

		return "ServiceDeliveryToken key is empty", nil
	}

	if !sdt.Expiry.After(time.Now()) {

		return "ServiceDeliveryToken has expired", nil
	}

	log.Debugf("ORDER:: %+v", _order)
	log.Debugf("REQ SDT Key: %s", sdt.Key)

	if !strings.EqualFold(_order.DeliveryToken.Key, sdt.Key) {

		return "Invalid ServiceDeliveryToken key", nil
	}

	unitsAvailable := _order.SelectedNumberOfUnits - _order.ConsumedUnits
	if request.UnitsToSupply > unitsAvailable {

		return fmt.Sprintf("Requested units (%d) not available for selected order. Units available = %d", request.UnitsToSupply, unitsAvailable), nil
	}

	return "", nil
}
