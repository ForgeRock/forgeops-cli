package install

import (
	"fmt"
	"strings"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/internal/utils"
	"github.com/ForgeRock/forgeops-cli/pkg/delete"
	"github.com/ForgeRock/forgeops-cli/pkg/doctor"
	"github.com/ForgeRock/forgeops-cli/pkg/get"
	"github.com/ForgeRock/forgeops-cli/pkg/health"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

// ForgeRockComponent Installs the given component in the namespace provided
func ForgeRockComponent(clientFactory factory.Factory, ghRepo, fileName, version, fqdn string) error {
	fPath := fmt.Sprintf("https://github.com/%s/releases/latest/download/%s", ghRepo, fileName)
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		fPath = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", ghRepo, version, fileName)
	}
	config := qsConfig{
		placeholderFQDN:      "default.iam.example.com",
		placeholderNamespace: "default",
	}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	if len(fqdn) == 0 {
		fqdn = fmt.Sprintf("%s.iam.example.com", ns)
	}

	if strings.Contains(fileName, "base") || strings.Contains(fileName, "quickstart") {
		if err := checkDependencies(clientFactory, doctor.SecretAgentOperatorHealth); err != nil {
			return err
		}
	}
	if strings.Contains(fileName, "ds") || strings.Contains(fileName, "quickstart") {
		if err := checkDependencies(clientFactory, doctor.DSOperatorHealth); err != nil {
			return err
		}
	}
	printer.NoticeHif("Targeting namespace: %q", ns)
	printer.NoticeHif("Installing %q from %q version: %q ", fileName, ghRepo, version)
	manifestStr, err := utils.DownloadTextFile(fPath)
	if err != nil {
		return err
	}

	manifestStr = strings.ReplaceAll(manifestStr, config.placeholderFQDN, fqdn)
	manifestStr = strings.ReplaceAll(manifestStr, "namespace: "+config.placeholderNamespace, "namespace: "+ns)
	if err := ManifestStr(clientFactory, manifestStr, standardTransforms()...); err != nil {
		return err
	}
	printer.Noticef("Installed %q from %q version: %q ", fileName, ghRepo, version)
	return nil

}

// Quickstart Installs the quickstart in the namespace provided
func Quickstart(clientFactory factory.Factory, ghRepo, version, fqdn string) error {

	gvrDeployment := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvrJob := schema.GroupVersionResource{Group: "batch", Version: "v1", Resource: "jobs"}
	gvrStatefulsets := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "statefulsets"}

	// TODO: We should obtain settings like these from a config that can be ingested at runtime.
	// Storing these here for now until we have a solution
	// https://github.com/ForgeRock/forgeops-cli/issues/58
	config := qsConfig{
		placeholderFQDN:      "default.iam.example.com",
		placeholderNamespace: "default",
		importantSecrets: []qsSecret{
			{
				secretName: "am-env-secrets",
				keyName:    []string{"AM_PASSWORDS_AMADMIN_CLEAR"},
				printName:  []string{"amadmin user"},
			},
			{
				secretName: "ds-passwords",
				keyName:    []string{"dirmanager.pw"},
				printName:  []string{"uid=admin user"},
			},
			{
				secretName: "rcs-agent-env-secrets",
				keyName:    []string{"AGENT_IDM_SECRET", "AGENT_RCS_SECRET"},
				printName:  []string{"rcs-agent IDM secret", "rcs-agent RCS secret"},
			},
		},
	}
	// BEGIN TIERED DEPLOYMENT
	// DEPLOY BASE
	if err := ForgeRockComponent(clientFactory, ghRepo, "base.yaml", version, fqdn); err != nil {
		return err
	}
	printer.Noticef("Waiting for secrets to be generated")
	if err := waitForSecrets(clientFactory, config.importantSecrets); err != nil {
		return err
	}
	printer.Noticef("Waiting for git-server to become available")
	if _, err := waitForCondition(clientFactory, gvrDeployment, "git-server", "status.availableReplicas>=1", 120); err != nil {
		return err
	}
	// DEPLOY DS
	if err := ForgeRockComponent(clientFactory, ghRepo, "ds.yaml", version, fqdn); err != nil {
		return err
	}
	printer.Noticef("Waiting for DS deployment to become available. This can take several minutes")
	if _, err := waitForCondition(clientFactory, gvrStatefulsets, "ds-idrepo",
		"status.readyReplicas==spec.replicas", 600); err != nil {
		return err
	}
	// DEPLOY APPS
	if err := ForgeRockComponent(clientFactory, ghRepo, "apps.yaml", version, fqdn); err != nil {
		return err
	}
	printer.Noticef("Waiting for AM deployment to become available. This can take several minutes")
	if _, err := waitForCondition(clientFactory, gvrDeployment, "am", "status.availableReplicas>=1", 600); err != nil {
		return err
	}
	printer.Noticef("Waiting for amster job to complete. This can take several minutess")
	if _, err := waitForCondition(clientFactory, gvrJob, "amster", "status.succeeded>=1", 300); err != nil {
		return err
	}
	// DELETE AMSTER
	if err := delete.ForgeRockComponent(clientFactory, ghRepo, "amster.yaml", version, true); err != nil {
		return err
	}
	// DEPLOY UI
	if err := ForgeRockComponent(clientFactory, ghRepo, "ui.yaml", version, fqdn); err != nil {
		return err
	}
	// END TIERED DEPLOYMENT

	if err := get.Secrets(clientFactory); err != nil {
		return err
	}
	if err := get.URLs(clientFactory, "forgerock"); err != nil {
		return err
	}
	printer.Noticef("CDQ Deployment Complete. Enjoy!")
	return nil
}

func checkDependencies(clientFactory factory.Factory, hlthCheck []byte) error {
	hlth, err := health.GetHealthFromBytes(hlthCheck)
	if err != nil {
		return err
	}
	operAllHealthy, err := health.Run(clientFactory, hlth, true)
	if !operAllHealthy {
		return health.ErrNotAllHealthy
	}
	return err
}

func waitForSecrets(clientFactory factory.Factory, importantSecrets []qsSecret) error {
	// kubectl api-resources -o wide
	gvr := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	for _, secret := range importantSecrets {
		if _, err := k8sCntMgr.WaitForResource(30, ns, secret.secretName, gvr); err != nil {
			return err
		}
	}
	return nil
}

func waitForCondition(clientFactory factory.Factory, gvr schema.GroupVersionResource, name, expr string, timeout int) (bool, error) {
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return false, err
	}
	met, err := k8sCntMgr.WatchEventsForCondition(timeout, ns, name, gvr, k8s.ConditionExpression(expr))
	return met, nil
}
