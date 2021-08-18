package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/fleetdm/fleet/v4/orbit/pkg/certificate"
	"github.com/fleetdm/fleet/v4/pkg/secure"
	"github.com/fleetdm/fleet/v4/server/service"
)

func debugCommand() *cli.Command {
	return &cli.Command{
		Name:  "debug",
		Usage: "Tools for debugging Fleet",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Subcommands: []*cli.Command{
			debugProfileCommand(),
			debugCmdlineCommand(),
			debugHeapCommand(),
			debugGoroutineCommand(),
			debugTraceCommand(),
			debugArchiveCommand(),
			debugConnectionCommand(),
		},
	}
}

func writeFile(filename string, bytes []byte, mode os.FileMode) error {
	if err := ioutil.WriteFile(filename, bytes, mode); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Output written to %s\n", filename)
	return nil
}

func outfileName(name string) string {
	return fmt.Sprintf("fleet-%s-%s", name, time.Now().Format(time.RFC3339))
}

func debugProfileCommand() *cli.Command {
	return &cli.Command{
		Name:      "profile",
		Usage:     "Record a CPU profile from the Fleet server.",
		UsageText: "Record a 30-second CPU profile. The output can be analyzed with go tool pprof.",
		Flags: []cli.Flag{
			outfileFlag(),
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			profile, err := fleet.DebugPprof("profile")
			if err != nil {
				return err
			}

			outfile := getOutfile(c)
			if outfile == "" {
				outfile = outfileName("profile")
			}

			if err := writeFile(outfile, profile, defaultFileMode); err != nil {
				return errors.Wrap(err, "write profile to file")
			}

			return nil
		},
	}
}

func joinCmdline(cmdline string) string {
	var tokens []string
	for _, token := range strings.Split(string(cmdline), "\x00") {
		tokens = append(tokens, fmt.Sprintf("'%s'", token))
	}
	return fmt.Sprintf("[%s]", strings.Join(tokens, ", "))
}

func debugCmdlineCommand() *cli.Command {
	return &cli.Command{
		Name:  "cmdline",
		Usage: "Get the command line used to invoke the Fleet server.",
		Flags: []cli.Flag{
			outfileFlag(),
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			cmdline, err := fleet.DebugPprof("cmdline")
			if err != nil {
				return err
			}

			out := joinCmdline(string(cmdline))

			if outfile := getOutfile(c); outfile != "" {
				if err := writeFile(outfile, []byte(out), defaultFileMode); err != nil {
					return errors.Wrap(err, "write cmdline to file")
				}
				return nil
			}

			fmt.Println(out)

			return nil
		},
	}
}

func debugHeapCommand() *cli.Command {
	name := "heap"
	return &cli.Command{
		Name:      name,
		Usage:     "Report the allocated memory in the Fleet server.",
		UsageText: "Report the heap-allocated memory. The output can be analyzed with go tool pprof.",
		Flags: []cli.Flag{
			outfileFlag(),
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			profile, err := fleet.DebugPprof(name)
			if err != nil {
				return err
			}

			outfile := getOutfile(c)
			if outfile == "" {
				outfile = outfileName(name)
			}

			if err := writeFile(outfile, profile, defaultFileMode); err != nil {
				return errors.Wrapf(err, "write %s to file", name)
			}

			return nil
		},
	}
}

func debugGoroutineCommand() *cli.Command {
	name := "goroutine"
	return &cli.Command{
		Name:      name,
		Usage:     "Get stack traces of all goroutines (threads) in the Fleet server.",
		UsageText: "Get stack traces of all current goroutines (threads). The output can be analyzed with go tool pprof.",
		Flags: []cli.Flag{
			outfileFlag(),
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			profile, err := fleet.DebugPprof(name)
			if err != nil {
				return err
			}

			outfile := getOutfile(c)
			if outfile == "" {
				outfile = outfileName(name)
			}

			if err := writeFile(outfile, profile, defaultFileMode); err != nil {
				return errors.Wrapf(err, "write %s to file", name)
			}

			return nil
		},
	}
}

func debugTraceCommand() *cli.Command {
	name := "trace"
	return &cli.Command{
		Name:      name,
		Usage:     "Record an execution trace on the Fleet server.",
		UsageText: "Record a 1 second execution trace. The output can be analyzed with go tool trace.",
		Flags: []cli.Flag{
			outfileFlag(),
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			profile, err := fleet.DebugPprof(name)
			if err != nil {
				return err
			}

			outfile := getOutfile(c)
			if outfile == "" {
				outfile = outfileName(name)
			}

			if err := writeFile(outfile, profile, defaultFileMode); err != nil {
				return errors.Wrapf(err, "write %s to file", name)
			}

			return nil
		},
	}
}

func debugArchiveCommand() *cli.Command {
	return &cli.Command{
		Name:  "archive",
		Usage: "Create an archive with the entire suite of debug profiles.",
		Flags: []cli.Flag{
			outfileFlag(),
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Action: func(c *cli.Context) error {
			fleet, err := clientFromCLI(c)
			if err != nil {
				return err
			}

			profiles := []string{
				"allocs",
				"block",
				"cmdline",
				"goroutine",
				"heap",
				"mutex",
				"profile",
				"threadcreate",
				"trace",
			}

			outpath := getOutfile(c)
			if outpath == "" {
				outpath = outfileName("profiles-archive")
			}
			outfile := outpath + ".tar.gz"

			f, err := secure.OpenFile(outfile, os.O_CREATE|os.O_WRONLY, defaultFileMode)
			if err != nil {
				return errors.Wrap(err, "open archive for output")
			}
			defer f.Close()
			gzwriter := gzip.NewWriter(f)
			defer gzwriter.Close()
			tarwriter := tar.NewWriter(gzwriter)
			defer tarwriter.Close()

			for _, profile := range profiles {
				res, err := fleet.DebugPprof(profile)
				if err != nil {
					// Don't fail the entire process on errors. We'll take what
					// we can get if the servers are in a bad state and not
					// responding to all requests.
					fmt.Fprintf(os.Stderr, "Failed %s: %v\n", profile, err)
					continue
				}
				fmt.Fprintf(os.Stderr, "Ran %s\n", profile)

				if err := tarwriter.WriteHeader(
					&tar.Header{
						Name: outpath + "/" + profile,
						Size: int64(len(res)),
						Mode: defaultFileMode,
					},
				); err != nil {
					return errors.Wrapf(err, "write %s header", profile)
				}

				if _, err := tarwriter.Write(res); err != nil {
					return errors.Wrapf(err, "write %s contents", profile)
				}
			}

			fmt.Fprintf(os.Stderr, "Archive written to %s\n", outfile)

			return nil
		},
	}
}

func debugConnectionCommand() *cli.Command {
	const timeoutPerCheck = 10 * time.Second

	return &cli.Command{
		Name:      "connection",
		ArgsUsage: "[<address>]",
		Usage:     "Investigate the cause of a connection failure to the Fleet server.",
		Description: `Run a number of checks to debug a connection failure to the Fleet
server.

If <address> is provided, this is the address that is investigated,
otherwise the address of the provided context is used, with
the default context used if none is explicitly specified.`,
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
			debugFlag(),
			fleetCertificateFlag(),
		},
		Action: func(c *cli.Context) error {
			var addr string
			if narg := c.NArg(); narg > 0 {
				if narg > 1 {
					return errors.New("too many arguments")
				}
				addr = c.Args().First()
			}

			// ensure there is an address to debug
			cc, err := clientConfigFromCLI(c)
			if err != nil {
				return err
			}
			if addr != "" {
				cc.Address = addr
			}
			if cc.Address == "" {
				return errors.New(`set the Fleet API address with: fleetctl config set --address https://localhost:8080
or provide an <address> argument to debug: fleetctl debug connection localhost:8080`)
			}

			// it's ok if there is no scheme specified, add it automatically (to debug a non-https localhost address,
			// the scheme must be explicitly set).
			if !strings.Contains(cc.Address, "://") {
				cc.Address = "https://" + cc.Address
			}

			fleet, err := unauthenticatedClientFromConfig(cc, getDebug(c))
			if err != nil {
				return err
			}

			// print a summary of the address and TLS context that is investigated
			baseURL := fleet.BaseURL()
			fmt.Fprintf(c.App.Writer, "Debugging connection to %s; Configuration context: %s; ", baseURL.Hostname(), c.String("context"))
			rootCA := "(system)"
			if cc.RootCA != "" {
				rootCA = cc.RootCA
			}
			fmt.Fprintf(c.App.Writer, "Root CA: %s; ", rootCA)
			tlsMode := "secure"
			if cc.TLSSkipVerify {
				tlsMode = "insecure"
			}
			fmt.Fprintf(c.App.Writer, "TLS: %s.\n", tlsMode)

			// 1. Check that the url's host resolves to an IP address or is otherwise
			// a valid IP address directly. The ips may be used in a later check to
			// verify if the certificate is for one of them instead of the hostname.
			ips, err := resolveHostname(c.Context, timeoutPerCheck, baseURL.Hostname())
			if err != nil {
				return errors.Wrap(err, "Fail: resolve host")
			}
			fmt.Fprintf(c.App.Writer, "Success: can resolve host %s.\n", baseURL.Hostname())
			_ = ips

			// 2. Attempt a raw TCP connection to host:port.
			if err := dialHostPort(c.Context, timeoutPerCheck, baseURL.Host); err != nil {
				return errors.Wrap(err, "Fail: dial server")
			}
			fmt.Fprintf(c.App.Writer, "Success: can dial server at %s.\n", baseURL.Host)

			if cert := getFleetCertificate(c); cert != "" {
				// 3. Is the certificate valid at all (x509.Certificate.ParseCertificate?)
				// 4. Is the certificate valid for the hostname/IP address (x509.Certificate.VerifyHostname?)
				if err := checkFleetCert(c.Context, timeoutPerCheck, cert, baseURL.Hostname(), ips); err != nil {
					return errors.Wrap(err, "Fail: TLS certificate")
				}
				fmt.Fprintln(c.App.Writer, "Success: TLS certificate seems valid.")
			}

			// 5. Check that the server responds with expected responses (by
			// making a POST to /api/v1/osquery/enroll with an invalid
			// secret).
			if err := checkAPIEndpoint(c.Context, timeoutPerCheck, fleet); err != nil {
				return errors.Wrap(err, "Fail: agent API endpoint")
			}
			fmt.Fprintln(c.App.Writer, "Success: agent API endpoints are available.")

			return nil
		},
	}
}

func resolveHostname(ctx context.Context, timeout time.Duration, host string) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var r net.Resolver
	return r.LookupIP(ctx, "ip", host)
}

func dialHostPort(ctx context.Context, timeout time.Duration, addr string) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err == nil {
		conn.Close()
	}
	return err
}

func checkAPIEndpoint(ctx context.Context, timeout time.Duration, client *service.Client) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var enrollRes struct {
		Error       string `json:"error"`
		NodeInvalid bool   `json:"node_invalid"`
	}
	// make an enroll request with a deliberately invalid secret,
	// to see if we get the expected error json payload.
	res, err := client.DoContext(ctx, "POST",
		"/api/v1/osquery/enroll", "", map[string]string{
			"enroll_secret": "--invalid--",
		})
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&enrollRes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if res.StatusCode != http.StatusUnauthorized || enrollRes.Error == "" || !enrollRes.NodeInvalid {
		return fmt.Errorf("unexpected %d response", res.StatusCode)
	}
	return nil
}

func checkFleetCert(ctx context.Context, timeout time.Duration, certPath, host string, ips []net.IP) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// TODO: is it ok to use an orbit package from fleet? I remember reading that
	// those used to be distinct repos and can both be used independently, so
	// maybe we don't want to depend on each other in code either.
	certPool, err := certificate.LoadPEM(certPath)
	if err != nil {
		return err
	}

	// TODO: validation would ideally take a context so we can apply a timeout
	if err := certificate.ValidateConnection(certPool, "https://"+host); err != nil {
		return err
	}
	// TODO: ValidateConnection checks that it can connect with
	// InsecureSkipVerify, add a step that connects without skipping (if the
	// fleetctl config doesn't skip it)?

	return nil
}
