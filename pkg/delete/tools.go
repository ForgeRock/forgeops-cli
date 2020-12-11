package delete

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
)

// SecretAgent Installs the SecretAgent operator
func SecretAgent(clientFactory factory.Factory, version string, skipUserQ bool) error {
	fPath := "https://github.com/ForgeRock/secret-agent/releases/latest/download/secret-agent.yaml"
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		fPath = fmt.Sprintf("https://github.com/ForgeRock/secret-agent/releases/download/%s/secret-agent.yaml", version)
	}
	printer.Warnf("Danger zone: You're about to delete a shared operator which may be required by other deployments in this cluster.")
	printer.Warnf("You wouldn't normally want to delete the secret-agent if you share this Kubernetes cluster with other users.")
	if err := Manifest(clientFactory, fPath, skipUserQ); err != nil {
		if err == errDidNotAccept {
			return nil
		}
		return err
	}
	return nil
}
