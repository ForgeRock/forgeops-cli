package install

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
)

// Quickstart Installs the quickstart in the namespace provided
func Quickstart(clientFactory factory.Factory, version string) error {
	var err error
	quickstartPath := "https://github.com/ForgeRock/forgeops/releases/latest/download/quickstart.yaml"
	if version != "latest" && version != "" {
		quickstartPath = fmt.Sprintf("https://github.com/ForgeRock/forgeops/releases/download/%s/quickstart.yaml", version)
	}

	if err = checkDependencies(); err != nil {
		return err
	}
	return Install(clientFactory, quickstartPath)
}

// TODO : need to implement checks before applying.
func checkDependencies() error {
	return nil
}
