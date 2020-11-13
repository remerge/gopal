# go-makefile

Shared Makefile includes ([common, app, divert](#includes)) for remerge Go projects.

> For how to inlcude this in a new project see [Setup](#setup).

Provides the following targets (if not specified - defined in [common](#includes)).

# Building binaries ([app](#includes))

To build a binary matching the local architecture.

```
make local
```

The resulting binary can be found in `.build/`.


To build a binary that matches the architecture of our servers (run by the CI as well).


```
make dist
```

For a specific OS/architecture the hidden target `.build/<app>.<os>.<arch>` can be used.

## Go versions

One of easiest ways to use multiple Go versions is [golang.org/dl](https://github.com/golang/dl):

```
GO111MODULE=off go get golang.org/dl/go1.11.6
go1.11.6 download
PATH=$HOME/sdk/go1.11.6/bin make lint test
```

Go version is not counted by Make targets. Use `make clean` to ensure that 
artifacts are built against correct version.

```
make clean && PATH=$HOME/sdk/go1.11.6/bin make lint test
```

# Generating artifacts

We use `go generate` to generate various artefacts (parsers and more). The `gen` target can be used to trigger the generation.

    make gen


## Testing

From the box you can use three targets:

- `test` Run tests
- `race` Run tests with race detection. Recommended over `test`
- `bench` Run benchmarks

Behaviour can be altered with following variables:

- `TESTS` Tests pattern
- `TEST_TIMEOUT` Tests timeout
- `TEST_TAGS` Use specific tags for tests delimited by comma
- `TEST_ARGS` Additional test arguments

```
$ make test TESTS=TestSomething TEST_TAGS="integration,postgres" TEST_ARGS=-v
```

To run the tests without test caching use `test-nocache` (and `race-nocache`)
It's highly recommends to use the `*-nocache` targets in CI to detect fragile tests.


## watching changes

Watch your Go code for changes, rebuild and run test on change.

    make watch T=NameOfTestWithoutTestPrefix


# Linting

This target lints the source code using various tools: `go fmt`, `goimports`, the modules consistency check, `go check`, `go vet` and `revive`.

    make lint

> The target will not change sources.

`lint-mod-outdated` checks all modules are up to date. It's disabled until 
we migrate to upstream pq and sarama.

## Buildsystem

`lint-mkf` checks that `go-makefile` is up to date. Use `make update-makefile` 
to update if `lint-mkf` target fails. Set `MKF_BRANCH` to test using a 
different `go-makefile` branch. 

If your repository is public and CI is unable to access `go-makefile`, remove 
the `lint-mkf` step in the CI config.

## vet configuration

To use specific vet flags use `VET_FLAGS` variable.

```
VET_FLAGS = -unsafeptr=false
include mkf/Makefile.common
``` 

## revive configuration

Revive will be installed automatically if it us not present. To 
override revive config put `revive.toml` in root of build tree. 
Use `REVIVELINTER_EXCLUDES` variable to add excludes.

```
REVIVELINTER_EXCLUDES = $(foreach p,$(wildcard **/*_fsm.go),-exclude $(p))
include mkf/Makefile.common
```

## gofmt and goimports configuration

Use `GOFMT_EXCLUDES` variable to exclude files from import checking.

```
GOFMT_EXCLUDES = -not -path "./vendor/*"
include mkf/Makefile.common
```

# Deploying ([app](#includes))

We use Chef to deploy binaries to our servers. The git branch that is used to build bianries from is `production`.
If changes to the production branch are detected, the CI executes our test suite and if succesfull builds a new binary.
Usually the new binary is roled out to all servers within a window of 30 minutes.

The `release` target can be used to trigger this process. It pushes the current `master` branch to `production`.

    make release

If the deployment needs to be done faster (enforced) the `deploy` target can be used. **This should only be used if there is a good reason!**

    make deploy


# Modules

To clean update the modules include via `go.mod` you can use the `mod-tidy` target.

    make mod-tidy



# Isolated test deployment ([divert](#includes))

In very rare case a custom build needs to be deployed in production. Diversions make this possible. 
**This should be use carefully and only under special curcumstances**

## Requirements

1. SSH access and sudo rights on target machine
1. Minimal knowledge about SystemD
1. Go with crossplatform build on developer machine

## Setup and teardown

For each operation you need to define the `DIVERT_SSH` environment variable.

```shell
DIVERT_SSH=user@app.machine make ...
make ... DIVERT_SSH=user@app.machine
```

Alternatively `DIVERT_SSH` can be defined globally:

```shell
export DIVERT_SSH=user@app.machine
```

> Diversion status may be checked by `.CHECK-divert-on` and
  `.CHECK-divert-on` targets.

Now you can run the `divert-setup` target. This target will prepare the basic setup and
stop `chef-client`.

To revert diversion environment and start `chef-client` use `divert-teardown` target.

## Journal

Use `divert-journal` to follow application log. This target is not depends on
diverted environment.

### Deploy

Use `divert-do` target to deploy dev version.

# Setup

For a new project that does not use the common Makefile includes yet, add this repository as a subtree to the project repository.

```
git remote add makefile https://github.com/remerge/go-makefile.git
git subtree add --squash --prefix mkf/ makefile master
```

Afterwards `mkf/Makefile.common` can be included in the parent project. If the project is a service that is compiled into a binary `mkf/Makefile.app` and `mkf/Makefile.divert` should be included as well.

#### Includes

* **common** (`Makefile.common`) basic Go targets 
* **app** (`Makefile.app`): common targets for building binaries and deploy them
* **divert** (`Makefile.divert`) for deploying temporarly deploying development binaries to production

### Example top level Makefile

```Makefile
REVIVELINTER_EXCLUDES = $(foreach p,$(wildcard **/*_fsm.go),-exclude $(p))

include mkf/Makefile.common mkf/Makefile.app mkf/Makefile.divert
```

## Updating

To update the Makefile includes in the current repository.

    make update-makefile


## Travis CI configuration

Every project should have a Travis CI configuration. [This example can be used as a starting point.](https://github.com/remerge/go-makefile/blob/master/travis.yaml)

## CircleCI configuration

To setup a project for CirlceCI please read the [Setup CircleCI guide in Confluence](https://remerge.atlassian.net/wiki/spaces/tech/pages/4030889/Creating+a+new+Go+project). This uses the config file in the `.circleci` folder as a starting point.
