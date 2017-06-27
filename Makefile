PROJECT := gopal
PACKAGE := github.com/remerge/$(PROJECT)

GOMETALINTER_OPTS = --enable-all --tests --fast --errors

include Makefile.common
