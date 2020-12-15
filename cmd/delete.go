package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/pkg/delete"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var deleteFlags *genericclioptions.ConfigFlags

var deleteQuickstart = &cobra.Command{
	Use:     "quickstart",
	Aliases: []string{"qs"},
	Short:   "Remove the ForgeRock Cloud Deployment Quickstart (CDQ)",
	Long: `
    Remove the ForgeRock Cloud Deployment Quickstart (CDQ):
    * Delete the quickstart deployment
    * Delete all the persistent volumes requested by the CDQ`,
	Example: `
    # Delete the CDQ from the "default" namespace.
    forgeops delete quickstart
    
    # Delete the CDQ from a given namespace.
    forgeops delete quickstart --namespace mynamespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := delete.Quickstart(clientFactory, tag, skipUserConfirmation)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteSecretAgent = &cobra.Command{
	Use:     "secret-agent",
	Aliases: []string{"sa"},
	Short:   "Remove the ForgeRock Secret Agent",
	Long: `
    Remove the ForgeRock secret-agent:
    * Delete the secret-agent deployment`,
	Example: `
    # Delete the secret-agent from the cluster.
    forgeops delete secret-agent`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := delete.GHResource(clientFactory, "ForgeRock/secret-agent", "secret-agent.yaml", tag, true, skipUserConfirmation)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteDsOperator = &cobra.Command{
	Use:     "ds",
	Aliases: []string{"ds-operator"},
	Short:   "Remove the ForgeRock DS operator",
	Long: `
    Remove the ForgeRock ds-operator:
    * Delete the ds-operator deployment`,
	Example: `
    # Delete the ds-operator from the cluster.
    forgeops delete ds`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := delete.GHResource(clientFactory, "ForgeRock/ds-operator", "ds-operator.yaml", tag, true, skipUserConfirmation)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove common platform components",
	Long: `
    Remove common platform components`,
	Example: `
    # Delete the CDQ from the "default" namespace.
    forgeops delete quickstart
    
    # Delete the CDQ from a given namespace.
    forgeops delete quickstart --namespace mynamespace
    
    # Delete the secret-agent from the cluster.
    forgeops delete secret-agent`,
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
	deleteCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Release tag  of the component to be deployed")
	deleteCmd.PersistentFlags().BoolVarP(&skipUserConfirmation, "yes", "y", false, "Do not prompt for confirmation")

	deleteCmd.AddCommand(deleteQuickstart)
	deleteCmd.AddCommand(deleteSecretAgent)
	deleteCmd.AddCommand(deleteDsOperator)

	rootCmd.AddCommand(deleteCmd)
}
