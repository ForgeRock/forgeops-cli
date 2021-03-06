## forgeops status

Diagnose common cluster and platform deployments

### Synopsis


		Diagnose issues related to running and deploying the ForgeRock platform.
		

```
forgeops status [flags]
```

### Examples

```

		# run all health checks
		forgeops status
		
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
  -h, --help                           help for status
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
      --password string                Password for basic authentication to the API server
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use
      --username string                Username for basic authentication to the API server
```

### Options inherited from parent commands

```
      --log-level string   (options: none|debug|info|warn|error) log statement level. When output=text and level is not 'none' the level is debug (default "none")
  -o, --output string      (options: text|json) command output type. Type json is intended for use in scripting, text is for interactive usage. Not all commands provide both types of output (default "text")
```

### SEE ALSO

* [forgeops](forgeops.md)	 - forgeops is a tool for managing ForgeRock Identity Platform deployments
* [forgeops status operators](forgeops_status_operators.md)	 - Verify that operators are installed and ready
* [forgeops status platform](forgeops_status_platform.md)	 - Verify that operators are installed and ready

