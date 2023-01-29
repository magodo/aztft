package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/magodo/aztft/aztft"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		flagEnvironment    string
		flagSubscriptionId string
		flagAPI            bool
		flagImport         bool
	)

	app := &cli.App{
		Name:      "aztft",
		Version:   getVersion(),
		Usage:     "Find Azure resource's Terraform AzureRM provider resource type or/and id, together with any property-like resources, by its Azure resource ID",
		UsageText: "aztft [option] <ID>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "env",
				EnvVars:     []string{"AZTFT_ENV"},
				Usage:       `The environment. Can be one of "public", "china", "usgovernment".`,
				Destination: &flagEnvironment,
				Value:       "public",
			},
			&cli.StringFlag{
				Name:        "subscription-id",
				EnvVars:     []string{"AZTFT_SUBSCRIPTION_ID", "ARM_SUBSCRIPTION_ID"},
				Aliases:     []string{"s"},
				Required:    true,
				Usage:       "The subscription id",
				Destination: &flagSubscriptionId,
			},
			&cli.BoolFlag{
				Name:        "api",
				EnvVars:     []string{"AZTFT_API"},
				Usage:       `Allow to use Azure API to disambiguate matching results (e.g. whether a VM is a Linux VM or Windows VM)`,
				Destination: &flagAPI,
				Value:       false,
			},
			&cli.BoolFlag{
				Name:        "import",
				EnvVars:     []string{"AZTFT_IMPORT"},
				Usage:       `Print the TF import instruction`,
				Destination: &flagImport,
				Value:       false,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return fmt.Errorf("No ID specified")
			}
			if ctx.NArg() > 1 {
				return fmt.Errorf("More than one IDs specified")
			}

			var opt *aztft.APIOption
			if flagAPI {
				cloudCfg := cloud.AzurePublic
				switch strings.ToLower(flagEnvironment) {
				case "public":
					cloudCfg = cloud.AzurePublic
				case "usgovernment":
					cloudCfg = cloud.AzureGovernment
				case "china":
					cloudCfg = cloud.AzureChina
				default:
					return fmt.Errorf("unknown environment specified: %q", flagEnvironment)
				}

				if v, ok := os.LookupEnv("ARM_TENANT_ID"); ok {
					os.Setenv("AZURE_TENANT_ID", v)
				}
				if v, ok := os.LookupEnv("ARM_CLIENT_ID"); ok {
					os.Setenv("AZURE_CLIENT_ID", v)
				}
				if v, ok := os.LookupEnv("ARM_CLIENT_SECRET"); ok {
					os.Setenv("AZURE_CLIENT_SECRET", v)
				}
				if v, ok := os.LookupEnv("ARM_CLIENT_CERTIFICATE_PATH"); ok {
					os.Setenv("AZURE_CLIENT_CERTIFICATE_PATH", v)
				}

				clientOpt := arm.ClientOptions{
					ClientOptions: policy.ClientOptions{
						Cloud: cloudCfg,
						Telemetry: policy.TelemetryOptions{
							ApplicationID: "aztft",
							Disabled:      false,
						},
						Logging: policy.LogOptions{
							IncludeBody: true,
						},
					},
				}

				cred, err := azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{
					ClientOptions: clientOpt.ClientOptions,
					TenantID:      os.Getenv("ARM_TENANT_ID"),
				})
				if err != nil {
					return fmt.Errorf("failed to obtain a credential: %v", err)
				}

				opt = &aztft.APIOption{
					Cred:         cred,
					ClientOption: clientOpt,
				}
			}

			id := ctx.Args().First()
			var output []string
			if flagImport {
				types, ids, _, err := aztft.QueryTypeAndId(id, opt)
				if err != nil {
					log.Fatal(err)
				}
				for i := 0; i < len(types); i++ {
					output = append(output, fmt.Sprintf("terraform import %s.example %s", types[i].TFType, ids[i]))
				}
			} else {
				rts, _, err := aztft.QueryType(id, opt)
				if err != nil {
					log.Fatal(err)
				}
				for _, t := range rts {
					output = append(output, t.TFType)
				}
			}
			if len(output) == 0 {
				fmt.Println("No match")
				return nil
			}
			for _, line := range output {
				fmt.Println(line)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
