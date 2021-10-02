package ip

import "net/http"

// HTTPClient wraps the HTTP client used to make calls
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}
