package gce

// Client is an API client for Google Compute Engine.
type Client struct {
	project string
	domain  string
}

// NewClient initializes a Client.
func NewClient(project, domain string) *Client {
	return &Client{
		project: project,
		domain:  domain,
	}
}
