package ip

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DefaultClient implements the ip address client
type DefaultClient struct {
	IPServiceURL string
	// Client       *http.Client
	Client HTTPClient
}

// NewClient returns a new ip address client
func NewClient(ipServiceURL string) DefaultClient {
	return DefaultClient{
		IPServiceURL: ipServiceURL,
		Client:       &http.Client{},
	}
}

// GetExternalIPAddress returns the preferred outbound IP address used by this machine
func (c *DefaultClient) GetExternalIPAddress(ctx context.Context) (string, error) {
	response, err := c.Client.Get(c.IPServiceURL)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = response.Body.Close()
	}()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	trimmedBody := string(bytes.TrimSpace(body))

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Received status code %v: %v", response.StatusCode, trimmedBody)
	}
	return trimmedBody, nil
}
