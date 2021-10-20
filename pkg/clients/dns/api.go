package dns

import (
	"context"
	"reflect"
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
