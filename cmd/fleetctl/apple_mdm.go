package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fleetdm/fleet/v4/server/mdm/apple/scep/scep_ca"
	"github.com/micromdm/nanodep/tokenpki"
	"github.com/urfave/cli/v2"
)

func appleMDMCommand() *cli.Command {
	return &cli.Command{
		Name:  "apple-mdm",
		Usage: "Apple MDM functionality",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
			debugFlag(),
		},
		Subcommands: []*cli.Command{
			appleMDMSetupCommand(),
			appleMDMEnrollmentsCommand(),
		},
	}
}

func appleMDMSetupCommand() *cli.Command {
	return &cli.Command{
		Name:  "setup",
		Usage: "Setup commands for Apple MDM",
		Subcommands: []*cli.Command{
			appleMDMSetupSCEPCommand(),
			appleMDMSetupAPNSCommand(),
			appleMDMSetupDEPCommand(),
		},
	}
}

func appleMDMSetupSCEPCommand() *cli.Command {
	// TODO(lucas): Define workflow when SCEP CA certificate expires.
	var (
		validityYears      int
		cn                 string
		organization       string
		organizationalUnit string
		country            string
	)
	return &cli.Command{
		Name:  "scep",
		Usage: "Create SCEP certificate authority",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "validity-years",
				Usage:       "Validity of the SCEP CA certificate in years",
				Required:    true,
				Destination: &validityYears,
			},
			&cli.StringFlag{
				Name:        "cn",
				Usage:       "Common name to set in the SCEP CA certificate",
				Required:    true,
				Destination: &cn,
			},
			&cli.StringFlag{
				Name:        "organization",
				Usage:       "Organization to set in the SCEP CA certificate",
				Required:    true,
				Destination: &organization,
			},
			&cli.StringFlag{
				Name:        "organizational-unit",
				Usage:       "Organizational unit to set in the SCEP CA certificate",
				Required:    true,
				Destination: &organizationalUnit,
			},
			&cli.StringFlag{
				Name:        "country",
				Usage:       "Country to set in the SCEP CA certificate",
				Required:    true,
				Destination: &country,
			},
		},
		Action: func(c *cli.Context) error {
			certPEM, keyPEM, err := scep_ca.Create(validityYears, cn, organization, organizationalUnit, country)
			if err != nil {
				return fmt.Errorf("creating SCEP CA: %w", err)
			}
			const (
				certPath = "fleet-mdm-apple-scep.crt"
				keyPath  = "fleet-mdm-apple-scep.key"
			)
			if err := os.WriteFile(certPath, certPEM, 0o600); err != nil {
				return fmt.Errorf("write %s: %w", certPath, err)
			}
			if err := os.WriteFile(keyPath, keyPEM, 0o600); err != nil {
				return fmt.Errorf("write %s: %w", keyPath, err)
			}
			fmt.Printf("Successfully generated SCEP CA: %s, %s.\n", certPath, keyPath)
			fmt.Printf("Set FLEET_MDM_APPLE_SCEP_CA_CERT_PEM=$(cat %s) FLEET_MDM_APPLE_SCEP_CA_KEY_PEM=$(cat %s) when running Fleet.\n", certPath, keyPath)
			return nil
		},
	}
}

func appleMDMSetupAPNSCommand() *cli.Command {
	return &cli.Command{
		Name:  "apns",
		Usage: "Commands to setup APNS certificate",
		Subcommands: []*cli.Command{
			appleMDMSetupAPNSInitCommand(),
			appleMDMSetupAPNSFinalizeCommand(),
		},
	}
}

func appleMDMSetupAPNSInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Start APNS certificate configuration",
		Action: func(c *cli.Context) error {
			// TODO(lucas): Implement command.
			fmt.Println("Not implemented yet.")
			fmt.Printf("TODO(lucas): Add environment variables to set with these generated files.\n")
			return nil
		},
	}
}

func appleMDMSetupAPNSFinalizeCommand() *cli.Command {
	var encryptedReq string
	return &cli.Command{
		Name:  "finalize",
		Usage: "Finalize APNS certificate configuration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "encrypted-req",
				Usage:       "File path of the encrypted .req p7m file",
				Destination: &encryptedReq,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			// TODO(lucas): Implement command.
			fmt.Println("Not implemented yet.")
			return nil
		},
	}
}

func appleMDMSetupDEPCommand() *cli.Command {
	return &cli.Command{
		Name:  "dep",
		Usage: "Configure DEP token",
		Subcommands: []*cli.Command{
			appleMDMSetDEPTokenInitCommand(),
			appleMDMSetDEPTokenFinalizeCommand(),
		},
	}
}

func appleMDMSetDEPTokenInitCommand() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "Start DEP token configuration",
		Flags: []cli.Flag{
			configFlag(),
			contextFlag(),
		},
		Action: func(c *cli.Context) error {
			// TODO(lucas): Check validity days default value.
			const (
				cn           = "fleet"
				validityDays = 1
				pemCertPath  = "fleet-mdm-apple-dep.crt"
				pemKeyPath   = "fleet-mdm-apple-dep.key"
			)
			key, cert, err := tokenpki.SelfSignedRSAKeypair(cn, validityDays)
			if err != nil {
				return fmt.Errorf("generate encryption keypair: %w", err)
			}
			pemCert := tokenpki.PEMCertificate(cert.Raw)
			pemKey := tokenpki.PEMRSAPrivateKey(key)
			if err := os.WriteFile(pemCertPath, pemCert, defaultFileMode); err != nil {
				return fmt.Errorf("write certificate: %w", err)
			}
			if err := os.WriteFile(pemKeyPath, pemKey, defaultFileMode); err != nil {
				return fmt.Errorf("write private key: %w", err)
			}
			fmt.Printf("Successfully generated DEP public and private key: %s, %s\n", pemCertPath, pemKeyPath)
			fmt.Printf("Upload %s to your Apple Business MDM server. (Don't forget to click \"Save\" after uploading it.)", pemCertPath)
			return nil
		},
	}
}

func appleMDMSetDEPTokenFinalizeCommand() *cli.Command {
	var (
		pemCertPath        string
		pemKeyPath         string
		encryptedTokenPath string
	)
	return &cli.Command{
		Name:  "finalize",
		Usage: "Finalize DEP token configuration for an automatic enrollment",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "certificate",
				Usage:       "Path to the certificate generated in the init step",
				Destination: &pemCertPath,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "private-key",
				Usage:       "Path to the private key file generated in the init step",
				Destination: &pemKeyPath,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "encrypted-token",
				Usage:       "Path to the encrypted token file downloaded from Apple Business (*.p7m)",
				Destination: &encryptedTokenPath,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			pemCert, err := os.ReadFile(pemCertPath)
			if err != nil {
				return fmt.Errorf("read certificate: %w", err)
			}
			depCert, err := tokenpki.CertificateFromPEM(pemCert)
			if err != nil {
				return fmt.Errorf("parse certificate: %w", err)
			}
			pemKey, err := os.ReadFile(pemKeyPath)
			if err != nil {
				return fmt.Errorf("read private key: %w", err)
			}
			depKey, err := tokenpki.RSAKeyFromPEM(pemKey)
			if err != nil {
				return fmt.Errorf("parse private key: %w", err)
			}
			encryptedToken, err := os.ReadFile(encryptedTokenPath)
			if err != nil {
				return fmt.Errorf("read encrypted token: %w", err)
			}
			token, err := tokenpki.DecryptTokenJSON(encryptedToken, depCert, depKey)
			if err != nil {
				return fmt.Errorf("decrypt token: %w", err)
			}
			tokenPath := "fleet-mdm-apple-dep.token"
			if err := os.WriteFile(tokenPath, token, defaultFileMode); err != nil {
				return fmt.Errorf("write token file: %w", err)
			}
			fmt.Printf("Successfully generated token file: %s.\n", tokenPath)
			// TODO(lucas): Delete pemCertPath, pemKeyPath and encryptedTokenPath files?
			fmt.Printf("Set FLEET_MDM_APPLE_DEP_TOKEN=$(cat %s) when running Fleet.\n", tokenPath)
			return nil
		},
	}
}

func appleMDMEnrollmentsCommand() *cli.Command {
	return &cli.Command{
		Name:  "enrollments",
		Usage: "Commands to manage enrollments",
		Subcommands: []*cli.Command{
			appleMDMEnrollmentsCreateAutomaticCommand(),
			appleMDMEnrollmentsCreateManualCommand(),
		},
	}
}

func appleMDMEnrollmentsCreateAutomaticCommand() *cli.Command {
	var (
		enrollmentName string
		configPath     string
		depConfigPath  string
	)
	return &cli.Command{
		Name:  "create-automatic",
		Usage: "Create a new automatic enrollment",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Name of the automatic enrollment",
				Destination: &enrollmentName,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "enroll-config",
				Usage:       "JSON file with enrollment config",
				Destination: &configPath,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "profile",
				Usage:       "JSON file with fields defined in https://developer.apple.com/documentation/devicemanagement/profile",
				Destination: &depConfigPath,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			// TODO(lucas): Document behavior: For MVP-Dogfood, we will only support one
			// automatic enrollment per team (or global).
			config, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("read enrollment config: %w", err)
			}
			profile, err := os.ReadFile(depConfigPath)
			if err != nil {
				return fmt.Errorf("read dep profile: %w", err)
			}
			fleet, err := clientFromCLI(c)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}
			depProfile := json.RawMessage(profile)
			enrollment, err := fleet.CreateEnrollment(enrollmentName, config, &depProfile)
			if err != nil {
				return fmt.Errorf("create enrollment: %w", err)
			}
			fmt.Printf("Automatic enrollment created, id: %d\n", enrollment.ID)
			return nil
		},
	}
}

func appleMDMEnrollmentsCreateManualCommand() *cli.Command {
	var (
		enrollmentName string
		configPath     string
	)
	return &cli.Command{
		Name:  "create-manual",
		Usage: "Create a new manual enrollment",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Usage:       "Name of the manual enrollment",
				Destination: &enrollmentName,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "enroll-config",
				Usage:       "JSON file with enrollment config",
				Destination: &configPath,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			config, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("read enrollment config: %w", err)
			}
			fleet, err := clientFromCLI(c)
			if err != nil {
				return fmt.Errorf("create client: %w", err)
			}
			enrollment, err := fleet.CreateEnrollment(enrollmentName, config, nil)
			if err != nil {
				return fmt.Errorf("create enrollment: %w", err)
			}
			fmt.Printf("Manual enrollment created, id: %d\n", enrollment.ID)
			return nil
		},
	}
}
