package health

import (
	"sigs.k8s.io/yaml"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
)

var healthobj = []byte(`
---
kind: health
version: v1alpha
metadata:
  name: forgeops-platform-config
spec:
  resources:
    - resource: secretagentconfigurations
      name: forgerock-sac
      apiversion: v1alpha1
      group: secret-agent.secrets.forgerock.io
      checks:
        - expression: status.totalManagedObjects == 12
          timeout: 0s
        - expression: status.state == "Completed"
          timeout: 0s
`)

// HealthFromBytes deserialize from bytes
func HealthFromBytes(hbytes []byte) (*Health, error) {
	hlth := &Health{}
	err := yaml.Unmarshal(healthobj, hlth)
	if err != nil {
		return &Health{}, err
	}
	return hlth, nil
}

// Run complete a check on a health object - for CLI based use
func Run(clientFactory factory.Factory, hlth *Health) error {
	clientMgr := k8s.NewK8sClientMgr(clientFactory)
	allHealthy, err := hlth.CheckResources(clientMgr)
	if err != nil {
		return err
	}
	if !allHealthy {
		numHealthy := len(hlth.healthy)
		totalNum := len(hlth.Spec.Resources)
		printer.Warnf("Health check %s has %d / %d healthy resources", hlth.Name, numHealthy, totalNum)
		for _, resourceName := range hlth.healthy {
			printer.Noticef("Resource %s healthy", resourceName)

		}
		for _, resourceName := range hlth.unhealthy {
			printer.Warnf("Resource %s is not healthy", resourceName)
		}
		return nil
	}
	printer.Noticef("Health check %s has passed", hlth.Metadata.Name)
	for _, resourceName := range hlth.healthy {
		printer.Noticef("Resource %s is healthy", resourceName)
	}
	return nil
}
