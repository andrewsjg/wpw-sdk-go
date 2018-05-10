package hte

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/wpwerrors"
)

// ClientHTTP an interface for HTTP requests
type ClientHTTP interface {
	Get(url string) ([]byte, error)
	PostJSON(url string, postBody []byte) ([]byte, int, error)
}

// ClientHTTPImpl implementation of HTEClientHTTP
type ClientHTTPImpl struct{}

// NewHTEClientHTTP create a new instance of ClientHTTPImpl
func NewHTEClientHTTP() (ClientHTTP, error) {

	return &ClientHTTPImpl{}, nil
}

// Get Helper function to make a HTTP GET request
func (client *ClientHTTPImpl) Get(url string) ([]byte, error) {

	response, err := http.Get(url)

	if err != nil {
		return nil, wpwerrors.GetError(wpwerrors.HTTP_GET, err.Error())
	}

	byteResponse, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, wpwerrors.GetError(wpwerrors.IO_READ, err.Error())
	}

	return byteResponse, nil
}

// PostJSON Helper function to make a http POST request
func (client *ClientHTTPImpl) PostJSON(url string, postBody []byte) ([]byte, int, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))

	if err != nil {
		return nil, 0, wpwerrors.GetError(wpwerrors.HTTP_REQUEST_POST, err)
	}

	req.Header.Add("Content-Type", "application/json")

	_client := &http.Client{}

	resp, err := _client.Do(req)

	if err != nil {
		return nil, 0, wpwerrors.GetError(wpwerrors.HTTP_REQUEST_DO, err)
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		return nil, resp.StatusCode, wpwerrors.GetError(wpwerrors.IO_READ, err)
	}

	return bodyBytes, resp.StatusCode, nil
}
