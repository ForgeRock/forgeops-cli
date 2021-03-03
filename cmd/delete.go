package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/pkg/delete"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var deleteFlags *genericclioptions.ConfigFlags

var deleteQuickstart = &cobra.Command{
	Use:     "quickstart",
	Aliases: []string{"qs"},
	Short:   "Delete the ForgeRock Cloud Deployment Quickstart (CDQ)",
	Long: `
    Delete the ForgeRock Cloud Deployment Quickstart (CDQ):
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
	Short:   "Delete the ForgeRock Secret Agent",
	Long: `
    Delete the ForgeRock secret-agent:
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
	Use:     "ds-operator",
	Aliases: []string{"dso"},
	Short:   "Delete the ForgeRock DS operator",
	Long: `
    Delete the ForgeRock ds-operator:
    * Delete the ds-operator deployment`,
	Example: `
    # Delete the ds-operator from the cluster.
    forgeops delete ds-operator`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := delete.GHResource(clientFactory, "ForgeRock/ds-operator", "ds-operator.yaml", tag, true, skipUserConfirmation)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteForgeopsBase = &cobra.Command{
	Use:     "base",
	Aliases: []string{"fb"},
	Short:   "Delete the ForgeRock base resources",
	Long: `
    Delete the base resources of the ForgeRock cloud deployment:
    * Delete the base resources of ForgeRock cloud deployment`,
	Example: `
      # Delete the base resources from the "default" namespace.
      forgeops delete base

      # Delete the base resources from a given namespace.
      forgeops delete base --namespace mynamespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Noticeln("This command is not implemented yet")
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteForgeopsDirectory = &cobra.Command{
	Use:     "directory",
	Aliases: []string{"fd"},
	Short:   "Delete the ForgeRock DS resources",
	Long: `
    Delete the directory service resources of the ForgeRock cloud deployment:
    * Delete the directory service resources of ForgeRock cloud deployment`,
	Example: `
      # Delete the directory service resources from the "default" namespace.
      forgeops delete directory

      # Delete the directory service resources from a given namespace.
      forgeops delete directory --namespace mynamespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Noticeln("This command is not implemented yet")
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var deleteForgeopsApps = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"fa"},
	Short:   "Delete the ForgeRock apps (AM, IDM, UI)",
	Long: `
    Delete the ForgeRock apps (AM, IDM, UI):
    * Delete the ForgeRock apps`,
	Example: `
      # Delete the ForgeRock apps from the "default" namespace.
      forgeops delete apps

      # Delete the ForgeRock apps from a given namespace.
      forgeops delete apps --namespace mynamespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Noticeln("This command is not implemented yet")
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete common platform components",
	Long: `
    Delete common platform components`,
	Example: `
    # Delete the CDQ from the "default" namespace.
    forgeops delete quickstart
    
    # Delete the CDQ from a given namespace.
    forgeops delete quickstart --namespace mynamespace
    
    # Delete the secret-agent from the cluster.
    forgeops delete secret-agent`,
	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(cmd.Parent(), args)
		clientFactory = factory.NewFactory(deleteFlags)
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	// Install k8s flags
	deleteFlags = initK8sFlags(deleteCmd.PersistentFlags())

	// Delete command-specific flags
	deleteCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "latest", "Release tag of the component to be deleted")
	deleteCmd.PersistentFlags().BoolVarP(&skipUserConfirmation, "yes", "y", false, "Do not prompt for confirmation")

	deleteCmd.AddCommand(deleteQuickstart)
	deleteCmd.AddCommand(deleteSecretAgent)
	deleteCmd.AddCommand(deleteDsOperator)
	deleteCmd.AddCommand(deleteForgeopsBase)
	deleteCmd.AddCommand(deleteForgeopsDirectory)
	deleteCmd.AddCommand(deleteForgeopsApps)

	rootCmd.AddCommand(deleteCmd)
}
