# `gcr.io/paketo-buildpacks/procfile`
The Procfile Buildpack is a Cloud Native Buildpack that turns the contents of a Procfile into process types.

## Behavior
This buildpack will participate if one or all of the following conditions are met:

* The application contains a `Procfile`
* A Binding exists with type `Procfile` and secret containing a `Procfile`
* The `BP_PROCFILE_DEFAULT_PROCESS` environment variable is set to a non-empty value

The buildpack will do the following:

* When `BP_PROCFILE_DEFAULT_PROCESS` is set, it will contribute the `web` process type to the image.
* Contribute the process types from one or both `Procfile` files to the image.
  * If process types are identified from both Binding _and_ file, the contents are merged into a single `Procfile`. Commands from the Binding take precedence if there are duplicate types.
  * If process types are identified from environment _and_ Binding _or_ file, the contents are merged into a single `Procfile`. Commands from Binding or file take precedence if there are duplicate types, with Binding taking precedence over file.
  * If the application's stack is `io.paketo.stacks.tiny` the contents of the `Procfile` must be single command with zero or more space delimited arguments. Argument values containing whitespace should be quoted. The resulting process will be executed directly and will not be parsed by the shell.
  * If the application's stack is not `io.paketo.stacks.tiny` the contents of `Procfile` will be executed as a shell script.
* If `BP_DIRECT_PROCESS` is set to `true`, the command will not be executed within a shell.
  * This behavior will become the default with the next major version, fulfilling [RFC-0093](https://github.com/buildpacks/rfcs/blob/main/text/0093-remove-shell-processes.md). Afterwards, this option will be deprecated and removed eventually.

The `BP_DIRECT_PROCESS` environment variable can be used to opt-in in starting processes directly. The next major version of this buildpack will no longer support indirect processes and all processes will be started directly. Once processes are no longer started indirectly by default, the configuration `BP_DIRECT_PROCESS` will be removed since it will have no effect.

## Bindings

The buildpack optionally accepts the following bindings:

|Key                   | Type | Value   | Description
|----------------------|------|---------|------------
|`Procfile` |`Procfile` |List of`<process-type>: <command>` entries | The entries from this Binding will be merged with those from the application's `Procfile`, if both are present. The commands from this Binding take precedence over the application's `Procfile` if there are duplicate process-types.



## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

