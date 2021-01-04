package cmd

import (
	"context"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/pkg/doctor"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

// cmd globals config
var kubeConfig string
var overrides = &clientcmd.ConfigOverrides{}
var ctx context.Context
var doctorFlags *genericclioptions.ConfigFlags

var ignoreProducts = []string{"ig"}

var ds = &cobra.Command{
	Use:   "directoryserver",
	Short: "Check the status of Directory Server deployment",
	Long: `
	Check the status of Directory Server deployment by checking ready state and configuration.
	  * check workload state
	  * should we call an endpoint/ldap?
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printer.Warnln("Not implemented")
		return nil
	},
	DisableAutoGenTag: true,
	SilenceErrors:     true, //We format and print errors ourselves during Execute().
}

var platform = &cobra.Command{
	Use:     "platform",
	Short:   "Check the status of platform deployment",
	Aliases: []string{"ds"},
	Long: `
	Check the status of platform deployment by checking ready state and configuration.
		* check secrets deployed - should we check for backups?
		* check configs deployed
		* check DS deployment - check backups?
		* check AM deployment - all "Ready" - any other checks e.g. curl?
		* amster? completed - and date?
		* check IDM
		* IG?
	`,
	DisableAutoGenTag: true,
	SilenceErrors:     true, //We format and print errors ourselves during Execute().
}

// Operators
var ignoreOperators []string

var operators = &cobra.Command{
	Use:     "operator",
	Aliases: []string{"op"},
	Short:   "Check Operators Installed and Running",
	Long: `
	Checks to ensure that required operators are installed and ready.
	Searches all namespaces for the default deployment of secret agent, nginx-ingress, cert-manager
	Checks for a minimum ready count of one.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := clientFactory.StaticClient()
		if err != nil {
			return err
		}
		err = doctor.CheckOperators(ctx, client)
		if err != nil {
			return err
		}
		return nil
	},
	DisableAutoGenTag: true,
	SilenceErrors:     true, //We format and print errors ourselves during Execute().
}

var doctorCmd = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{"dr"},
	Short:   "Diagnose common cluster and platform deployments",
	Long: `
	Diagnose common cluster and platform deployments
    `,
	// Configure Client Mgr for all subcommands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		clientFactory = factory.NewFactory(doctorFlags)
	},
	DisableAutoGenTag: true,
	SilenceErrors:     true, //We format and print errors ourselves during Execute().
}

func init() {
	ctx = context.Background()

	// Install k8s flags
	doctorFlags = initK8sFlags(doctorCmd.PersistentFlags())

	// operators
	operators.LocalFlags().StringSlice("ignore-operators", ignoreOperators, "comma seperated list of operators that should ignored during checks")

	//	platform
	platform.LocalFlags().StringSlice("ignore-products", ignoreProducts, "comma seperated list of products that should ignored during checks")
	platform.AddCommand(ds)

	// module command
	doctorCmd.AddCommand(operators)
	doctorCmd.AddCommand(platform)

	// root command
	rootCmd.AddCommand(doctorCmd)
}
