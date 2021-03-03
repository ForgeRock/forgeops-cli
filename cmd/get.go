package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/pkg/get"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var getFlags *genericclioptions.ConfigFlags

var getSecrets = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secret"},
	Short:   "Get the relevant ForgeRock Identity Platform secrets",
	Long: `
    Get the relevant ForgeRock Identity Platform secrets:
    * Reads secrets from the Kubernetes API
    * Returns json format or prints them in the console`,
	Example: `
    # Get the platform secrets from the default namespace in text format.
    forgeops get secrets
    
    # Get the secrets from the given namespace in json format.
    forgeops get secrets -o json -n myns`,

	RunE: func(cmd *cobra.Command, args []string) error {
		err := get.Secrets(clientFactory)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var getURLs = &cobra.Command{
	Use:     "urls",
	Aliases: []string{"url"},
	Short:   "Get the relevant ForgeRock Identity Platform URLs",
	Long: `
    Get the relevant ForgeRock Identity Platform URLs:
    * Obtains the FQDN from the "forgerock" ingress in the Kubernetes
    * Prints out the URLs in the console`,
	Example: `
    # Get the platform URLs from the default namespace in text format.
    forgeops get urls`,

	RunE: func(cmd *cobra.Command, args []string) error {
		err := get.URLs(clientFactory, "forgerock")
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get platform information",
	Long: `
    Get relevant ForgeRock Identity Platform information`,
	Example: `
    # Get platform secrets from the default namespace.
    forgeops get secrets`,

	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(cmd.Parent(), args)
		clientFactory = factory.NewFactory(getFlags)
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	// Install k8s flags
	getFlags = initK8sFlags(getCmd.PersistentFlags())

	getCmd.AddCommand(getSecrets)
	getCmd.AddCommand(getURLs)

	rootCmd.AddCommand(getCmd)
}
