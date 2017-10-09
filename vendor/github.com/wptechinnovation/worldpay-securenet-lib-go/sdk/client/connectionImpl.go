package client

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk"
	"github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/utils"

	"github.com/WPTechInnovation/worldpay-securenet-lib-go/sdk/service"
)

// ConnectionImpl Required details to make a connection to the API server
type ConnectionImpl struct {
	APIBaseURL  string
	AppVersion  string
	SecureNetID string
	SecureKey   string
	Timeout     int
	SSLCheck    bool
	UserAgent   string
	proxy       string
}

// NewConnection Create a new Connection instance
// APIBaseURL string, as per Worldpay documentation
// APIVersion e.g. v1
// serviceKey string, as per the API Keys section of your Worldpay account
// timeout int, number of milliseconds to wait before throwing a timeout error on requests
func NewConnection(apiBaseURL, appVersion, secureNetID, secureKey string, timeout int, sslCheck bool, proxy string) (service.Connection, error) {

	uaString := utils.BuildUserAgentString(runtime.GOOS, "?", runtime.GOARCH, runtime.Version(), sdk.LibVersion, sdk.APIVersion, sdk.LibLang, sdk.LibOwner)

	log.WithFields(
		log.Fields{
			"APIBaseURL":      apiBaseURL,
			"AppVersion":      appVersion,
			"SecureNetID":     secureNetID,
			"SecureKey":       "**obfuscated**",
			"timeout":         timeout,
			"UserAgentString": uaString,
			"Proxy":           proxy,
			"sslCheck":        sslCheck}).Debug("begin connectionImpl.NewConnection")

	result := &ConnectionImpl{
		APIBaseURL:  apiBaseURL,
		AppVersion:  appVersion,
		SecureNetID: secureNetID,
		SecureKey:   secureKey,
		Timeout:     timeout,
		SSLCheck:    sslCheck,
		UserAgent:   uaString,
		proxy:       proxy,
	}

	return result, nil
}

// Post a message to the worldpay server
// body is to be marshalled into JSON and submitted as the POST body
// url is appended to the connection APIEndpoint
// auth specifies whether the authorisation header (service key) should be added (true) or not (false)
func (conn *ConnectionImpl) Post(body []byte, url string, auth bool) ([]byte, error) {

	log.WithFields(log.Fields{"url": url, "auth": auth}).Debug("begin connectionImpl.Post()")

	finalURL := fmt.Sprintf("%s%s", conn.APIBaseURL, url)

	defer log.Debug("end connectionImpl.Post()")

	return conn.request("POST", finalURL, auth, body)
}

func (conn *ConnectionImpl) request(method, reqURL string, auth bool, body []byte) ([]byte, error) {

	log.Info(fmt.Sprintf("HTTP Request: %s %s", method, reqURL))

	log.WithFields(log.Fields{"method": method, "url": reqURL, "auth": auth, "body": string(body)}).Debug("begin connectionImpl.request()")

	reader := bytes.NewReader(body)

	req, err := http.NewRequest(method, reqURL, reader)

	if err != nil {

		log.WithField("Error", err.Error()).Error("Error creating new http request")

		return nil, err
	}

	log.Debug("Did create new HTTP request")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "http://www.worldpay.com")

	if auth {

		b64Auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", conn.SecureNetID, conn.SecureKey)))

		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64Auth))

		log.Debug("Did set authorisation key (service key)")
	}

	httpClient := http.Client{}

	if !strings.EqualFold(conn.proxy, "") {

		proxyURL, errParse := url.Parse(conn.proxy)

		if errParse != nil {

			return nil, errParse
		}

		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	for hk, hv := range req.Header {

		log.Debugf("Request Header. %s : %s", hk, hv)
	}

	log.Debugf("Request Body: %s", string(body))

	resp, err := httpClient.Do(req)

	if err != nil {

		log.WithField("Error", err.Error()).Error("Error making http request")

		return nil, err
	}

	log.Debug("Did make HTTP request")

	defer resp.Body.Close()
	defer log.Debug("Did close response body")

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	for hk, hv := range resp.Header {

		log.Debugf("Response Header. %s : %s", hk, hv)
	}

	log.Debugf("HTTP Response body: \n%s", string(bodyBytes))

	log.Debug("end connectionImpl.request()")

	return bodyBytes, err
}
