package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/pkg/install"
	"github.com/spf13/cobra"
)

// cmd globals config

var quickstart = &cobra.Command{
	Use:     "quickstart",
	Aliases: []string{"qs"},
	Short:   "Installs the ForgeRock Cloud Deployment Quickstart (CDQ)",
	Long: `
	Installs the ForgeRock Cloud Deployment Quickstart (CDQ):
	  * Applies the latest quickstart manifest
	  * use --version to specify a specific CDQ version to install
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := install.Quickstart(clientFactory, tag)
		return err
	},
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install common platform components",
	Long: `
	Install common platform components
    `,
	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		clientFactory = factory.NewFactory(kubeConfigFlags)
	},
}

func init() {
	// Install command-specific flags
	installCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Tag/version to install")

	installCmd.AddCommand(quickstart)
	rootCmd.AddCommand(installCmd)
}
