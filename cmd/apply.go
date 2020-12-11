package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/pkg/apply"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var applyFlags *genericclioptions.ConfigFlags

var quickstart = &cobra.Command{
	Use:     "quickstart",
	Aliases: []string{"qs"},
	Short:   "Installs the ForgeRock Cloud Deployment Quickstart (CDQ)",
	Long: `
    Installs the ForgeRock Cloud Deployment Quickstart (CDQ):
    * Applies the latest quickstart manifest
    * use --tag to specify a specific CDQ version to install`,
	Example: `
      # Install the "latest" CDQ in the "default" namespace.
      forgeops apply quickstart
    
      # Install the CDQ in the "default" namespace.
      forgeops apply quickstart -t 2020.10.28-AlSugoDiNoci
      
      # Install the CDQ in a different namespace.
      forgeops apply quickstart -t 2020.10.28-AlSugoDiNoci -n mynamespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := apply.Quickstart(clientFactory, tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var secretAgent = &cobra.Command{
	Use:     "sa",
	Aliases: []string{"secret-agent"},
	Short:   "Installs the ForgeRock secret-agent",
	Long: `
    Installs the ForgeRock secret-agent:
    * Applies the latest secret-agent manifest
    * use --tag to specify a specific secret-agent version to install`,
	Example: `
      # Install the "latest" secret-agent.
      forgeops apply sa

      # Install a specific version of the secret-agent.
      forgeops apply sa -t v0.2.1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := apply.SecretAgent(clientFactory, tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Install common platform components",
	Long: `
	Apply common platform components
    `,
	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		clientFactory = factory.NewFactory(applyFlags)
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	// Install k8s flags
	applyFlags = initK8sFlags(applyCmd.PersistentFlags())

	// Apply command-specific flags
	applyCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Tag/version to apply")

	applyCmd.AddCommand(quickstart)
	applyCmd.AddCommand(secretAgent)

	rootCmd.AddCommand(applyCmd)
}
