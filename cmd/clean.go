package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"

	"github.com/ForgeRock/forgeops-cli/pkg/clean"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var cleanFlags *genericclioptions.ConfigFlags

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove any remaining platform components from the given namespace",
	Long: `
    Remove any remaining platform components from the given namespace`,
	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(cmd.Parent(), args)
		clientFactory = factory.NewFactory(cleanFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := clean.Clean(clientFactory, skipUserConfirmation)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	// Install k8s flags
	cleanFlags = initK8sFlags(cleanCmd.PersistentFlags())

	// clean command-specific flags
	cleanCmd.PersistentFlags().BoolVarP(&skipUserConfirmation, "yes", "y", false, "Do not prompt for confirmation")

	rootCmd.AddCommand(cleanCmd)
}
