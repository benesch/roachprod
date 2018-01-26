package gce

import (
	"os/exec"

	"github.com/pkg/errors"
)

// CleanSSH removes entries in ~/.ssh/config.
func (c *Client) CleanSSH() error {
	args := []string{"compute", "config-ssh", "--project", c.project, "--quiet", "--remove"}
	cmd := exec.Command("gcloud", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Command: gcloud %s\nOutput: %s", args, output)
	}
	return nil
}

// ConfigSSH adds entries for VMs to ~/.ssh/config.
func (c *Client) ConfigSSH() error {
	args := []string{"compute", "config-ssh", "--project", c.project, "--quiet"}
	cmd := exec.Command("gcloud", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Command: gcloud %s\nOutput: %s", args, output)
	}
	return nil
}
