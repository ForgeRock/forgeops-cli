package delete

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// Quickstart Installs the quickstart in the namespace provided
func Quickstart(clientFactory factory.Factory, version string, skipUserQ bool) error {
	var errs []error
	fPath := "https://github.com/ForgeRock/forgeops/releases/latest/download/quickstart.yaml"
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		fPath = fmt.Sprintf("https://github.com/ForgeRock/forgeops/releases/download/%s/quickstart.yaml", version)
	}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	printer.NoticeHif("Targeting namespace: %q", ns)

	// Delete the quickstart resources listed in the manifest
	if err := Manifest(clientFactory, fPath, skipUserQ); err != nil {
		if err == errDidNotAccept {
			return nil
		}
		errs = append(errs, err)
	}

	// Delete the PVCs
	infos, err := k8sCntMgr.GetObjectsFromServer("pvc", "")
	if err != nil {
		return err
	}
	if err := Resources(clientFactory, infos, true); err != nil {
		if err == errDidNotAccept {
			return nil
		}
		errs = append(errs, err)
	}

	// Aggregate of errors from Manifests + PVCs
	if len(errs) == 1 {
		return errs[0]
	}
	if len(errs) > 1 {
		return utilerrors.NewAggregate(errs)
	}
	return nil
}
