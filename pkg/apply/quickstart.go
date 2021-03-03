package apply

import (
	"fmt"
	"strings"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	"github.com/ForgeRock/forgeops-cli/internal/utils"
	"github.com/ForgeRock/forgeops-cli/pkg/get"
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

// Quickstart Installs the quickstart in the namespace provided
func Quickstart(clientFactory factory.Factory, version, fqdn string) error {
	quickstartPath := "https://github.com/ForgeRock/forgeops/releases/latest/download/quickstart.yaml"
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		quickstartPath = fmt.Sprintf("https://github.com/ForgeRock/forgeops/releases/download/%s/quickstart.yaml", version)
	}

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
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	if len(fqdn) == 0 {
		fqdn = fmt.Sprintf("%s.iam.example.com", ns)
	}
	if err := checkDependencies(); err != nil {
		return err
	}
	printer.NoticeHif("Targeting namespace: %q", ns)
	printer.NoticeHif("Installing CDQ version: %q", version)
	manifestStr, err := utils.DownloadTextFile(quickstartPath)
	if err != nil {
		return err
	}
	manifestStr = strings.ReplaceAll(manifestStr, config.placeholderFQDN, fqdn)
	manifestStr = strings.ReplaceAll(manifestStr, "namespace: "+config.placeholderNamespace, "namespace: "+ns)
	if err := ManifestStr(clientFactory, manifestStr, standardTransforms()...); err != nil {
		return err
	}
	printer.Noticef("Deployed CDQ version: %q", version)
	printer.Noticef("Waiting for secrets to be generated")
	if err := waitForSecrets(clientFactory, config.importantSecrets); err != nil {
		return err
	}
	if err := get.Secrets(clientFactory); err != nil {
		return err
	}
	if err := get.URLs(clientFactory, "forgerock"); err != nil {
		return err
	}
	printer.Noticef("CDQ Deployment Complete. Enjoy!")
	return nil
}

// TODO : need to implement checks before applying.
// Will use the doctor/status command once development is complete.
func checkDependencies() error {
	return nil
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
