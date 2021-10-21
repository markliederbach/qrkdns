package dns

import (
	"context"
	"reflect"
)

// ProviderType labels specific supported DNS providers
type ProviderType string

const (
	// ProviderTypeCloudflare is a supported DNS client
	ProviderTypeCloudflare ProviderType = "cloudflare"
)

var (
	// SuportedProviders defines which providers this app supports
	SuportedProviders []ProviderType = []ProviderType{
		ProviderTypeCloudflare,
	}
)

// RecordType wraps the various DNS Record types
type RecordType string

const (
	// RecordTypeA is the DNS record type A
	RecordTypeA RecordType = "A"
)

// Record stores only the managed fields from a DNS record
type Record struct {
	ID      string     `json:"id"`
	Type    RecordType `json:"type"`
	Name    string     `json:"name"`
	Content string     `json:"content"`
	TTL     int        `json:"ttl"`
	Proxied bool       `json:"proxied"`
}

// Provider abstracts the interface necessary to call a downstream DNS provider's API
type Provider interface {
	// ApplyDNSARecord creates or updates a DNS record without creating a duplicate. It will also delete
	// other A records for the domain that don't match the provided IP address
	ApplyDNSARecord(ctx context.Context, subdomain, ipAddress string) (Record, error)
}

// Equal checks whether two records are equal (except for unmanaged fields)
func (d *Record) Equal(other Record, matchID bool) bool {
	if !matchID {
		// Temporarily copy ID
		other.ID = d.ID
	}
	return reflect.DeepEqual(*d, other)
}
