package cmd

import (
	"os"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cliflag "k8s.io/component-base/cli/flag"
)

var cfgFile string
var kubeConfigFlags *genericclioptions.ConfigFlags
var clientFactory factory.Factory
var tag string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "forgeops",
	Short: "forgeops is a tool for managing ForgeRock platform deployments",
	Long: `
	This tool helps deploying the ForgeRock platform, debug common issues, and validate environments.
    `,
}

// Doc Generate Documents
func Doc() {
	doc.GenMarkdownTree(rootCmd, "./docs")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printer.Errorln(err.Error())
		os.Exit(1)
	}
}

func init() {
	// Install k8s flags
	flags := installCmd.PersistentFlags()
	flags.SetNormalizeFunc(cliflag.WarnWordSepNormalizeFunc) // Warn for "_" flags
	// Normalize all flags coming from other packages. a.k.a. change all "_" to "-". e.g. glog package
	flags.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	kubeConfigFlags = genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	kubeConfigFlags.AddFlags(flags)
}
