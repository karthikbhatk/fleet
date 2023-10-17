package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/service"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/urfave/cli/v2"
)

func mdmCommand() *cli.Command {
	return &cli.Command{
		Name:  "mdm",
		Usage: "Run MDM commands against your hosts",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Subcommands: []*cli.Command{
			mdmRunCommand(),
		},
	}
}

func mdmRunCommand() *cli.Command {
	return &cli.Command{
		Name:    "run-command",
		Aliases: []string{"run_command"},
		Usage:   "Run a custom MDM command on macOS and Windows hosts.",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "hosts",
				Usage:    "Hosts specified by hostname, uuid, osquery_host_id or node_key that you want to target.",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "payload",
				Usage:    "A path to an XML file containing the raw MDM request payload.",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			client, err := clientFromCLI(c)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}

			// print an error if MDM is not configured
			if err := client.CheckAnyMDMEnabled(); err != nil {
				return err
			}

			// dedupe and remove any empty host identifier
			hostIdents := c.StringSlice("hosts")
			slices.Sort(hostIdents)
			hostIdents = slices.Compact(hostIdents)
			if len(hostIdents) > 0 && hostIdents[0] == "" {
				// because it is sorted, an empty ident can only be at the start
				hostIdents = hostIdents[1:]
			}
			if len(hostIdents) == 0 {
				return errors.New(`Required flag "hosts" not set`)
			}

			payloadFile := c.String("payload")
			payload, err := os.ReadFile(payloadFile)
			if err != nil {
				return fmt.Errorf("read payload: %w", err)
			}

			// fetch all specified hosts by their identifier
			var (
				hostUUIDs     []string
				notFoundCount int
				platform      string
			)
			for i, ident := range hostIdents {
				host, err := client.HostByIdentifier(ident)
				if err != nil {
					var nfe service.NotFoundErr
					if errors.As(err, &nfe) {
						notFoundCount++
						continue
					}
					var sce kithttp.StatusCoder
					if errors.As(err, &sce) {
						if sce.StatusCode() == http.StatusForbidden {
							return fmt.Errorf("Permission denied. You don't have permission to run an MDM command on this host: %w", err)
						}
					}
					return err
				}

				if host.Platform != platform && i > 0 {
					return errors.New(`Command can't run on hosts with different platforms. Make sure the hosts specified in the "hosts" flag are either all macOS or all Windows hosts.`)
				}
				platform = host.Platform

				// TODO(mna): this "On" check is brittle, but looks like it's the only
				// enrollment indication we have right now...
				if host.MDM.EnrollmentStatus == nil || !strings.HasPrefix(*host.MDM.EnrollmentStatus, "On") ||
					host.MDM.Name != fleet.WellKnownMDMFleet {
					return errors.New("Can't run the MDM command because the host doesn't have MDM turned on. Run the following command to see a list of hosts with MDM on: fleetctl get hosts --mdm")
				}

				hostUUIDs = append(hostUUIDs, host.UUID)
			}

			if len(hostUUIDs) == 0 {
				return errors.New("No hosts targeted. Make sure you provide a valid hostname, UUID, osquery host ID, or node key.")
			}
			if notFoundCount > 0 {
				// at least one was not found
				return errors.New("One or more targeted hosts don't exist. Make sure you provide a valid hostname, UUID, osquery host ID, or node key.")
			}

			result, err := client.RunMDMCommand(hostUUIDs, payload, platform)
			if err != nil {
				var sce kithttp.StatusCoder
				if errors.As(err, &sce) {
					if sce.StatusCode() == http.StatusForbidden {
						return fmt.Errorf("Permission denied. You don't have permission to run an MDM command on this host: %w", err)
					}
					if sce.StatusCode() == http.StatusUnsupportedMediaType {
						// this status code is only returned for an Apple MDM command, so this message is correct.
						return fmt.Errorf("The payload isn't valid. Please provide a valid MDM command in the form of a plist-encoded XML file: %w", err)
					}
				}
				return err
			}

			fmt.Fprintf(c.App.Writer, `
Hosts will run the command the next time they check into Fleet.

Copy and run this command to see results:

fleetctl get mdm-command-results --id=%v
`, result.CommandUUID)

			return nil
		},
	}
}
