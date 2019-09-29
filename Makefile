NAME=jira-ui
GO?=go

DIST=$(CWD)$(SEP)dist

GOBIN ?= $(CWD)

CURVER ?= $(patsubst v%,%,$(shell [ -d .git ] && git describe --abbrev=0 --tags || grep ^\#\# CHANGELOG.md | awk '{print $$2; exit}'))
LDFLAGS:= -w

build:
	$(GO) build -gcflags="-e" -v -ldflags "$(LDFLAGS) -s" -o '$(BIN)' jira-ui/main.go

vet:
	@$(GO) vet .
	@$(GO) vet ./jira-ui

lint:
	@$(GO) get github.com/golang/lint/golint
	@golint .
	@golint ./jira-ui

all:
	GO111MODULE=off $(GO) get -u github.com/mitchellh/gox
	rm -rf dist
	mkdir -p dist
	gox -ldflags="-w -s" -output="dist/github.com/go-jira/jira-{{.OS}}-{{.Arch}}" -osarch="darwin/amd64 linux/386 linux/amd64 windows/386 windows/amd64" ./...
	_t/test_binaries.sh

install:
	${MAKE} GOBIN=$$HOME/bin build

NEWVER ?= $(shell echo $(CURVER) | awk -F. '{print $$1"."$$2"."$$3+1}')
TODAY  := $(shell date +%Y-%m-%d)

release:
	make update-usage
	git diff --exit-code --quiet README.md || git commit -m "Updated Usage" README.md
	git commit -m "Updated Changelog" CHANGELOG.md
	git commit -m "version bump" jira.go
	git tag v$(NEWVER)
	git push --tags

version:
	@echo $(CURVER)

clean:
	rm -rf ./$(NAME)


