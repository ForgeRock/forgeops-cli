package get

import (
	"context"
	"fmt"
	"strings"

	"github.com/ForgeRock/forgeops-cli/api"
	"github.com/ForgeRock/forgeops-cli/internal/factory"
	"github.com/ForgeRock/forgeops-cli/internal/k8s"
	"github.com/ForgeRock/forgeops-cli/internal/printer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// Secrets returns relevant secrets
func Secrets(clientFactory factory.Factory) error {
	// TODO: We should obtain settings like these from a config that can be ingested at runtime.
	// Storing these here for now until we have a solution
	// https://github.com/ForgeRock/forgeops-cli/issues/58
	config := qsConfig{
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
	return printSecret(clientFactory, config.importantSecrets)
}

func printSecret(clientFactory factory.Factory, importantSecrets []qsSecret) error {
	ctx := context.Background()
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	secretKeyPair := []string{}
	errs := []error{}
	sclient, err := k8sCntMgr.Factory().StaticClient()
	if err != nil {
		return err
	}
	if printer.CommandOut == printer.OutText {
		printer.Noticef("Relevant passwords:")
	}
	for _, s := range importantSecrets {
		k8sSecret, err := sclient.CoreV1().Secrets(ns).Get(ctx, s.secretName, metav1.GetOptions{})
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for idx, key := range s.keyName {
			switch printer.CommandOut {
			case printer.OutJson:
				secretKeyPair = append(secretKeyPair, strings.ReplaceAll(s.printName[idx], " ", "_"))
				secretKeyPair = append(secretKeyPair, string(k8sSecret.Data[key]))
			case printer.OutText:
				printer.NoticeHiln(fmt.Sprintf("%s (%s)", string(k8sSecret.Data[key]), s.printName[idx]))
			}
		}
	}
	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}
	if printer.CommandOut == printer.OutJson {
		results, _ := api.NewResultFromKeyPair(secretKeyPair...)
		printer.JsonResult("forgeops secrets", results)
	}
	return nil
}

// URLs returns releval URLs
func URLs(clientFactory factory.Factory, platformIngressName string) error {
	ctx := context.Background()
	k8sCntMgr := k8s.NewK8sClientMgr(clientFactory)
	ns, err := k8sCntMgr.Namespace()
	if err != nil {
		return err
	}
	sclient, err := k8sCntMgr.Factory().StaticClient()
	if err != nil {
		return err
	}

	forgeopsIngress, err := sclient.ExtensionsV1beta1().Ingresses(ns).Get(ctx, platformIngressName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	fqdn := ""
	for ruleIdx, rule := range forgeopsIngress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			if path.Backend.ServiceName == "am" {
				fqdn = forgeopsIngress.Spec.Rules[ruleIdx].Host
				break
			}
		}
		if len(fqdn) > 0 {
			break
		}
	}
	if len(fqdn) == 0 {
		return fmt.Errorf("Could not find the \"am\" rule in the \"%s\" ingress. Review ingress rules", platformIngressName)
	}
	printURLs(fqdn)
	return nil
}

func printURLs(fqdn string) error {
	baseURL := fmt.Sprintf("https://%s/", fqdn)
	printer.Noticef("Relevant URLs:")
	printer.NoticeHif(baseURL + "platform")
	printer.NoticeHif(baseURL + "admin")
	printer.NoticeHif(baseURL + "am")
	printer.NoticeHif(baseURL + "enduser")
	return nil

}
