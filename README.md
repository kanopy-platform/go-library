# go-library is a monorepo of modules

The go-library is a monorepo of golang modules designed for reusability within CLI and daemon applications.

## Conventions

A module should be self contained within its own folder off the root of the repository and contain its own `go.mod`.

For example:
```
mod1/
    pkg
    go.mod
mod2/
    pkg
    go.mod
```

* All modules MUST have a `pkg` directory with the exported interfaces and functions
* All modules MAY have an `internal` directory to organize interfaces and functions not meant for external use
* All modules MUST NOT expose a `main` function


## Release

Modules are individually released using semantic versioning.  The module path + version is used as the release name.

For example:

1. mod1/v0.1.0
1. mod2/v0.1.0
