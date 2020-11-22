package cmd

import (
	"os"

	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var cfgFile string

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
