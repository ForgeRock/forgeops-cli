package doctor

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/ForgeRock/forgeops-cli/internal/printer"
)

type operatorCheck struct {
	minRequired int
	readyCount  int
	installed   bool
}

// CheckOperators validate operators are installed and running
func CheckOperators(ctx context.Context, client *kubernetes.Clientset) error {
	operatorReadyReplicas := map[string]*operatorCheck{
		"secret-agent-controller-manager": &operatorCheck{1, 0, false},
		"ingress-nginx-controller":        &operatorCheck{1, 0, false},
		"cert-manager":                    &operatorCheck{1, 0, false}}

	// Check Installed
	deploys, err := client.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, d := range deploys.Items {
		if checkStatus, ok := operatorReadyReplicas[d.Name]; ok {
			checkStatus.readyCount += int(d.Status.ReadyReplicas)
			checkStatus.installed = true
		}
	}
	operatorErrors := false
	for operatorName, status := range operatorReadyReplicas {
		if !status.installed {
			printer.Warnf("Operator %s is not installed", operatorName)
			operatorErrors = true
		} else if status.minRequired > status.readyCount {
			printer.Warnf("Operator %s installed but has %s ready when %s is required", operatorName, fmt.Sprint(status.readyCount), fmt.Sprint(status.minRequired))
			operatorErrors = true
		}
		printer.Noticef("Operator %s installed and has %s replicas ready", operatorName, fmt.Sprint(status.readyCount))
	}
	if operatorErrors {
		printer.Warnf("Not all operators are installed and ready. Please install them and make sure they are healthy.")
		return nil
	}
	printer.Noticeln("All operators found to be installed and ready.")

	// Check Versions?
	return nil
}
