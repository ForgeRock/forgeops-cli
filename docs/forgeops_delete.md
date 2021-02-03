## forgeops delete

Delete common platform components

### Synopsis


    Delete common platform components

### Examples

```

    # Delete the CDQ from the "default" namespace.
    forgeops delete quickstart
    
    # Delete the CDQ from a given namespace.
    forgeops delete quickstart --namespace mynamespace
    
    # Delete the secret-agent from the cluster.
    forgeops delete secret-agent
```

### Options

```
      --as string                      Username to impersonate for the operation
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
  -h, --help                           help for delete
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --password string                Password for basic authentication to the API server
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
  -t, --tag string                     Release tag of the component to be deleted (default "latest")
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server
  -y, --yes                            Do not prompt for confirmation
```

### SEE ALSO

* [forgeops](forgeops.md)	 - forgeops is a tool for managing ForgeRock Identity Platform deployments
* [forgeops delete apps](forgeops_delete_apps.md)	 - Delete the ForgeRock apps (AM, IDM, UI)
* [forgeops delete base](forgeops_delete_base.md)	 - Delete the ForgeRock base resources
* [forgeops delete directory](forgeops_delete_directory.md)	 - Delete the ForgeRock DS resources
* [forgeops delete ds-operator](forgeops_delete_ds-operator.md)	 - Delete the ForgeRock DS operator
* [forgeops delete quickstart](forgeops_delete_quickstart.md)	 - Delete the ForgeRock Cloud Deployment Quickstart (CDQ)
* [forgeops delete secret-agent](forgeops_delete_secret-agent.md)	 - Delete the ForgeRock Secret Agent

