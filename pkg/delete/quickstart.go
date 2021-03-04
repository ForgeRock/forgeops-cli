package delete

import (
	"fmt"
	"strings"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/internal/utils"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

type qsSecret struct {
	secretName string
	keyName    []string
	printName  []string
}

type qsConfig struct {
	placeholderFQDN      string
	placeholderNamespace string
	importantSecrets     []qsSecret
}

// ForgeRockComponent Deletes the given component from the namespace provided
func ForgeRockComponent(clientFactory factory.Factory, ghRepo, fileName, version string, skipUserQ bool) error {
	var errs []error
	fPath := fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", ghRepo, fileName)
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		fPath = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", ghRepo, version, fileName)
	}
	config := qsConfig{
		placeholderNamespace: "default",
	}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	printer.NoticeHif("Targeting namespace: %q", ns)
	manifestStr, err := utils.DownloadTextFile(fPath)
	if err != nil {
		return err
	}
	manifestStr = strings.ReplaceAll(manifestStr, "namespace: "+config.placeholderNamespace, "namespace: "+ns)
	// Delete the quickstart resources listed in the manifest
	if err := ManifestStr(clientFactory, manifestStr, skipUserQ); err != nil {
		if err == errDidNotAccept {
			return nil
		}
		errs = append(errs, err)
	}
	if len(errs) == 1 {
		return errs[0]
	}
	if len(errs) > 1 {
		return utilerrors.NewAggregate(errs)
	}
	return nil

}

// Quickstart Installs the quickstart in the namespace provided
func Quickstart(clientFactory factory.Factory, ghRepo, version string, skipUserQ bool) error {
	var errs []error
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	if err := ForgeRockComponent(clientFactory, ghRepo, "quickstart.yaml", version, skipUserQ); err != nil {
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
