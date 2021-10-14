package cloudflare

import (
	"context"
	"fmt"
	"reflect"

	sdk "github.com/cloudflare/cloudflare-go"
	log "github.com/sirupsen/logrus"
)

// DefaultClient implements the cloudflare client
type DefaultClient struct {
	// Client *sdk.API
	Client     SDKClient
	AccountID  string
	DomainName string
	ZoneID     string
}

// RecordType wraps the various DNS Record types
type RecordType string

const (
	// RecordTypeA is the DNS record type A
	RecordTypeA RecordType = "A"
)

// DNSRecord stores only the managed fields from the
// Cloudflare DNSRecord struct
type DNSRecord struct {
	ID      string     `json:"id"`
	Type    RecordType `json:"type"`
	Name    string     `json:"name"`
	Content string     `json:"content"`
	TTL     int        `json:"ttl"`
	Proxied bool       `json:"proxied"`
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

// ListDNSARecords returns all DNS records for the provided subdomain
func (c *DefaultClient) ListDNSARecords(ctx context.Context, subdomain string) ([]sdk.DNSRecord, error) {
	records, err := c.Client.DNSRecords(ctx, c.ZoneID, sdk.DNSRecord{Type: string(RecordTypeA), Name: fqdn(subdomain, c.DomainName)})
	if err != nil {
		return []sdk.DNSRecord{}, err
	}

	return records, nil
}

// GetDNSRecord retrieves a DNS record by ID
func (c *DefaultClient) GetDNSRecord(ctx context.Context, recordID string) (DNSRecord, error) {
	response, err := c.Client.DNSRecord(ctx, c.ZoneID, recordID)
	if err != nil {
		return DNSRecord{}, err
	}

	return FromCloudFlareDNSRecord(response), nil
}

// CreateDNSARecord creates a new DNS A record for the provided subdomain and IP Address
func (c *DefaultClient) CreateDNSARecord(ctx context.Context, record DNSRecord) (DNSRecord, error) {
	response, err := c.Client.CreateDNSRecord(ctx, c.ZoneID, record.ToCloudFlareDNSRecord())
	if err != nil {
		return DNSRecord{}, err
	}

	return FromCloudFlareDNSRecord(response.Result), nil
}

// UpdateDNSARecord updates an existing DNS A record for the provided subdomain and IP Address
func (c *DefaultClient) UpdateDNSARecord(ctx context.Context, recordID string, record DNSRecord) error {
	err := c.Client.UpdateDNSRecord(ctx, c.ZoneID, recordID, record.ToCloudFlareDNSRecord())
	if err != nil {
		return err
	}

	return nil
}

// DeleteDNSARecord deletes an existing DNS A record for the provided record ID
func (c *DefaultClient) DeleteDNSARecord(ctx context.Context, record DNSRecord) error {
	err := c.Client.DeleteDNSRecord(ctx, c.ZoneID, record.ID)
	if err != nil {
		return err
	}

	return nil
}

// ApplyDNSARecord creates or updates a DNS record without creating a duplicate. It will also delete
// other A records for the domain that don't match the provided IP address
func (c *DefaultClient) ApplyDNSARecord(ctx context.Context, subdomain, ipAddress string) (DNSRecord, error) {
	expectedRecord := BuildDNSARecord(subdomain, c.DomainName, ipAddress)
	contextLog := log.WithField("expected_record", expectedRecord)

	sdkRecords, err := c.ListDNSARecords(ctx, subdomain)
	if err != nil {
		return DNSRecord{}, err
	}

	existingRecords := ConvertDNSRecordList(sdkRecords)

	chosenRecord := DNSRecord{}

	// First, look for any record with a matching IP address because
	// Cloudflare's unique key is (name, content)
	for _, record := range existingRecords {
		existingRecordLog := contextLog.WithField("existing_record", record)
		if record.Content == ipAddress {
			existingRecordLog.Debugf("Reusing found record")
			chosenRecord = record
			break
		}
	}

	// If we found a matching record, check if it needs to be updated
	if !chosenRecord.Equal(DNSRecord{}, false) {

		contextLog = contextLog.WithField("chosen_record", chosenRecord)

		if !chosenRecord.Equal(expectedRecord, false) {
			contextLog.Debugf("Updating record")
			err = c.UpdateDNSARecord(ctx, chosenRecord.ID, expectedRecord)
			if err != nil {
				return DNSRecord{}, err
			}

			// Update local copy of record
			chosenRecord, err = c.GetDNSRecord(ctx, chosenRecord.ID)
			if err != nil {
				return DNSRecord{}, err
			}
		} else {
			contextLog.Debugf("Record is already up to date")
		}
	} else {
		// Otherwise, we need to create a new record
		contextLog.Debugf("Creating new record")
		chosenRecord, err = c.CreateDNSARecord(ctx, expectedRecord)
		if err != nil {
			return DNSRecord{}, err
		}
		contextLog = contextLog.WithField("chosen_record", chosenRecord)
	}

	// At this point, we've either updated or created an new record
	// with the correct data. All other records should be removed.
	for _, record := range existingRecords {
		if record.ID == chosenRecord.ID {
			// Skip the record we've chosen to manage
			continue
		}
		contextLog.WithField("existing_record", record).Debugf("Deleting extra record")
		err = c.DeleteDNSARecord(ctx, record)
		if err != nil {
			return DNSRecord{}, err
		}
	}

	return chosenRecord, nil
}

// BuildDNSARecord constructs a consistent DNS record across the client
func BuildDNSARecord(subdomain, domainName, ipAddress string) DNSRecord {
	return DNSRecord{
		Type:    RecordTypeA,
		Name:    fqdn(subdomain, domainName),
		Content: ipAddress,
		TTL:     1,
		Proxied: false,
	}
}

// fqdn concatenates a subdomain name with the base domain and returns the FQDN
func fqdn(subdomain, domainName string) string {
	return fmt.Sprintf("%v.%v", subdomain, domainName)
}

// ConvertDNSRecordList converts a list of Cloudflare DNS records to
// locally-managed DNS Records
func ConvertDNSRecordList(sdkRecords []sdk.DNSRecord) []DNSRecord {
	records := []DNSRecord{}
	for _, record := range sdkRecords {
		records = append(records, FromCloudFlareDNSRecord(record))
	}
	return records
}

// FromCloudFlareDNSRecord converts a CloudFlare DNS Record struct to
// one managed and controlled by this client
func FromCloudFlareDNSRecord(record sdk.DNSRecord) DNSRecord {
	return DNSRecord{
		ID:      record.ID,
		Type:    RecordType(record.Type),
		Name:    record.Name,
		Content: record.Content,
		TTL:     record.TTL,
		Proxied: *record.Proxied,
	}
}

// ToCloudFlareDNSRecord converts a local DNS record to one accepted
// by the CloudFlare SDK
func (d *DNSRecord) ToCloudFlareDNSRecord() sdk.DNSRecord {
	return sdk.DNSRecord{
		ID:      d.ID,
		Type:    string(d.Type),
		Name:    d.Name,
		Content: d.Content,
		TTL:     d.TTL,
		Proxied: &d.Proxied,
	}
}

// Equal checks whether two records are equal (except for unmanaged fields)
func (d *DNSRecord) Equal(other DNSRecord, matchID bool) bool {
	if !matchID {
		// Temporarily copy ID
		other.ID = d.ID
	}
	return reflect.DeepEqual(*d, other)
}
