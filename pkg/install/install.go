package install

import (
	"fmt"
	"strings"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/pkg/version"
	"k8s.io/apimachinery/pkg/api/meta"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/cli-runtime/pkg/resource"
)

// TransformInfoFunc is used to modify the resource.Info
type TransformInfoFunc func(*resource.Info) error

// Manifest obtains the manifest from the given path or URL and applies it in the namespace provided
func Manifest(clientFactory factory.Factory, path string, transformFunctions ...TransformInfoFunc) error {
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	infos, err := k8sCntMgr.GetObjectsFromPath(path)
	if err != nil {
		return err
	}
	return Resources(clientFactory, infos, transformFunctions...)
}

// ManifestStr Applies the given manifest in the namespace provided
func ManifestStr(clientFactory factory.Factory, manifestContents string, transformFunctions ...TransformInfoFunc) error {
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	infos, err := k8sCntMgr.GetObjectsFromStream(strings.NewReader(manifestContents))
	if err != nil {
		return err
	}
	return Resources(clientFactory, infos, transformFunctions...)
}

// Resources applies the resources provided
func Resources(clientFactory factory.Factory, infos []*resource.Info, transformFunctions ...TransformInfoFunc) error {
	errs := []error{}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	if len(infos) == 0 {
		return fmt.Errorf("no objects found")
	}
	for _, tf := range transformFunctions {
		for _, info := range infos {
			if err := tf(info); err != nil {
				return err
			}
		}
	}
	// Iterate through all objects, applying each one.
	for _, info := range infos {
		if err := k8sCntMgr.ApplyObject(info); err != nil {
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

// Provides a set of standard transforms applied to resource.Info objects
func standardTransforms() []TransformInfoFunc {

	manageLabels := func(info *resource.Info) error {
		var metadataAccessor = meta.NewAccessor()
		labels, err := metadataAccessor.Labels(info.Object)
		if err != nil {
			return err
		}
		if labels == nil {
			labels = make(map[string]string)
		}
		labels["forgeops-cli.forgerock.com/version"] = version.Version
		metadataAccessor.SetLabels(info.Object, labels)
		return nil
	}

	return []TransformInfoFunc{manageLabels}

}
