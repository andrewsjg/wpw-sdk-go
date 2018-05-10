package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/psp/onlineworldpay/types"
	"github.com/gorilla/mux"
	"github.com/nanobox-io/golang-scribble"
	"github.com/nu7hatch/gouuid"
)

var trpmMap map[string]types.TokenResponsePaymentMethod
var db *scribble.Driver
var orderInformation OrderInformation

const DB_NAME = "orders"

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./static/js/", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	db, _ = scribble.New(".", nil)
	cleandb()
	port := ":8080"
	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(dir))))
	router.HandleFunc("/v1/tokens", Tokens)
	router.HandleFunc("/v1/orders", Orders)
	router.HandleFunc("/v1/transactions", Transactions)
	router.HandleFunc("/api/", WTRoot)
	router.HandleFunc("/api/PreVault/Card", WTPrevaultCard)
	router.HandleFunc("/", HomePage)
	http.Handle("/", router)
	fmt.Println("Serving worldpay web service mock on port " + port)
	trpmMap = make(map[string]types.TokenResponsePaymentMethod)
	log.Fatal(http.ListenAndServe(port, router))
}
func WTRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/api/ request received from " + r.RemoteAddr)
	defer r.Body.Close()

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(requestBody)
}
func WTPrevaultCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/api/PreVault/Card request received from " + r.RemoteAddr)
	defer r.Body.Close()

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(requestBody)
	fmt.Println(r.Body)
	fmt.Println("%v\n",r)
}

func Tokens(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/v1/tokens purchase request received from " + r.RemoteAddr)

	defer r.Body.Close()

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	tokenRequest := types.TokenRequest{}
	json.Unmarshal(requestBody, &tokenRequest)
	unmaskedCardPart := regexp.MustCompile("[0-9]{4}$")
	//fmt.Println("tokenRequest %v", tokenRequest)
	//fmt.Println("---")
	fmt.Println("tokenRequest.PaymentMethod %v", tokenRequest.PaymentMethod)
	fmt.Println("---")
	//fmt.Println("r %v", r)
	//fmt.Println("---")
	//fmt.Println("requestBody %v", requestBody)
	//fmt.Println("---")
	trpm := types.TokenResponsePaymentMethod{
		Type:                              tokenRequest.PaymentMethod.Type,
		Name:                              tokenRequest.PaymentMethod.Name,                                                              // TODO
		ExpiryMonth:                       tokenRequest.PaymentMethod.ExpiryMonth,                                                       // TODO
		ExpiryYear:                        tokenRequest.PaymentMethod.ExpiryYear,                                                        // TODO
		CardType:                          tokenRequest.PaymentMethod.Type,                                                              // TODO
		MaskedCardNumber:                  strings.Repeat("*", 12) + unmaskedCardPart.FindString(tokenRequest.PaymentMethod.CardNumber), // TODO
		CardSchemeType:                    "",
		CardSchemeName:                    "",
		CardIssuer:                        "",
		CountryCode:                       "",
		CardClass:                         "",
		CardProductTypeDescNonContactless: "",
		CardProductTypeDescContactless:    "",
		Prepaid: "",
	}
	uuid, _ := uuid.NewV4()
	tokenResponse := types.TokenResponse{
		Reusable: false,
		Token:    uuid.String(), // TODO
		TokenResponsePaymentMethod: trpm,
	}
	trpmMap[uuid.String()] = trpm
	json.NewEncoder(w).Encode(tokenResponse)
	fmt.Println("")
}
func Orders(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/v1/orders request received from " + r.RemoteAddr)

	defer r.Body.Close()

	respBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println("BODY:")
	fmt.Println(string(respBody))
	orderRequest := types.OrderRequest{}
	json.Unmarshal(respBody, &orderRequest)

	fmt.Println("%v", orderRequest)

	orpr := types.OrderResponsePaymentResponse{
		Type:                              trpmMap[orderRequest.Token].Type,
		Name:                              trpmMap[orderRequest.Token].Name,
		ExpiryMonth:                       trpmMap[orderRequest.Token].ExpiryMonth,
		ExpiryYear:                        trpmMap[orderRequest.Token].ExpiryYear,
		CardType:                          trpmMap[orderRequest.Token].CardType,
		MaskedCardNumber:                  trpmMap[orderRequest.Token].MaskedCardNumber,
		CardSchemeType:                    trpmMap[orderRequest.Token].CardSchemeType,
		CardSchemeName:                    trpmMap[orderRequest.Token].CardSchemeName,
		CardIssuer:                        trpmMap[orderRequest.Token].CardIssuer,
		CountryCode:                       trpmMap[orderRequest.Token].CountryCode,
		CardClass:                         trpmMap[orderRequest.Token].CardClass,
		CardProductTypeDescNonContactless: trpmMap[orderRequest.Token].CardProductTypeDescNonContactless,
		CardProductTypeDescContactless:    trpmMap[orderRequest.Token].CardProductTypeDescContactless,
		Prepaid: trpmMap[orderRequest.Token].Prepaid, //IssueNumber int32 ,StartMonth int32 , StartYear int32
	}

	delete(trpmMap, orderRequest.Token)
	orrs := types.OrderResponseRiskScore{
		Value: "1",
	}
	t := types.OrderResponse{
		OrderCode:         "Test_code",                   // TODO
		Token:             orderRequest.Token,            // TODO
		OrderDescription:  orderRequest.OrderDescription, // TODO
		Amount:            int32(orderRequest.Amount),    // TODO
		CurrencyCode:      orderRequest.CurrencyCode,     // TODO
		PaymentStatus:     "SUCCESS",
		PaymentResponse:   orpr,
		CustomerOrderCode: orderRequest.CustomerOrderCode, // TODO
		Environment:       "TEST",
		RiskScore:         orrs, // ResultCodes not supported by current WPW API
	}

	orderInformation.MaskedCard = orpr.MaskedCardNumber
	orderInformation.TransactionDateAndTime = time.Now()
	orderInformation.TotalPrice = orderRequest.Amount
	orderInformation.CurrencyCode = orderRequest.CurrencyCode
	orderInformation.OrderDescription = orderRequest.OrderDescription

	_, err := json.Marshal(orderInformation)
	if err != nil {
		fmt.Println(err)
	}
	if err := db.Write(DB_NAME, orderRequest.CustomerOrderCode, orderInformation); err != nil {
		fmt.Println("Error", err)
	}
	json.NewEncoder(w).Encode(t)
}

func Transactions(w http.ResponseWriter, r *http.Request) {
	records, err := db.ReadAll(DB_NAME)
	if err != nil {
		fmt.Println("DB readAll error: ", err)
	}
	transactions := []OrderInformation{}

	for _, f := range records {
		transactionFound := OrderInformation{}
		if err := json.Unmarshal([]byte(f), &transactionFound); err != nil {
			fmt.Println("Error", err)
		}
		transactions = append(transactions, transactionFound)
	}

	response, err := json.Marshal(transactions)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(response)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func cleandb() {
	db.Delete(DB_NAME, "")
}
