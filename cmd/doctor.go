package cmd

import (
	"context"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/pkg/doctor"
)

// kubectl config
var kubeConfig string
var overrides = &clientcmd.ConfigOverrides{}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose local and cluster issues",
	Long: `Diagnose local and cluster issues that may be
	preventing the ForgeRock platform from deploy or running properly`,
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()
		client, err := k8s.GetClient(kubeConfig, overrides)
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
	if home := homedir.HomeDir(); home != "" {
		doctorCmd.Flags().StringVar(&kubeConfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		doctorCmd.Flags().StringVar(&kubeConfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
	clientcmd.BindOverrideFlags(overrides, doctorCmd.Flags(), clientcmd.RecommendedConfigOverrideFlags(""))
	rootCmd.AddCommand(doctorCmd)
}
