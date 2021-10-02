package cloudflare

import (
	"context"

	sdk "github.com/cloudflare/cloudflare-go"
)

// DefaultClient implements the cloudflare client
type DefaultClient struct {
	// Client     *sdk.API
	Client     SDKClient
	AccountID  string
	DomainName string
	ZoneID     string
}

// LoadOption allows for modifying the client after it's created
type LoadOption func(client *DefaultClient) error

// WithTokenLoader is a load option for initializing the client with a token
func withTokenLoader(token string) LoadOption {
	loadOption := func(client *DefaultClient) error {
		cloudflareClient, err := sdk.NewWithAPIToken(token)
		if err != nil {
			return err
		}
		client.Client = cloudflareClient
		return nil
	}
	return loadOption
}

// NewClientWithToken is an initializer specifically for using an API token
func NewClientWithToken(ctx context.Context, accountID, domain, token string, opts ...LoadOption) (DefaultClient, error) {
	newOpts := []LoadOption{withTokenLoader(token)}
	newOpts = append(newOpts, opts...)
	return newClient(ctx, accountID, domain, newOpts...)
}

// newClient returns a new cloudflare client based on credentials
func newClient(ctx context.Context, accountID, domain string, opts ...LoadOption) (DefaultClient, error) {
	client := DefaultClient{
		Client:     &sdk.API{},
		AccountID:  accountID,
		DomainName: domain,
		ZoneID:     "",
	}

	for _, opt := range opts {
		if err := opt(&client); err != nil {
			return DefaultClient{}, err
		}
	}

	// Preload Zone ID
	_, err := client.GetZoneID(ctx)
	if err != nil {
		return DefaultClient{}, err
	}

	return client, nil
}

// GetZoneID returns and caches the Zone ID for the current client
func (c *DefaultClient) GetZoneID(ctx context.Context) (string, error) {
	if c.ZoneID != "" {
		return c.ZoneID, nil
	}

	zoneID, err := c.Client.ZoneIDByName(c.DomainName)
	if err != nil {
		return "", err
	}

	c.ZoneID = zoneID
	return c.ZoneID, nil
}

// ListDNSRecords returns all DNS records for the current zone
func (c *DefaultClient) ListDNSRecords(ctx context.Context) ([]sdk.DNSRecord, error) {
	records, err := c.Client.DNSRecords(ctx, c.ZoneID, sdk.DNSRecord{})
	if err != nil {
		return []sdk.DNSRecord{}, err
	}

	return records, nil
}
