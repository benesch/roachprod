package gce

import (
	"fmt"
	"strings"
)

// ActiveAccount determines the username of the account that the client has
// authenticated with.
func (c *Client) ActiveAccount() (string, error) {
	args := []string{"auth", "list", "--format", "json", "--filter", "status~ACTIVE"}

	type jsonAuth struct {
		Account string
		Status  string
	}

	accounts := []jsonAuth{}
	if err := runJSONCommand(args, &accounts); err != nil {
		return "", err
	}

	if len(accounts) != 1 {
		return "", fmt.Errorf("no active accounts found, please configure gcloud")
	}

	if !strings.HasSuffix(accounts[0].Account, c.domain) {
		return "", fmt.Errorf("active account %q does no belong to domain %s",
			accounts[0].Account, c.domain)
	}

	username := strings.Split(accounts[0].Account, "@")[0]
	return username, nil
}
