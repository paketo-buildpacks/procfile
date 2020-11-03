# `gcr.io/paketo-buildpacks/procfile`
The Procfile Buildpack is a Cloud Native Buildpack that turns the contents of a Procfile into process types.

## Behavior
This buildpack will participate if all the following conditions are met

* The application contains a `Procfile`

The buildpack will do the following:

* Contribute the process types from `Procfile` to the image.
  * If the application's stack is `io.paketo.stacks.tiny` the contents of the `Procfile` must be single commands with arguments (no inline environment variable declarations, environment variable references, or joined processes) and will be shellparsed to run in `direct` mode.
  * If the application's stack is not `io.paketo.stacks.tiny` the contents of `Procfile` have no limiataions and will be executed in a shell as-is.

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
