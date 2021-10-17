package mocks

import (
	"io"
	"net/http"
	"strings"

	"github.com/markliederbach/configr/mocks"
	configrmocks "github.com/markliederbach/configr/mocks"
	"github.com/markliederbach/qrkdns/pkg/clients/ip"
)

var (

	// Assert mock client matches the correct interface
	_ ip.HTTPClient = &MockHTTPClient{}

	// Assert mock reader matches the correct interface
	_ io.ReadCloser = &ErrorReader{}

	// DefaultExternalIPAddress is the default IP address returned
	DefaultExternalIPAddress = "1.2.3.4"

	// DefaultGetResponse is the default response for this function
	DefaultGetResponse *http.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(DefaultExternalIPAddress)),
	}
)

// MockHTTPClient mocks the internal client for http.Client
type MockHTTPClient struct{}

// ErrorReader is a mocked io.ReadCloser that returns erorr when reading
type ErrorReader struct {
	Error error
}

func init() {
	sdkFunctions := []string{
		"Do",
	}
	for _, functionName := range sdkFunctions {
		mocks.ObjectChannels[functionName] = make(chan interface{}, 100)
		mocks.ErrorChannels[functionName] = make(chan error, 100)
		mocks.DefaultObjects[functionName] = struct{}{}
		mocks.DefaultErrors[functionName] = nil
	}
}

// Read returns an error when reading
func (e *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, e.Error
}

// Close returns an error when closing the reader
func (e *ErrorReader) Close() error {
	return e.Error
}

// Do implements corresponding client function
func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	functionName := "Do"
	obj := configrmocks.GetObject(functionName)
	err := mocks.GetError(functionName)
	switch obj := obj.(type) {
	case *http.Response:
		return obj, err
	default:
		return DefaultGetResponse, err
	}
}
