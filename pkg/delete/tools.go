package delete

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
)

// GHResource Uninstalls resources listed in manifests publised on github
func GHResource(clientFactory factory.Factory, ghRepo, fileName, version string, sharedWarn, skipUserQ bool) error {
	fPath := fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", ghRepo, fileName)
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		fPath = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", ghRepo, version, fileName)
	}
	if sharedWarn {
		printer.Warnf("Danger zone: You're about to delete a shared operator which may be required by other deployments in this cluster.")
		printer.Warnf("You normally do not want to delete this if you share this Kubernetes cluster with other users.")
	}
	if err := Manifest(clientFactory, fPath, skipUserQ); err != nil {
		if err == errDidNotAccept {
			return nil
		}
		return err
	}
	return nil
}
