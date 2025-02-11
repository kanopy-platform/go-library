# go-library is a monorepo of modules

The go-library is a monorepo of golang packages designed for reusability within CLI and daemon applications.

## Conventions

A package should be self contained within its own folder off the root of the repository.

If a package grows rapidly or versioning complexity arises due to frequent package changes, you may want to add a `go.mod` file to make it a separate Go module that is versioned and tagged separately.

Further, as module complexity or test execution time increases, you may also want to consider moving modules to a separate repo as needed.

For example:
```
pkg1/
    pkg/
    internal/
pkg2/
    cool.go
mod2/
    pkg/
    internal/
    go.mod
```

Packages may be organized following golang standards.  Use a model that best fits the requirements and complexity of the individual package.  Below are some guidelines to consider: 

* All packages MUST define flat `.go` files in the package directory with the exported interfaces and functions. Complex packages MAY also have a `pkg` directory to organize exported sub-packages.
* Packages MAY have an `internal` directory to organize interfaces and functions not meant for external use
* Packages MUST NOT expose a `main` function
* Packages MAY have their own `go.mod`.  It may need its own `go.mod` if
    - The package needs an isolated release cycle
    - The number of dependencies is causing a negative impact for other packages in the monorepo.


## Release

The monorepo should be released following semantic versioning and golang standards.  e.g. v0.1.0

If a package is a module (contains its own `go.mod`) it MUST follow an individual release cycle using the module path and version.

For example:

1. mod1/v0.1.0
1. mod2/v0.1.0
