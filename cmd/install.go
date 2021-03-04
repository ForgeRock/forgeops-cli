package cmd

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
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
		err := install.Quickstart(clientFactory, "ForgeRock/forgeops", tag, fqdn)
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

func generateFRComponentInstallCommands() {
	type componentProperties struct {
		artifactName string
		hidden       bool
		aliases      []string
	}
	var componentList = map[string]componentProperties{
		"base":        {"base.yaml", false, []string{}},
		"directory":   {"ds.yaml", false, []string{"ds"}},
		"apps":        {"apps.yaml", false, []string{}},
		"ui":          {"ui.yaml", false, []string{}},
		"ds-cts":      {"ds-cts.yaml", true, []string{}},
		"ds-idrepo":   {"ds-idrepo.yaml", true, []string{}},
		"am":          {"am.yaml", true, []string{}},
		"amster":      {"amster.yaml", true, []string{}},
		"idm":         {"idm.yaml", true, []string{}},
		"admin-ui":    {"admin-ui.yaml", true, []string{}},
		"end-user-ui": {"end-user-ui.yaml", true, []string{"enduser-ui"}},
		"login-ui":    {"login-ui.yaml", true, []string{}},
		"rcs-agent":   {"rcs-agent.yaml", true, []string{}},
	}
	newCmd := func(componentName string, componentProperty componentProperties) *cobra.Command {
		return &cobra.Command{
			Use:     componentName,
			Aliases: componentProperty.aliases,
			Short:   fmt.Sprintf("Install the ForgeRock %[1]s", componentName),
			Long: fmt.Sprintf(`
            Install the ForgeRock Identity Platform %[1]s:
            * Install the ForgeRock Identity Platform %[1]q
            * Use --tag to specify a different version to install`, componentName),
			Example: fmt.Sprintf(`
            # Install the ForgeRock %[1]q in the default namespace.
            forgeops install %[1]s
            # Install the ForgeRock %[1]q in a given namespace.
            forgeops install %[1]s --namespace mynamespace`, componentName),
			RunE: func(cmd *cobra.Command, args []string) error {
				err := install.ForgeRockComponent(clientFactory, "ForgeRock/forgeops", componentProperty.artifactName, tag, fqdn)
				return err
			},
			Hidden:            componentProperty.hidden,
			SilenceUsage:      true,
			DisableAutoGenTag: true,
		}
	}
	for componentName, componentProperty := range componentList {
		cmd := newCmd(componentName, componentProperty)
		installCmd.AddCommand(cmd)
		if componentName == "base" {
			cmd.PersistentFlags().StringVar(&fqdn, "fqdn", "", "FQDN used in the deployment. (default \"[NAMESPACE].iam.example.com\")")
		}
	}
}

func init() {
	// Install k8s flags
	installFlags = initK8sFlags(installCmd.PersistentFlags())

	// Install command-specific flags
	installCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "latest", "Release tag  of the component to be deployed")
	quickstart.PersistentFlags().StringVar(&fqdn, "fqdn", "", "FQDN used in the deployment. (default \"[NAMESPACE].iam.example.com\")")
	installCmd.AddCommand(quickstart)
	installCmd.AddCommand(secretAgent)
	installCmd.AddCommand(dsOperator)
	generateFRComponentInstallCommands()
	rootCmd.AddCommand(installCmd)
}
