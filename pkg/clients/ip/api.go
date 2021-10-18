package ip

import "net/http"

// HTTPClient wraps the HTTP client used to make calls
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
