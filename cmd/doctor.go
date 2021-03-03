package cmd

import (
	"context"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/pkg/doctor"
	"github.com/ForgeRock/forgeops-cli/pkg/health"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	// cmd globals config
	kubeConfig    string
	overrides     = &clientcmd.ConfigOverrides{}
	ctx           context.Context
	doctorFlags   *genericclioptions.ConfigFlags
	operatorFlags *genericclioptions.ConfigFlags
	allNamespaces bool

	ds = &cobra.Command{
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
		Hidden:            true,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
	}

	platform = &cobra.Command{
		Use:   "platform",
		Short: "Verify that operators are installed and ready",
		Long: `
		Checks that the platform is running.
	    `,
		Example: `
		# validate the platform is running in the current namespace
		forgeops doctor platform
		# validate the platform is running in the "prod" namespace
		forgeops doctor platform -n prod
		`,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		RunE: func(cmd *cobra.Command, args []string) error {
			configHealth, err := health.GetHealthFromBytes(doctor.DefaultConfigCheck)
			if err != nil {
				return err
			}
			platformHlth, err := health.GetHealthFromBytes(doctor.DefaultPlatformHealth)
			if err != nil {
				return err
			}

			_, confErr := health.Run(clientFactory, configHealth, true)
			_, platErr := health.Run(clientFactory, platformHlth, false)
			if confErr != nil && platErr != nil {
				return errors.Wrap(confErr, platErr.Error())

			} else if confErr != nil {
				return confErr
			} else if platErr != nil {
				return platErr
			}
			return err
		},
	}

	operators = &cobra.Command{
		Use:     "operators",
		Aliases: []string{"op", "operator"},
		Short:   "Verify that operators are installed and ready",
		Long: `
	    Checks to ensure that required operators are installed and ready.
	    `,
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		Example: `
		# check for operators in any namespaces
		forgeops doctor operators
		# check for operators in single namespace
		forgeops doctor operators --all-namespaces=false
		`,

		// Configure Client Mgr for all subcommands
		RunE: func(cmd *cobra.Command, args []string) error {
			hlth, err := health.GetHealthFromBytes(doctor.DefaultOperatorHealth)
			if err != nil {
				return err
			}
			_, err = health.Run(clientFactory, hlth, allNamespaces)
			return err
		},
	}

	doctorCmd = &cobra.Command{
		Use:               "doctor",
		Aliases:           []string{"dr"},
		DisableAutoGenTag: true,
		SilenceUsage:      true,
		Short:             "Diagnose common cluster and platform deployments",
		Long: `
		Diagnose issues related to running and deploying the ForgeRock platform.
		`,
		Example: `
		# run all health checks
		forgeops doctor
		`,
		// Configure Client Mgr for all subcommands
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.Parent().PersistentPreRun(cmd.Parent(), args)
			clientFactory = factory.NewFactory(doctorFlags)

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			operatorHlth, err := health.GetHealthFromBytes(doctor.DefaultOperatorHealth)
			if err != nil {
				return err
			}
			platformHlth, err := health.GetHealthFromBytes(doctor.DefaultPlatformHealth)
			if err != nil {
				return err
			}

			_, operErr := health.Run(clientFactory, operatorHlth, true)
			_, platErr := health.Run(clientFactory, platformHlth, false)
			if operErr != nil && platErr != nil {
				return errors.Wrap(operErr, platErr.Error())

			} else if operErr != nil {
				return operErr
			} else if platErr != nil {
				return platErr
			}
			return err
		},
	}
)

func init() {
	ctx = context.Background()

	// Install k8s flags
	doctorFlags = initK8sFlags(doctorCmd.PersistentFlags())

	// operators
	operators.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", true, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")

	//	platform
	platform.AddCommand(ds)

	// module command
	doctorCmd.AddCommand(operators)
	doctorCmd.AddCommand(platform)

	// root command
	rootCmd.AddCommand(doctorCmd)
}
