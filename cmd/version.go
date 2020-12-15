package cmd

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/pkg/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the build information",
	Long: `
    Print the build information.
    Please provide the output of this command when reporting issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
