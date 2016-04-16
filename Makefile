LAST_TAG := $(shell git describe --abbrev=0 --tags)

USER := kagux
EXECUTABLE := remote_cli
REPOSITORY := go-remote-cli

# only include the amd64 binaries, otherwise the github release will become
# too big
UNIX_EXECUTABLES := \
	darwin/amd64/$(EXECUTABLE) \
	freebsd/amd64/$(EXECUTABLE) \
	linux/amd64/$(EXECUTABLE)

COMPRESSED_EXECUTABLES=$(UNIX_EXECUTABLES:%=%.tar.bz2)
COMPRESSED_EXECUTABLE_TARGETS=$(COMPRESSED_EXECUTABLES:%=bin/%)

UPLOAD_CMD = github-release upload -u $(USER) -r $(REPOSITORY) -t $(LAST_TAG) -n $(subst /,-,$(FILE)) -f bin/$(FILE)

all: $(EXECUTABLE)

# amd64
bin/freebsd/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=freebsd go build -o "$@"
bin/darwin/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=darwin go build -o "$@"
bin/linux/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=linux go build -o "$@"

# compressed artifacts, makes a huge difference (Go executable is ~9MB,
# after compressing ~2MB)
%.tar.bz2: %
	tar -jcvf "$<.tar.bz2" "$<"

github-release:
	go get github.com/aktau/github-release

# git tag -a v$(RELEASE) -m 'release $(RELEASE)'
release: clean github-release $(COMPRESSED_EXECUTABLE_TARGETS)
	git push && git push --tags
	github-release release -u $(USER) -r $(REPOSITORY) -t $(LAST_TAG) -n $(LAST_TAG) || true
	$(foreach FILE,$(COMPRESSED_EXECUTABLES),$(UPLOAD_CMD);)

# install and/or update all dependencies, run this from the project directory
# go get -u ./...
# go test -i ./
dep:
	go list -f '{{join .Deps "\n"}}' | xargs go list -e -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | xargs go get -u

$(EXECUTABLE): dep
	go build -o "$@"

install:
	go install

clean:
	rm go-app || true
	rm $(EXECUTABLE) || true
	rm -rf bin/

.PHONY: clean release dep install
