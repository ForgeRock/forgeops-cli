package apply

import (
	"context"
	"fmt"

	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type qsSecret struct {
	secretName string
	keyName    []string
	printName  []string
}

// Quickstart Installs the quickstart in the namespace provided
func Quickstart(clientFactory factory.Factory, version string) error {
	quickstartPath := "https://github.com/ForgeRock/forgeops/releases/latest/download/quickstart.yaml"
	if len(version) == 0 {
		version = "latest"
	}
	if version != "latest" {
		quickstartPath = fmt.Sprintf("https://github.com/ForgeRock/forgeops/releases/download/%s/quickstart.yaml", version)
	}

	// TODO: We should obtain settings like these from a config that can be ingested at runtime.
	// Storing these here for now until we have a solution
	importantSecrets := []qsSecret{
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
	}

	if err := checkDependencies(); err != nil {
		return err
	}
	printer.NoticeHif("Installing CDQ version: %q", version)
	if err := Manifest(clientFactory, quickstartPath); err != nil {
		return err
	}
	printer.Noticef("Deployed CDQ version: %q", version)

	printer.Noticef("Waiting for secrets to be generated")
	if err := waitForSecrets(clientFactory, importantSecrets); err != nil {
		return err
	}
	printer.Noticef("Relevant passwords:")
	if err := printSecret(clientFactory, importantSecrets); err != nil {
		return err
	}
	printURLs(clientFactory)
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
	cfg, err := k8sCntMgr.GetConfigFlags()
	if err != nil {
		return err
	}
	for _, secret := range importantSecrets {
		if _, err := k8sCntMgr.WaitForResource(30, *cfg.Namespace, secret.secretName, gvr); err != nil {
			return err
		}
	}
	return nil
}

func printSecret(clientFactory factory.Factory, importantSecrets []qsSecret) error {
	ctx := context.Background()
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	cfg, err := k8sCntMgr.GetConfigFlags()
	if err != nil {
		return err
	}
	sclient, err := k8sCntMgr.StaticClient()
	if err != nil {
		return err
	}
	for _, s := range importantSecrets {
		k8sSecret, err := sclient.CoreV1().Secrets(*cfg.Namespace).Get(ctx, s.secretName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		for idx, key := range s.keyName {
			printer.NoticeHiln(fmt.Sprintf("%s (%s)", string(k8sSecret.Data[key]), s.printName[idx]))
		}
	}
	return nil
}

func printURLs(clientFactory factory.Factory) error {
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	cfg, err := k8sCntMgr.GetConfigFlags()
	if err != nil {
		return err
	}
	fqdn := fmt.Sprintf("https://%s.iam.example.com/", *cfg.Namespace)
	printer.Noticef("Relevant URLs:")
	printer.NoticeHif(fqdn + "platform")
	printer.NoticeHif(fqdn + "admin")
	printer.NoticeHif(fqdn + "am")
	printer.NoticeHif(fqdn + "enduser")
	return nil

}
