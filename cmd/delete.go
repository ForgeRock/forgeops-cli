package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/pkg/delete"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var deleteFlags *genericclioptions.ConfigFlags
var skipUserDelQ bool

var deleteQuickstart = &cobra.Command{
	Use:     "quickstart",
	Aliases: []string{"qs"},
	Short:   "Uninstalls the ForgeRock Cloud Deployment Quickstart (CDQ)",
	Long: `
    Uninstalls the ForgeRock Cloud Deployment Quickstart (CDQ):
    * Deletes the quickstart deployment`,
	Example: `
    # Delete the CDQ from the "default" namespace.
    forgeops delete quickstart
    
    # Delete the CDQ from a given namespace.
    forgeops delete quickstart -n mynamespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := delete.Quickstart(clientFactory, tag, skipUserDelQ)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteSecretAgent = &cobra.Command{
	Use:     "sa",
	Aliases: []string{"secret-agent"},
	Short:   "Uninstalls the ForgeRock secret-agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := delete.GHResource(clientFactory, "ForgeRock/secret-agent", "secret-agent.yaml", tag, true, skipUserDelQ)
		return err
	},
	// Hide this command from help docs
	Hidden:            true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteDsOperator = &cobra.Command{
	Use:     "ds",
	Aliases: []string{"ds-operator"},
	Short:   "Uninstalls the ForgeRock ds-operator",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := delete.GHResource(clientFactory, "ForgeRock/ds-operator", "ds-operator.yaml", tag, true, skipUserDelQ)
		return err
	},
	// Hide this command from help docs
	Hidden:            true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Uninstalls common platform components",
	Long: `
	Uninstalls common platform components`,
	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		clientFactory = factory.NewFactory(deleteFlags)
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	// Install k8s flags
	deleteFlags = initK8sFlags(deleteCmd.PersistentFlags())

	// Delete command-specific flags
	deleteCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Tag/version to parse for delete")
	deleteCmd.PersistentFlags().BoolVarP(&skipUserDelQ, "yes", "y", false, "Do not prompt for confirmation")

	deleteCmd.AddCommand(deleteQuickstart)
	deleteCmd.AddCommand(deleteSecretAgent)
	deleteCmd.AddCommand(deleteDsOperator)

	rootCmd.AddCommand(deleteCmd)
}
