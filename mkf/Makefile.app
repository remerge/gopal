clean::
	-rm -rf .build/*


# Buidling binaries

# local development fallbacks if there are no CI variables
DEV_COMMIT := dev.$(shell whoami).$(shell git rev-parse --short HEAD)
DEV_REPO := $(PACKAGE)
DEV_NUM := $(shell whoami)

CI_COMMIT := $(or $(CIRCLE_SHA1), $(TRAVIS_COMMIT), $(DEV_COMMIT))
CI_REPO := $(or $(CIRCLE_PROJECT_REPONAME), $(TRAVIS_REPO_SLUG), $(DEV_REPO))
CI_NUM := $(or $(CIRCLE_BUILD_NUM), $(TRAVIS_JOB_NUMBER), $(DEV_NUM))

LDIMPORT=github.com/remerge/go-service
LDFLAGS=-X $(LDIMPORT).CodeVersion=$(CI_COMMIT) -X $(LDIMPORT).CodeBuild=$(CI_REPO)\#$(CI_NUM)@$(shell date -u +%FT%TZ)
MAIN ?= main

# NOTE: logically pull out binaries from source tree
.build/$(PROJECT).%: $(GO_SOURCES) go.mod
	mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(basename $*) GOARCH=$(patsubst .%,%,$(suffix $*)) go build -trimpath -o $@ -ldflags "$(LDFLAGS)" $(PACKAGE)/$(MAIN)

dist: .build/$(PROJECT).linux.amd64	## linux amd64 binary

.BIN_LOCAL = $(PROJECT).$(shell go env GOOS).$(shell go env GOARCH)

local: .build/$(.BIN_LOCAL)

# Deployment

release:
	git push origin master:production

deploy:
	/bin/bash -c 'cd ../chef.new && knife ssh roles:$(PROJECT) sudo chef-client'
