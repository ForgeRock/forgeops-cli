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
	Short:   "Deploy the ForgeRock Cloud Deployment Quickstart (CDQ)",
	Long: `
    Deploy the ForgeRock Cloud Deployment Quickstart (CDQ):
    * Apply the latest quickstart manifest
    * use --tag to specify a different CDQ version to deploy`,
	Example: `
      # Deploy the "latest" CDQ in the "default" namespace.
      forgeops apply quickstart
    
      # Deploy the CDQ in the "default" namespace.
      forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci
      
      # Deploy the CDQ in a given namespace.
      forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := apply.Quickstart(clientFactory, tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var secretAgent = &cobra.Command{
	Use:     "secret-agent",
	Aliases: []string{"sa"},
	Short:   "Deploy the ForgeRock Secret Agent",
	Long: `
    Deploy the ForgeRock secret-agent:
    * Apply the latest secret-agent manifest
    * use --tag to specify a different secret-agent version to deploy`,
	Example: `
      # Deploy the "latest" secret-agent.
      forgeops apply sa

      # Deploy a specific version of the secret-agent.
      forgeops apply sa --tag v0.2.1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := apply.GHResource(clientFactory, "ForgeRock/secret-agent", "secret-agent.yaml", tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var dsOperator = &cobra.Command{
	Use:     "ds",
	Aliases: []string{"ds-operator"},
	Short:   "Deploy the ForgeRock DS operator",
	Long: `
    Deploy the ForgeRock ds-operator:
    * Apply the latest ds-operator manifest
    * use --tag to specify a different ds-operator version to deploy`,
	Example: `
      # Deploy the "latest" ds-operator.
      forgeops apply ds

      # Deploy a specific version of the ds-operator.
      forgeops apply ds --tag v0.0.4`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := apply.GHResource(clientFactory, "ForgeRock/ds-operator", "ds-operator.yaml", tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Deploy common platform components",
	Long: `
	Deploy common platform components`,
	Example: `
    # Deploy the "latest" ds-operator.
    forgeops apply ds

    # Deploy the "latest" secret-agent.
    forgeops apply sa

    # Deploy the CDQ in a given namespace.
    forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace`,
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
	applyCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Release tag  of the component to be deployed")

	applyCmd.AddCommand(quickstart)
	applyCmd.AddCommand(secretAgent)
	applyCmd.AddCommand(dsOperator)

	rootCmd.AddCommand(applyCmd)
}
