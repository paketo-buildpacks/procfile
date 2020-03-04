# `procfile`
The Procfile Buildpack is a Cloud Native Buildpack that turns the contents of a Procfile into process types.

## Behavior
This buildpack will participate any of the following conditions are met

* The application contains a `Procfile`

The buildpack will do the following:

* Contribute the process types from `Procfile` to the image.

## License
This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
