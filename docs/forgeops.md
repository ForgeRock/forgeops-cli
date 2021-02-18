## forgeops

forgeops is a tool for managing ForgeRock Identity Platform deployments

### Synopsis


    This tool helps deploy the ForgeRock platform, debug common issues, and validate environments.

### Options

```
  -h, --help               help for forgeops
      --log-level string   (options: none|debug|info|warn|error) log statement level. When output=text and level is not 'none' the level is debug (default "none")
  -o, --output string      (options: text|json) command output type. Type json is intended for use in scripting, text is for interactive usage. Not all commands provide both types of output (default "text")
```

### SEE ALSO

* [forgeops apply](forgeops_apply.md)	 - Apply common platform components
* [forgeops clean](forgeops_clean.md)	 - Remove any remaining platform components from the given namespace
* [forgeops delete](forgeops_delete.md)	 - Delete common platform components
* [forgeops docs](forgeops_docs.md)	 - Generate docs
* [forgeops doctor](forgeops_doctor.md)	 - Diagnose common cluster and platform deployments
* [forgeops status](forgeops_status.md)	 - Diagnose common cluster and platform deployments
* [forgeops version](forgeops_version.md)	 - Print the build information

