# go-makefile

Common Makefile for remerge Go projects. Various targets are provided:

*TBD: update while working on: https://trello.com/c/vROA2Va8/115-improve-and-share-makefilecommon-revivetoml-between-go-services*

## Setup

### If the makefile is not yet part of the repository

Add this repository as a subtree to the go project repository.

```
git remote add makefile https://github.com/remerge/go-makefile.git
git subtree add --squash --prefix mkf/ makefile master
```
Afterwards `mkf/Makefile.common` can be included in the parent project.

### If the makefile is already available

`make update-makefile`

## Example

```
REVIVELINTER_EXCLUDES = $(foreach p,$(wildcard **/*_fsm.go),-exclude $(p))

include mkf/Makefile.common
```

## Build

Basically you always build locally using

```
make local
```

Build on the server (using CI most probably) is done using

```
make dist
```

## CI

Speaking about CI [here is example .travis.yml file](https://github.com/remerge/go-makefile/blob/master/travis.yaml)


## Testing

From the box you can use three targets:

- `test` Run tests
- `race` Run tests with race detection. Recommended over `test`
- `bench` Run benchmarks

Behaviour can be altered with following variables:

- `TESTS` Tests pattern
- `TEST_TIMEOUT` Tests timeout
- `TEST_TAGS` Use specific tags for tests delimited by comma

```
$ make test TESTS=TestSomething TEST_TAGS="integration,postgres"
```

Two additional `test-nocache` and `race-nocache` targets are same as `test` and
`race` but always runs each test. It's highly recommends to use `*-nocache` 
targets in CI to detect fragile tests. To maintain thread-safe code 
`race-nocache` is more preferable than `test-nocache`.

## Modules

Use `mod-tidy` to tidy Go modules.

From Go 1.11 modules used instead vendor. To maintain vendor directory use

```
GO111MODULES=on go mod vendor
```

## watch

Watches your go code for changes and constantly rebuild it and run tests

## lint

This target invokes Go format, modules consistency, check, vet and revive.

### vet

To use specific vet flags use `VET_FLAGS` variable.

```
VET_FLAGS = -unsafeptr=false
include mkf/Makefile.common
``` 

### revive

Revive will be installed automatically if it not present on host. To 
override revive config put `revive.toml` in root of build tree. 
Use `REVIVELINTER_EXCLUDES` variable to add excludes.

```
REVIVELINTER_EXCLUDES = $(foreach p,$(wildcard **/*_fsm.go),-exclude $(p))
include mkf/Makefile.common
```

`lint` target guarantees to never change sources.

`lint-mod-outdated` checks all modules are up to date. It's disabled until 
we migrate to upstream pq and sarama.

## Isolated test deployment (diversion)

Sometimes run application on isolated instance is the best way to test it. *This is
not safe regular deployment!* You must warn your colleagues.

### Requirements

1. SSH access and sudo rights on target machine
1. Minimal knowledge about SystemD
1. Go with crossplatform build on developer machine

### Setup and teardown

For each operation you need to define the `DIVERT_SSH` environment variable.

```shell
$ DIVERT_SSH=user@app.machine make ...
$ make ... DIVERT_SSH=user@app.machine
```

Alternatively `DIVERT_SSH` can be defined globally:

```shell
$ export DIVERT_SSH=user@app.machine
```

> Diversion status may be checked by `.CHECK-divert-on` and
  `.CHECK-divert-on` targets.

Now you can run the `divert-setup` target. This target will prepare the basic setup and
stop `chef-client`.

To revert diversion environment and start `chef-client` use `divert-teardown` target.

### Journal

Use `divert-journal` to follow application log. This target is not depends on
diverted environment.

### Deploy

Use `divert-do` target to deploy dev version.
