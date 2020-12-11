package apply

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
)

// SecretAgent Installs the SecretAgent operator
func SecretAgent(clientFactory factory.Factory, version string) error {
	fpath := "https://github.com/ForgeRock/secret-agent/releases/latest/download/secret-agent.yaml"
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		fpath = fmt.Sprintf("https://github.com/ForgeRock/secret-agent/releases/download/%s/secret-agent.yaml", version)
	}

	printer.Noticef("Installing secret-agent version: %q", version)
	if err := Apply(clientFactory, fpath); err != nil {
		return err
	}
	printer.Noticef("Installed secret-agent version: %q", version)
	return nil
}
