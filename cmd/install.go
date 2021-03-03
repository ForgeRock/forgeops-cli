package cmd

import (
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/pkg/install"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// cmd globals config
var installFlags *genericclioptions.ConfigFlags
var fqdn string

var quickstart = &cobra.Command{
	Use:     "quickstart",
	Aliases: []string{"qs"},
	Short:   "Install the ForgeRock Cloud Deployment Quickstart (CDQ)",
	Long: `
    Install the ForgeRock Cloud Deployment Quickstart (CDQ):
    * Install the latest quickstart manifest
    * Use --tag to specify a different CDQ version to install`,
	Example: `
      # Install the "latest" CDQ in the "default" namespace.
      forgeops install quickstart
    
      # Install the CDQ in the "default" namespace.
      forgeops install quickstart --tag 2020.10.28-AlSugoDiNoci
      
      # Install the CDQ in a given namespace.
      forgeops install quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace
      
      # Install the CDQ with a custom FQDN.
      forgeops install quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace --fqdn demo.customdomain.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := install.Quickstart(clientFactory, tag, fqdn)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var secretAgent = &cobra.Command{
	Use:     "secret-agent",
	Aliases: []string{"sa"},
	Short:   "Install the ForgeRock Secret Agent",
	Long: `
    Install the ForgeRock secret-agent:
    * Install the latest secret-agent manifest
    * Use --tag to specify a different secret-agent version to install`,
	Example: `
      # Install the "latest" secret-agent.
      forgeops install sa

      # Install a specific version of the secret-agent.
      forgeops install sa --tag v0.2.1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := install.GHResource(clientFactory, "ForgeRock/secret-agent", "secret-agent.yaml", tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var dsOperator = &cobra.Command{
	Use:     "ds-operator",
	Aliases: []string{"dso"},
	Short:   "Install the ForgeRock DS operator",
	Long: `
    Install the ForgeRock ds-operator:
    * Install the latest ds-operator manifest
    * Use --tag to specify a different ds-operator version to install`,
	Example: `
      # Install the "latest" ds-operator.
      forgeops install ds-operator

      # Install a specific version of the ds-operator.
      forgeops install ds-operator --tag v0.0.4`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := install.GHResource(clientFactory, "ForgeRock/ds-operator", "ds-operator.yaml", tag)
		return err
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var forgeopsBase = &cobra.Command{
	Use:     "base",
	Aliases: []string{"fb"},
	Short:   "Install the ForgeRock base resources",
	Long: `
    Install the base resources of the ForgeRock cloud deployment:
    * Install the base resources of ForgeRock cloud deployment
    * Use --tag to specify a different version to install`,
	Example: `
      # Install the base resources listed in the "latest" release of the forgeops repository.
      forgeops install base

      # Install the base resources listed in a specific release of the forgeops repository.
      forgeops install base --tag 2020.10.28-AlSugoDiNoci`,
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
	Short:   "Install the ForgeRock DS resources",
	Long: `
    Install the directory service resources of the ForgeRock cloud deployment:
    * Install the directory service resources of ForgeRock cloud deployment
    * Use --tag to specify a different version to install`,
	Example: `
      # Install the directory service resources listed in the "latest" release of the forgeops repository.
      forgeops install directory

      # Install the directory service resources listed in a specific release of the forgeops repository.
      forgeops install directory --tag 2020.10.28-AlSugoDiNoci`,
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
	Short:   "Install the ForgeRock apps (AM, IDM, UI)",
	Long: `
    Install the ForgeRock apps (AM, IDM, UI):
    * Install the ForgeRock apps
    * Use --tag to specify a different version to install`,
	Example: `
      # Install the ForgeRock apps listed in the "latest" release of the forgeops repository.
      forgeops install apps

      # Install the ForgeRock apps listed in a specific release of the forgeops repository.
      forgeops install apps --tag 2020.10.28-AlSugoDiNoci`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Noticeln("This command is not implemented yet")
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install common platform components",
	Long: `
	Install common platform components`,
	Example: `
    # Install the "latest" ds-operator.
    forgeops install ds-operator

    # Install the "latest" secret-agent.
    forgeops install sa

    # Install the CDQ in a given namespace.
    forgeops install quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace
    
    # Install the CDQ with a custom FQDN.
    forgeops install quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace --fqdn demo.customdomain.com`,
	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Parent().PersistentPreRun(cmd.Parent(), args)
		clientFactory = factory.NewFactory(installFlags)
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	// Install k8s flags
	installFlags = initK8sFlags(installCmd.PersistentFlags())

	// Install command-specific flags
	installCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "latest", "Release tag  of the component to be deployed")
	quickstart.PersistentFlags().StringVar(&fqdn, "fqdn", "", "FQDN used in the deployment. (default \"[NAMESPACE].iam.example.com\")")
	forgeopsBase.PersistentFlags().StringVar(&fqdn, "fqdn", "", "FQDN used in the deployment. (default \"[NAMESPACE].iam.example.com\")")

	installCmd.AddCommand(quickstart)
	installCmd.AddCommand(secretAgent)
	installCmd.AddCommand(dsOperator)
	installCmd.AddCommand(forgeopsBase)
	installCmd.AddCommand(forgeopsDirectory)
	installCmd.AddCommand(forgeopsApps)

	rootCmd.AddCommand(installCmd)
}
