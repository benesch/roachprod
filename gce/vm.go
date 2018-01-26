package gce

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// VM represents a Google Compute Engine virtual machine.
type VM struct {
	Name        string
	Lifetime    time.Duration
	MachineType string
	LocalSSD    bool
	Zone        string
}

// CreateVMs creates the specified virtual machines.
func (c *Client) CreateVMs(vms []VM) error {
	// Create GCE startup script file.
	filename, err := writeStartupScript()
	if err != nil {
		return errors.Wrapf(err, "could not write GCE startup script to temp file")
	}
	defer os.Remove(filename)

	// Fixed args.
	args := []string{
		"compute", "instances", "create",
		"--subnet", "default",
		"--maintenance-policy", "MIGRATE",
		"--service-account", "21965078311-compute@developer.gserviceaccount.com",
		"--scopes", "default,storage-rw",
		"--image", "ubuntu-1604-xenial-v20171002",
		"--image-project", "ubuntu-os-cloud",
		"--boot-disk-size", "10",
		"--boot-disk-type", "pd-ssd",
	}

	// Dynamic args.
	if opts.UseLocalSSD {
		args = append(args, "--local-ssd", "interface=SCSI")
	}
	args = append(args, "--machine-type", opts.MachineType)
	args = append(args, "--labels", fmt.Sprintf("lifetime=%s", opts.Lifetime))

	args = append(args, "--metadata-from-file", fmt.Sprintf("startup-script=%s", filename))
	args = append(args, "--project", project)

	var g errgroup.Group

	// This is calculating the number of machines to allocate per zone by taking the ceiling of the the total number
	// of machines left divided by the number of zones left. If the the number of machines isn't
	// divisible by the number of zones, then the extra machines will be allocated one per zone until there are
	// no more extra machines left.
	for i < len(names) {
		argsWithZone := append(args[:len(args):len(args)], "--zone", zones[ct])
		ct++
		argsWithZone = append(argsWithZone, names[i:i+nodesPerZone]...)
		i += nodesPerZone

		totalNodes -= float64(nodesPerZone)
		totalZones -= 1
		nodesPerZone = int(math.Ceil(totalNodes / totalZones))

		g.Go(func() error {
			cmd := exec.Command("gcloud", argsWithZone...)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return errors.Wrapf(err, "Command: gcloud %s\nOutput: %s", args, output)
			}
			return nil
		})

	}

	return g.Wait()
}

func ListVMs() ([]VM, error) {
	args := []string{"compute", "instances", "list", "--project", project, "--format", "json"}
	vms := []jsonVM{}

	if err := runJSONCommand(args, &vms); err != nil {
		return nil, err
	}

	for _, vms := range vms {

	}

	return vms, nil
}

func DeleteVMs(names []string, zones []string) error {
	zoneMap := make(map[string][]string)
	for i, name := range names {
		zoneMap[zones[i]] = append(zoneMap[zones[i]], name)
	}

	var g errgroup.Group

	for zone, names := range zoneMap {
		args := []string{
			"compute", "instances", "delete",
			"--delete-disks", "all",
		}

		args = append(args, "--project", project)
		args = append(args, "--zone", zone)
		args = append(args, names...)

		g.Go(func() error {
			cmd := exec.Command("gcloud", args...)

			output, err := cmd.CombinedOutput()
			if err != nil {
				return errors.Wrapf(err, "Command: gcloud %s\nOutput: %s", args, output)
			}
			return nil
		})
	}

	return g.Wait()
}

func ExtendVM(name string, lifetime time.Duration) error {
	args := []string{"compute", "instances", "add-labels"}

	args = append(args, "--project", project)
	args = append(args, "--zone", zone)
	args = append(args, "--labels", fmt.Sprintf("lifetime=%s", lifetime))
	args = append(args, name)

	cmd := exec.Command("gcloud", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Command: gcloud %s\nOutput: %s", args, output)
	}
	return nil
}

type jsonVM struct {
	Name              string
	Labels            map[string]string
	CreationTimestamp time.Time
	NetworkInterfaces []struct {
		NetworkIP     string
		AccessConfigs []struct {
			Name  string
			NatIP string
		}
	}
	Zone string
}
