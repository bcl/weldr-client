DESTDIR ?= /
VERSION ?= 36.1
TAG = v$(VERSION)
PREVTAG := $(shell git tag --sort=-creatordate | head -n 2 | tail -n 1)
COMMITS := $(shell git log --pretty=oneline --no-merges ${PREVTAG}..HEAD | wc -l)
GPGKEY ?= $(shell git config user.signingkey)
GITEMAIL := $(shell git config user.email)
GITNAME := $(shell git config user.name)
GOBUILDFLAGS ?= -ldflags="-X github.com/osbuild/weldr-client/v2/cmd/composer-cli/root.Version=${VERSION}"

build: composer-cli
composer-cli:
	go build ${GOBUILDFLAGS} ./cmd/composer-cli

check:
	go vet ./... && golangci-lint --build-tags=integration run ./...

test:
	go test ${GOBUILDFLAGS} -v -covermode=atomic -coverprofile=coverage.txt -coverpkg=./... ./...

integration: composer-cli-tests
composer-cli-tests:
	go test -c -tags=integration ${GOBUILDFLAGS} -o composer-cli-tests ./weldr/

install: composer-cli
	install -m 0755 -vd ${DESTDIR}/usr/bin/
	install -m 0755 -vp composer-cli ${DESTDIR}/usr/bin/
	install -m 0755 -vd ${DESTDIR}/etc/bash_completion.d/
	install -m 0644 -vp etc/bash_completion.d/composer-cli ${DESTDIR}/etc/bash_completion.d/
	install -m 0755 -vd ${DESTDIR}/usr/share/man/man1/
	./composer-cli doc ${DESTDIR}/usr/share/man/man1/

install-tests: composer-cli-tests
	install -m 0755 -vd ${DESTDIR}/usr/libexec/tests/composer-cli/
	install -m 0755 -vp composer-cli-tests ${DESTDIR}/usr/libexec/tests/composer-cli/

weldr-client.spec: weldr-client.spec.in
	sed -e "s/%%VERSION%%/$(VERSION)/g" -e "s/%%GPGKEY%%/$(GPGKEY)/g" < $< > $@
	$(MAKE) -s changelog >> clog
	@echo "Use the contents of ./clog when committing the new spec file"

tag:
	@if [ -z "$(GPGKEY)" ]; then echo "ERROR: The git config user.signingkey must be set" ; exit 1; fi
	git tag -u $(GPGKEY) -m "Tag as $(TAG)" -f $(TAG)
	@echo "Tagged as $(TAG)"

# Order matters, so run make for each step instead of declaring them as dependencies
release:
	@if [ -z "$(GPGKEY)" ]; then echo "ERROR: The git config user.signingkey must be set" ; exit 1; fi
	$(MAKE) test && $(MAKE) bumpver && $(MAKE) tag && $(MAKE) archive && $(MAKE) sign

sign:
	@if [ -z "$(GPGKEY)" ]; then echo "ERROR: The git config user.signingkey must be set" ; exit 1; fi
	gpg --armor --detach-sign -u $(GPGKEY) weldr-client-$(VERSION).tar.gz
	gpg2 --export --export-options export-minimal $(GPGKEY) > gpg-$(GPGKEY).key

changelog:
	@echo "* $(shell date '+%a %b %d %Y') ${GITNAME} <${GITEMAIL}> - ${VERSION}-1"
	@git log --no-merges --pretty="format:- %s (%ae)" ${PREVTAG}..HEAD |sed -e 's/@.*)/)/'

bumpver:
	@NEWSUBVER=$$((`echo $(VERSION) |cut -d . -f 2` + 1)) ; \
	NEWVERSION=`echo $(VERSION).$$NEWSUBVER |cut -d . -f 1,3` ; \
	sed -i "s/VERSION ?= $(VERSION)/VERSION ?= $$NEWVERSION/" Makefile ; \
	git add Makefile; \
	git commit -m "New release: $$NEWVERSION"

archive:
	git archive --prefix=weldr-client-$(VERSION)/ --format=tar.gz HEAD > weldr-client-$(VERSION).tar.gz

RPM_SPECFILE=rpmbuild/SPECS/weldr-client.spec
RPM_TARBALL=rpmbuild/SOURCES/weldr-client-$(VERSION).tar.gz
RPM_TARBALL_SIG=rpmbuild/SOURCES/weldr-client-$(VERSION).tar.gz.asc

$(RPM_SPECFILE): weldr-client.spec
	mkdir -p $(CURDIR)/rpmbuild/SPECS
	cp weldr-client.spec $(CURDIR)/rpmbuild/SPECS

$(RPM_TARBALL): archive sign
	mkdir -p $(CURDIR)/rpmbuild/SOURCES
	cp weldr-client-$(VERSION).tar.gz* gpg-$(GPGKEY).key rpmbuild/SOURCES/

builddep: $(RPM_SPECFILE)
	dnf builddep -y -D 'with 1' $(RPM_SPECFILE)

srpm: $(RPM_SPECFILE) $(RPM_TARBALL)
	rpmbuild -bs \
		--define "_topdir $(CURDIR)/rpmbuild" \
		--define "commit $(VERSION)" \
		--with tests \
		$(RPM_SPECFILE)

rpm: $(RPM_SPECFILE) $(RPM_TARBALL)
	rpmbuild -bb \
		--define "_topdir $(CURDIR)/rpmbuild" \
		--define "commit $(VERSION)" \
		--with tests \
		$(RPM_SPECFILE)

scratch-srpm: $(RPM_SPECFILE) archive
	mkdir -p $(CURDIR)/rpmbuild/SOURCES
	cp weldr-client-$(VERSION).tar.gz* rpmbuild/SOURCES/
	rpmbuild -bs \
		--define "_topdir $(CURDIR)/rpmbuild" \
		--define "commit $(VERSION)" \
		--with tests \
		--without signed \
		$(RPM_SPECFILE)

scratch-rpm: $(RPM_SPECFILE) archive
	mkdir -p $(CURDIR)/rpmbuild/SOURCES
	cp weldr-client-$(VERSION).tar.gz* rpmbuild/SOURCES/
	rpmbuild -bb \
		--define "_topdir $(CURDIR)/rpmbuild" \
		--define "commit $(VERSION)" \
		--with tests \
		--without signed \
		$(RPM_SPECFILE)

build-in-podman: archive
	podman build --security-opt seccomp=unconfined \
		--build-arg VERSION=$(VERSION) \
		-t weldr-client:$(VERSION) \
		-f Containerfile.test .

CONTAINER_IMAGE_CLI ?= osbuild-cli_dev
CONTAINERFILE_CLI ?= Containerfile.dev

CONTAINER_EXECUTABLE ?= podman

SRC_DEPS_DIRS := cmd weldr examples
CLI_SRC_DEPS := $(shell find $(SRC_DEPS_DIRS) -name *.go -or -name *.toml)

container_cli_built.info: $(CLI_SRC_DEPS) $(CONTAINERFILE_CLI)
	$(CONTAINER_EXECUTABLE) build -t $(CONTAINER_IMAGE_CLI) -f $(CONTAINERFILE_CLI) --build-arg GOMODARGS="$(GOMODARGS)" --build-arg GCFLAGS="$(GCFLAGS)" .
	echo "CLI last built on" > $@
	date >> $@

# build a container with the cli from full source
.PHONY: container.dev
container.dev: container_cli_built.info

update-mods:
	go get -u ./...
	go mod vendor
	$(MAKE) test

clean:
	rm -f weldr-client.spec
	rm -f composer-cli
	rm -f container_cli_built.info
	rm -rf $(CURDIR)/rpmbuild

.PHONY: build check test integration install srpm rpm weldr-client.spec update-mods build-in-podman
