package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docFormat string
var outputDir string

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate docs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdTree := cmd.Parent()
		os.Mkdir(outputDir, 0755)
		switch docFormat {
		case "man":
			header := &doc.GenManHeader{
				Title:   "FORGEOPS",
				Section: "1",
			}
			err := doc.GenManTree(cmdTree, header, outputDir)
			if err != nil {
				return err
			}
			return nil
		case "md":
			err := doc.GenMarkdownTree(cmdTree, outputDir)
			if err != nil {
				return err
			}
		}
		return nil
	},
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func init() {
	docsCmd.Flags().StringVarP(&docFormat, "doc-type", "t", "md", "output can be md || man")
	docsCmd.Flags().StringVarP(&outputDir, "output-dir", "d", "./docs", "output path docs")
	rootCmd.AddCommand(docsCmd)
}
