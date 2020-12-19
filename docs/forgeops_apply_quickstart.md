## forgeops apply quickstart

Deploy the ForgeRock Cloud Deployment Quickstart (CDQ)

### Synopsis


    Deploy the ForgeRock Cloud Deployment Quickstart (CDQ):
    * Apply the latest quickstart manifest
    * use --tag to specify a different CDQ version to deploy

```
forgeops apply quickstart [flags]
```

### Examples

```

      # Deploy the "latest" CDQ in the "default" namespace.
      forgeops apply quickstart
    
      # Deploy the CDQ in the "default" namespace.
      forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci
      
      # Deploy the CDQ in a given namespace.
      forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace
      
      # Deploy the CDQ with a custom FQDN.
      forgeops apply quickstart --tag 2020.10.28-AlSugoDiNoci --namespace mynamespace --fqdn demo.customdomain.com
```

### Options

```
      --fqdn string   FQDN for CDQ deployment. (default "[NAMESPACE].iam.example.com")
  -h, --help          help for quickstart
```

### Options inherited from parent commands

```
      --as string                      Username to impersonate for the operation
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --cache-dir string               Default cache directory (default "/home/jcastillo/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --password string                Password for basic authentication to the API server
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
  -t, --tag string                     Release tag  of the component to be deployed (default "latest")
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server
```

### SEE ALSO

* [forgeops apply](forgeops_apply.md)	 - Deploy common platform components

