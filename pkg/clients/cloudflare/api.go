package cloudflare

import (
	"context"

	sdk "github.com/cloudflare/cloudflare-go"
)

// SDKClient wraps the SDK client for Cloudflare
type SDKClient interface {
	ZoneIDByName(zoneName string) (string, error)
	DNSRecords(ctx context.Context, zoneID string, rr sdk.DNSRecord) ([]sdk.DNSRecord, error)
}
