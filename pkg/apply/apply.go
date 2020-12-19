package apply

import (
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/resource"
)

// TransformInfoFunc is used to modify the resource.Info
type TransformInfoFunc func(*resource.Info) (*resource.Info, error)

// Manifest Installs the quickstart in the namespace provided
func Manifest(clientFactory factory.Factory, path string, transformFunctions ...TransformInfoFunc) error {
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	infos, err := k8sCntMgr.GetObjectsFromPath(path)
	if err != nil {
		return err
	}
	return Resources(clientFactory, infos, transformFunctions...)
}

// Resources applies the resources provided
func Resources(clientFactory factory.Factory, infos []*resource.Info, transformFunctions ...TransformInfoFunc) error {
	errs := []error{}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	cfg, err := k8sCntMgr.GetOverrideFlags()
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("no objects found")
	}
	for _, tf := range transformFunctions {
		for _, info := range infos {
			if _, err := tf(info); err != nil {
				return err
			}
		}
	}
	// Iterate through all objects, applying each one.
	for _, info := range infos {
		if err := k8sCntMgr.ApplyObjectInOtherNamespace(info, *cfg.Namespace); err != nil {
			errs = append(errs, err)
		}
	}
	// If any errors occurred during apply, then return error (or
	// aggregate of errors).
	if len(errs) == 1 {
		return errs[0]
	}
	if len(errs) > 1 {
		return utilerrors.NewAggregate(errs)
	}
	return nil
}
