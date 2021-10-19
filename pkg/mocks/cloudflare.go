package mocks

import (
	"context"

	sdk "github.com/cloudflare/cloudflare-go"
	"github.com/markliederbach/go-envy"
	"github.com/markliederbach/qrkdns/pkg/clients/cloudflare"
)

var (

	// Assert mock client matches the correct interface
	_ cloudflare.SDKClient = &MockCloudflareSDKClient{}

	// DefaultDNSRecord is used as the default option for the corresponding function
	DefaultDNSRecord sdk.DNSRecord = sdk.DNSRecord{
		ID:       "1234",
		Type:     "TXT",
		Name:     "test",
		Content:  "foobar",
		ZoneID:   "zone1234",
		ZoneName: "foo.bar",
		Proxied:  boolPtr(false),
	}

	// DefaultDNSRecords is used as the default option for the corresponding function
	DefaultDNSRecords []sdk.DNSRecord = []sdk.DNSRecord{DefaultDNSRecord}

	// DefaultZoneID is used as the default option for the corresponding function
	DefaultZoneID string = "zone1234"
)

// MockCloudflareSDKClient mocks the internal client from Cloudflare
type MockCloudflareSDKClient struct{}

func init() {
	sdkFunctions := []string{
		"ZoneIDByName",
		"DNSRecords",
		"DNSRecord",
		"CreateDNSRecord",
		"UpdateDNSRecord",
		"DeleteDNSRecord",
	}
	for _, functionName := range sdkFunctions {
		envy.ObjectChannels[functionName] = make(chan interface{}, 100)
		envy.ErrorChannels[functionName] = make(chan error, 100)
		envy.DefaultObjects[functionName] = struct{}{}
		envy.DefaultErrors[functionName] = nil
	}
}

// ZoneIDByName implements corresponding client function
func (c *MockCloudflareSDKClient) ZoneIDByName(zoneName string) (string, error) {
	functionName := "ZoneIDByName"
	obj := envy.GetObject(functionName)
	err := envy.GetError(functionName)
	switch obj := obj.(type) {
	case string:
		return obj, err
	default:
		return DefaultZoneID, err
	}
}

// DNSRecords implements corresponding client function
func (c *MockCloudflareSDKClient) DNSRecords(ctx context.Context, zoneID string, rr sdk.DNSRecord) ([]sdk.DNSRecord, error) {
	functionName := "DNSRecords"
	obj := envy.GetObject(functionName)
	err := envy.GetError(functionName)
	switch obj := obj.(type) {
	case []sdk.DNSRecord:
		return obj, err
	default:
		return DefaultDNSRecords, err
	}
}

// DNSRecord implements corresponding client function
func (c *MockCloudflareSDKClient) DNSRecord(ctx context.Context, zoneID string, recordID string) (sdk.DNSRecord, error) {
	functionName := "DNSRecord"
	obj := envy.GetObject(functionName)
	err := envy.GetError(functionName)
	switch obj := obj.(type) {
	case sdk.DNSRecord:
		return obj, err
	default:
		return DefaultDNSRecord, err
	}
}

// CreateDNSRecord implements corresponding client function
func (c *MockCloudflareSDKClient) CreateDNSRecord(ctx context.Context, zoneID string, rr sdk.DNSRecord) (*sdk.DNSRecordResponse, error) {
	functionName := "CreateDNSRecord"
	obj := envy.GetObject(functionName)
	err := envy.GetError(functionName)
	switch obj := obj.(type) {
	case *sdk.DNSRecordResponse:
		return obj, err
	default:
		return &sdk.DNSRecordResponse{Result: DefaultDNSRecord}, err
	}
}

// UpdateDNSRecord implements corresponding client function
func (c *MockCloudflareSDKClient) UpdateDNSRecord(ctx context.Context, zoneID string, recordID string, rr sdk.DNSRecord) error {
	functionName := "UpdateDNSRecord"
	err := envy.GetError(functionName)
	return err
}

// DeleteDNSRecord implements corresponding client function
func (c *MockCloudflareSDKClient) DeleteDNSRecord(ctx context.Context, zoneID string, recordID string) error {
	functionName := "DeleteDNSRecord"
	err := envy.GetError(functionName)
	return err
}

func boolPtr(val bool) *bool {
	return &val
}
