# go-library is a monorepo of modules

The go-library is a monorepo of golang modules designed for reusability within CLI and daemon applications.

## Conventions

A package should be self contained within its own folder off the root of the repository.

For example:
```
mod1/
    pkg/
    internal/
mod2/
    pkg/
    internal/
    go.mod
```

* All modules MUST have a `pkg` directory with the exported interfaces and functions
* All modules MAY have an `internal` directory to organize interfaces and functions not meant for external use
* All modules MUST NOT expose a `main` function
* Modules MAY have their own `go.mod`.  It may need its own `go.mod` if
    - The module needs an isolated release cycle
    - The number of dependencies is causing a negative impact for other modules in the monorepo.


## Release

The monorepo should be release following semantic versioning and golang standards.  e.g. v0.1.0

If a module contains its own `go.mod` it may follow an individual release cycle using the module path and version.

For example:

1. mod1/v0.1.0
1. mod2/v0.1.0
