PROJECT := gopal
PACKAGE := github.com/remerge/$(PROJECT)

# http://stackoverflow.com/questions/322936/common-gnu-makefile-directory-path#comment11704496_324782
TOP := $(dir $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST)))

GOFMT=gofmt -w -s

GOSRCDIR=$(GOPATH)/src/$(PACKAGE)
GOPATHS=$(shell glide novendor | grep -v /main/)
GOFILES=$(shell git ls-files | grep '\.go$$')
MAINGO=$(wildcard main/*.go)
MAIN=$(patsubst main/%.go,%,$(MAINGO))

CODE_VERSION=$(TRAVIS_COMMIT)
ifeq ($(CODE_VERSION),)
	CODE_VERSION=$(shell git rev-parse --short HEAD)-dev
endif

CODE_BUILD=$(TRAVIS_REPO_SLUG)\#$(TRAVIS_JOB_NUMBER)
ifeq ($(CODE_BUILD),\#)
	CODE_BUILD=$(PACKAGE)\#$(shell whoami)
endif

LDIMPORT=$(PACKAGE)/vendor/github.com/remerge/rex
LDFLAGS=-X $(LDIMPORT).CodeVersion=$(CODE_VERSION) -X $(LDIMPORT).CodeBuild=$(CODE_BUILD)@$(shell date -u +%FT%TZ)

.PHONY: build run clean lint test bench fmt dep init up gen release deploy

all: build

build: fmt
	cd $(GOSRCDIR) && \
		CGO_ENABLED=0 \
		go build -v -i -ldflags "$(LDFLAGS)" $(MAINGO)
	cd $(GOSRCDIR) && \
		CGO_ENABLED=0 \
		go test -v -i -c -ldflags "$(LDFLAGS)" $(GOPATHS)

run: build
	./$(MAIN)

clean:
	go clean
	rm -f $(MAIN)
	rm -rf $(TOP)/vendor/

lint:
	cd $(GOSRCDIR) && \
		gometalinter --vendor --errors --fast --deadline=60s -D gotype $(GOPATHS)

test: build lint
	./$(MAIN).test -test.timeout 60s

bench:
	cd $(GOSRCDIR) && \
		go test -bench=. -cpu 4 $(GOPATHS)

fmt:
	$(GOFMT) $(GOFILES)

dep:
	go get -u github.com/Masterminds/glide
	go get -u github.com/smartystreets/goconvey
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install --update
	cd $(GOSRCDIR) && glide install

init:
	cd $(GOSRCDIR) && \
		glide init

up:
	cd $(GOSRCDIR) && \
		glide update

gen:
	cd $(GOSRCDIR) && \
		go generate $(GOPATHS)
	$(GOFMT) $(GOFILES)

release:
	git push origin master master:production
