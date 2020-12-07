package doctor

import (
	"context"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var sacResource = schema.GroupVersionResource{
	Group:    "secret-agent.secrets.forgerock.io",
	Version:  "v1alpha1",
	Resource: "secretagentconfigurations",
}

// SACStatus state considered ready
var SACStatus = "Completed"

// CheckSAC validate SAC is ready
func CheckSAC(ctx context.Context, namespace string, f factory.Factory) (bool, error) {
	client, err := f.DynamicClient()
	if err != nil {
		return false, err
	}
	var resource dynamic.ResourceInterface
	if namespace != "" {
		resource = client.Resource(sacResource).Namespace(namespace)
	} else {
		resource = client.Resource(sacResource)
	}
	res, err := resource.List(ctx, metav1.ListOptions{})
	allPassed := 0
	numItems := len(res.Items)
	if numItems == 0 {
		printer.Warnf("No secret agent configuration found in namespace: %s", namespace)

	}
	for _, item := range res.Items {
		state, _, err := unstructured.NestedString(item.Object, "status", "state")
		checksFailed := 0
		if err != nil {
			return false, err
		}
		if state != SACStatus {
			checksFailed++
			printer.Warnf("Secret Agent Config: %s in namespace %s in state of %s", item.GetName(), item.GetNamespace(), state)
		} else {
			printer.Noticef("Secret Agent Config: %s in namespace %s in state of %s", item.GetName(), item.GetNamespace(), state)
		}
		secretManagerConfig, found, err := unstructured.NestedString(item.Object, "spec", "appConfig", "secretsManager")
		if err != nil {
			return false, err
		}
		if !found {
			checksFailed++
			printer.Warnf("Secret Agent Config: %s in namespace %s is missing secrets manager setting missing", item.GetName(), item.GetNamespace())
		}
		if secretManagerConfig == "none" {
			checksFailed++
			printer.Warnf("Secret Agent Config: %s in namespace %s is at risk for data loss since no secret manager is set", item.GetName(), item.GetNamespace())
		}
		// TODO should we check cloud creds to exist too?
		if checksFailed == 0 {
			allPassed++
			printer.Noticef("Secret Agent Config: %s in namespace %s passed all checks", item.GetName(), item.GetNamespace())
		} else {
			printer.Warnf("Secret Agent Config: %s in namespace %s failed %d checks", item.GetName(), item.GetNamespace(), checksFailed)
		}
	}
	return allPassed == numItems, nil
}

// CheckConfigMaps validate platform configmaps
func CheckConfigMaps(ctx context.Context, f factory.Factory) (bool, error) {
	return false, nil
}
