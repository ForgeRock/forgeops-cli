package clean

import (
	"errors"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/pkg/delete"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

var errDidNotAccept = errors.New("Did not accept prompt to delete")

// Clean deletes remaining forgeops resources from a given namespace
func Clean(clientFactory factory.Factory, skipUserQ bool) error {
	errs := []error{}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	printer.NoticeHif("Targeting namespace: %q", ns)
	printer.Warnf("Danger zone: You're about to delete persistent DS data. This action cannot be undone")
	printer.Warnf("Please back up your DS instance before proceeding.")

	// Delete the PVCs
	infos, err := k8sCntMgr.GetObjectsFromServer("pvc", "")
	if err != nil {
		return err
	}
	if err := delete.Resources(clientFactory, infos, skipUserQ); err != nil {
		errs = append(errs, err)
	}
	// If any errors occurred during Delete, then return error (or
	// aggregate of errors).
	if len(errs) == 1 {
		return errs[0]
	}
	if len(errs) > 1 {
		return utilerrors.NewAggregate(errs)
	}
	return nil
}
