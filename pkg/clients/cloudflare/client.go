package cloudflare

import (
	sdk "github.com/cloudflare/cloudflare-go"
)

// DefaultClient implements the cloudflare client
type DefaultClient struct {
	Client     *sdk.API
	AccountID  string
	DomainName string
	ZoneID     string
}

// NewCloudflareClient returns a new cloudflare client based on credentials
func NewCloudflareClient(accountID, domain, token string) (DefaultClient, error) {
	cloudflareClient, err := sdk.NewWithAPIToken(token)
	if err != nil {
		return DefaultClient{}, err
	}

	zoneID, err := cloudflareClient.ZoneIDByName(domain)
	if err != nil {
		return DefaultClient{}, err
	}

	return DefaultClient{
		Client:     cloudflareClient,
		AccountID:  accountID,
		DomainName: domain,
		ZoneID:     zoneID,
	}, nil
}
