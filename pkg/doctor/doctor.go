package doctor

var (

	// DefaultOperatorHealth default health definition for operators
	DefaultOperatorHealth = []byte(`
---
kind: health
version: v1alpha
metadata:
  name: forgeops-operators
spec:
  resources:
    - resource: deployments
      name: secret-agent-controller-manager
      apiversion: v1
      group: apps
      checks:
        - expression: status.availableReplicas == 1
          timeout: 0s
    - resource: deployments
      name: ingress-nginx-controller
      apiversion: v1
      group: apps
      checks:
        - expression: status.availableReplicas == 1
          timeout: 0s
    - resource: deployments
      name: cert-manager
      apiversion: v1
      group: apps
      checks:
        - expression: status.availableReplicas == 1
          timeout: 0s
`)
	// DefaultConfigCheck default configuration check
	DefaultConfigCheck = []byte(`
---
kind: health
version: v1alpha
metadata:
  name: forgerock-config-validation
spec:
  resources:
    - resource: secretagentconfigurations
      name: forgerock-sac
      apiversion: v1alpha1
      group: secret-agent.secrets.forgerock.io
      checks:
        - expression: 'spec.appConfig.secretsManager != "none"'
          timeout: 0s
    - resource: configmaps
      name: platform-config
      apiversion: v1
      group: ""
      checks:
        - expression: 'not (data.FQDN contains "example.com")'
          timeout: 0s
`)

	// DefaultPlatformHealth default definition of the platform
	DefaultPlatformHealth = []byte(`
---
kind: health
version: v1alpha
metadata:
  name: forgerock-platform
spec:
  resources:
    - resource: secretagentconfigurations
      name: forgerock-sac
      apiversion: v1alpha1
      group: secret-agent.secrets.forgerock.io
      checks:
        - expression: status.totalManagedObjects == len(spec.secrets)
          timeout: 0s
        - expression: status.state == "Completed"
          timeout: 0s
    - resource: statefulsets
      name: ds-cts
      apiversion: v1
      group: apps
      checks:
        - expression: status.readyReplicas == spec.replicas
          timeout: 0s
    - resource: statefulsets
      name: ds-idrepo
      apiversion: v1
      group: apps
      checks:
        - expression: status.readyReplicas == spec.replicas
          timeout: 0s
    - resource: jobs
      name: amster
      apiversion: v1
      group: batch
      checks:
        - expression: status.succeeded == spec.completions
          timeout: 0s
    - resource: deployments
      name: am
      apiversion: v1
      group: apps
      checks:
        - expression: status.availableReplicas >= 1
          timeout: 0s
    - resource: deployments
      name: admin-ui
      apiversion: v1
      group: apps
      checks:
        - expression: status.availableReplicas == spec.replicas
          timeout: 0s
    - resource: deployments
      name: end-user-ui
      apiversion: v1
      group: apps
      checks:
        - expression: status.availableReplicas == spec.replicas
          timeout: 0s
    - resource: deployments
      name: login-ui
      apiversion: v1
      group: apps
      checks:
        - expression: status.availableReplicas == spec.replicas
          timeout: 0s
    - resource: statefulsets
      name: idm
      apiversion: v1
      group: apps
      checks:
        - expression: status.readyReplicas >= 1
          timeout: 0s
`)
)
