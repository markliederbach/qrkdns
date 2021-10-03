package cloudflare

import (
	"context"

	sdk "github.com/cloudflare/cloudflare-go"
)

// SDKClient wraps the SDK client for Cloudflare
type SDKClient interface {
	ZoneIDByName(zoneName string) (string, error)
	DNSRecords(ctx context.Context, zoneID string, rr sdk.DNSRecord) ([]sdk.DNSRecord, error)
	CreateDNSRecord(ctx context.Context, zoneID string, rr sdk.DNSRecord) (*sdk.DNSRecordResponse, error)
	UpdateDNSRecord(ctx context.Context, zoneID string, recordID string, rr sdk.DNSRecord) error
	DeleteDNSRecord(ctx context.Context, zoneID string, recordID string) error
}
