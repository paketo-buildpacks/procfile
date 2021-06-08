# `gcr.io/paketo-buildpacks/procfile`
The Procfile Buildpack is a Cloud Native Buildpack that turns the contents of a Procfile into process types.

## Behavior
This buildpack will participate if all the following conditions are met

* The application contains a `Procfile`

The buildpack will do the following:

* Contribute the process types from `Procfile` to the image.
  * If the application's stack is `io.paketo.stacks.tiny` the contents of the `Procfile` must be single command with zero or more space delimited arguments. Argument values containing whitespace should be quoted. The resulting process will be executed directly and will not be parsed by the shell.
  * If the application's stack is not `io.paketo.stacks.tiny` the contents of `Procfile` will be executed as a shell script.

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0

