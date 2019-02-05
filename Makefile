PROJECT := gopal
PACKAGE := github.com/remerge/$(PROJECT)

include mkf/Makefile.common mkf/Makefile.app

ARTIFACTS_VERSION = $(shell git rev-parse --short HEAD)

.build/$(PROJECT)-$(ARTIFACTS_VERSION)-%_amd64.tar.gz: .build/$(PROJECT).%.amd64
	tar --show-transform --transform 's,$<,$(PROJECT),' -czvf $@ $<
	echo $@ $$(openssl sha256 < $@)
	tar -tvf $@

artifacts: \
	.build/$(PROJECT)-$(ARTIFACTS_VERSION)-linux_amd64.tar.gz \
	.build/$(PROJECT)-$(ARTIFACTS_VERSION)-darwin_amd64.tar.gz
