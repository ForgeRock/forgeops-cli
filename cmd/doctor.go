package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/pkg/doctor"
)

var kubeConf string

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose local and cluster issues",
	Long: `Diagnose local and cluster issues that may be
	preventing the ForgeRock platform from deploy or running properly`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := k8s.GetClient(kubeConf)
		if err != nil {
			return err
		}
		err = doctor.CheckOperators(ctx, client)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	doctorCmd.Flags().StringVar(&kubeConf, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests.")
	rootCmd.AddCommand(doctorCmd)
}
