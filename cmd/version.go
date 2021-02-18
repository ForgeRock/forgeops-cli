package cmd

import (
	"github.com/ForgeRock/forgeops-cli/api"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		log := printer.Logger()
		log.Debug().Msgf("gathering version information without output of %s", printer.CommandOut)
		switch printer.CommandOut {
		case printer.OutJson:
			results, _ := api.NewResultFromKeyPair(
				"build_date:", version.BuildDate,
				"git commit:", version.GitCommit,
				"version:", version.Version,
				"go_version:", version.GoVersion,
				"os_arch:", version.OsArch,
			)
			printer.JsonResult("forgeops version", results)
		case printer.OutText:
			printer.Printf("Build Date: %s", version.BuildDate)
			printer.Printf("Git Commit: %s", version.GitCommit)
			printer.Printf("Version: %s", version.Version)
			printer.Printf("Go Version: %s", version.GoVersion)
			printer.Printf("OS / Arch: %s", version.OsArch)
			return nil

		}
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
