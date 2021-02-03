package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/pkg/apply"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var applyFlags *genericclioptions.ConfigFlags
var fqdn string

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
      forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace
      
      # Deploy the CDQ with a custom FQDN.
      forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace --fqdn demo.customdomain.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := apply.Quickstart(clientFactory, tag, fqdn)
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
	Use:     "ds-operator",
	Aliases: []string{"dso"},
	Short:   "Deploy the ForgeRock DS operator",
	Long: `
    Deploy the ForgeRock ds-operator:
    * Apply the latest ds-operator manifest
    * use --tag to specify a different ds-operator version to deploy`,
	Example: `
      # Deploy the "latest" ds-operator.
      forgeops apply ds-operator

      # Deploy a specific version of the ds-operator.
      forgeops apply ds-operator --tag v0.0.4`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := apply.GHResource(clientFactory, "ForgeRock/ds-operator", "ds-operator.yaml", tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var forgeopsBase = &cobra.Command{
	Use:     "base",
	Aliases: []string{"fb"},
	Short:   "Deploy the ForgeRock base resources",
	Long: `
    Deploy the base resources of the ForgeRock cloud deployment:
    * Apply the base resources of ForgeRock cloud deployment
    * use --tag to specify a different version to deploy`,
	Example: `
      # Deploy the base resources listed in the "latest" release of Forgeops.
      forgeops apply base

      # Deploy the base resources listed in a specific release of Forgeops.
      forgeops apply base --tag 2020.10.28-AlSugoDiNoci`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Noticeln("This command is not implemented yet")
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var forgeopsDirectory = &cobra.Command{
	Use:     "directory",
	Aliases: []string{"fd"},
	Short:   "Deploy the ForgeRock DS resources",
	Long: `
    Deploy the directory service resources of the ForgeRock cloud deployment:
    * Apply the directory service resources of ForgeRock cloud deployment
    * use --tag to specify a different version to deploy`,
	Example: `
      # Deploy the directory service resources listed in the "latest" release of Forgeops.
      forgeops apply directory

      # Deploy the directory service resources listed in a specific release of Forgeops.
      forgeops apply directory --tag 2020.10.28-AlSugoDiNoci`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Noticeln("This command is not implemented yet")
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var forgeopsApps = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"fa"},
	Short:   "Deploy the ForgeRock apps (AM, IDM, UI)",
	Long: `
    Deploy the ForgeRock apps (AM, IDM, UI):
    * Apply the ForgeRock apps
    * use --tag to specify a different version to deploy`,
	Example: `
      # Deploy the ForgeRock apps listed in the "latest" release of Forgeops.
      forgeops apply apps

      # Deploy the ForgeRock apps listed in a specific release of Forgeops.
      forgeops apply apps --tag 2020.10.28-AlSugoDiNoci`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Noticeln("This command is not implemented yet")
		return nil
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
    forgeops apply ds-operator

    # Deploy the "latest" secret-agent.
    forgeops apply sa

    # Deploy the CDQ in a given namespace.
    forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace
    
    # Deploy the CDQ with a custom FQDN.
    forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace --fqdn demo.customdomain.com`,
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
	applyCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Release tag  of the component to be deployed (default \"latest\")")
	quickstart.PersistentFlags().StringVar(&fqdn, "fqdn", "", "FQDN for CDQ deployment. (default \"[NAMESPACE].iam.example.com\")")
	forgeopsBase.PersistentFlags().StringVar(&fqdn, "fqdn", "", "FQDN for CDQ deployment. (default \"[NAMESPACE].iam.example.com\")")

	applyCmd.AddCommand(quickstart)
	applyCmd.AddCommand(secretAgent)
	applyCmd.AddCommand(dsOperator)
	applyCmd.AddCommand(forgeopsBase)
	applyCmd.AddCommand(forgeopsDirectory)
	applyCmd.AddCommand(forgeopsApps)

	rootCmd.AddCommand(applyCmd)
}
