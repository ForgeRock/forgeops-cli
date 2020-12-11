package delete

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/resource"
)

var errDidNotAccept = errors.New("Did not accept prompt to delete")

// Manifest Deletes the resources listed in the manifest provided
func Manifest(clientFactory factory.Factory, path string, skipUserQ bool) error {
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	infos, err := k8sCntMgr.GetObjectsFromPath(path)
	if err != nil {
		return err
	}
	return Resources(clientFactory, infos, skipUserQ)
}

// Resources delete the resources provided
func Resources(clientFactory factory.Factory, infos []*resource.Info, skipUserQ bool) error {
	errs := []error{}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	cfg, err := k8sCntMgr.GetConfigFlags()
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("no objects found")
	}
	accepted, err := askForConfirmation(skipUserQ, infos)
	if err != nil {
		return err
	}
	if !accepted {
		return errDidNotAccept
	}
	// Iterate through all objects, applying each one.
	for _, info := range infos {
		if err := k8sCntMgr.DeleteObjectFromOtherNamespace(info, *cfg.Namespace); err != nil {
			errs = append(errs, err)
		}
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

func askForConfirmation(skipUserQ bool, infos []*resource.Info) (bool, error) {

	if skipUserQ {
		return true, nil
	}
	for _, info := range infos {
		printer.Noticef("Deleting: %s/%s", info.Object.GetObjectKind().GroupVersionKind().Kind, info.Name)
	}
	scanner := bufio.NewScanner(os.Stdin)
	printer.Printf("Do you want to continue? [Y/n]")
	if ok := scanner.Scan(); !ok {
		return false, scanner.Err()
	}
	text := scanner.Text()
	switch strings.ToLower(text) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}
