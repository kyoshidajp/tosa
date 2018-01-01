PROJECT_NAME=tosa
INTERNAL_BIN_DIR=_internal_bin
GOVERSION=$(shell go version)
GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
VERSION=$(patsubst "%",%,$(lastword $(shell grep 'const version' $(PROJECT_NAME).go)))
RELEASE_DIR=releases
ARTIFACTS_DIR=$(RELEASE_DIR)/artifacts/$(VERSION)
SRC_FILES = $(wildcard *.go cmd/$(PROJECT_NAME)/*.go)
HAVE_GLIDE:=$(shell which glide >/dev/null 2>&1)
GITHUB_USERNAME=kyoshidajp

build: $(RELEASE_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH)/$(PROJECT_NAME)$(SUFFIX)

installdeps: glide $(SRC_FILES)
	@echo "Installing dependencies..."
	@PATH=$(INTERNAL_BIN_DIR)/$(GOOS)/$(GOARCH):$(PATH) glide install

build-darwin-386:
	@$(MAKE) build GOOS=darwin GOARCH=386

$(RELEASE_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH)/$(PROJECT_NAME)$(SUFFIX):
	go build -o $(RELEASE_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH)/$(PROJECT_NAME)$(SUFFIX) cmd/$(PROJECT_NAME)/$(PROJECT_NAME).go

test:
	$(GOTEST) -v ./...

clean:
	-rm -rf $(RELEASE_DIR)/*/*
	-rm -rf $(ARTIFACTS_DIR)/*
