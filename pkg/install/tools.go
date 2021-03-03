package install

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
)

// GHResource Installs resources listed in manifests publised on github
func GHResource(clientFactory factory.Factory, ghRepo, fileName, version string) error {
	fPath := fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", ghRepo, fileName)
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		fPath = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", ghRepo, version, fileName)
	}
	printer.Noticef("Installing %q version: %q", ghRepo, version)
	if err := Manifest(clientFactory, fPath, standardTransforms()...); err != nil {
		return err
	}
	printer.Noticef("Installed %q version: %q", ghRepo, version)
	return nil
}
