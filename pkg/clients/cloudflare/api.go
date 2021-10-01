package cloudflare

// SDKClient wraps the SDK client for Cloudflare
type SDKClient interface {
	ZoneIDByName(zoneName string) (string, error)
}
