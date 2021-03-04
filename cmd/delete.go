package cmd

import (
	"fmt"

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
		err := delete.Quickstart(clientFactory, "ForgeRock/forgeops", tag, skipUserConfirmation)
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

func generateFRComponentDeleteCommands() {
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
			Short:   fmt.Sprintf("Delete the ForgeRock %[1]s", componentName),
			Long: fmt.Sprintf(`
            Delete the ForgeRock Identity Platform %[1]s:
            * Delete the ForgeRock Identity Platform %[1]q
            * Use --tag to specify a different version to delete`, componentName),
			Example: fmt.Sprintf(`
            # Delete the ForgeRock %[1]q in the default namespace.
            forgeops delete %[1]s
            # Delete the ForgeRock %[1]q in a given namespace.
            forgeops delete %[1]s --namespace mynamespace`, componentName),
			RunE: func(cmd *cobra.Command, args []string) error {
				err := delete.ForgeRockComponent(clientFactory, "ForgeRock/forgeops", componentProperty.artifactName, tag, skipUserConfirmation)
				return err
			},
			Hidden:            componentProperty.hidden,
			SilenceUsage:      true,
			DisableAutoGenTag: true,
		}
	}
	for componentName, componentProperty := range componentList {
		cmd := newCmd(componentName, componentProperty)
		deleteCmd.AddCommand(cmd)
		if componentName == "base" {
			cmd.PersistentFlags().StringVar(&fqdn, "fqdn", "", "FQDN used in the deployment. (default \"[NAMESPACE].iam.example.com\")")
		}
	}
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
	generateFRComponentDeleteCommands()
	rootCmd.AddCommand(deleteCmd)
}
