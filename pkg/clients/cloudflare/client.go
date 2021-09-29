package cloudflare

import (
	"github.com/cloudflare/cloudflare-go"
)

// DefaultClient implements the cloudflare client
type DefaultClient struct {
	Client cloudflare.API
}
