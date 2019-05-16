GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
BIN := japroxy

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

build:
	$(GOBUILD) -o ../../bin/$(OS)_$(ARCH)/$(BIN) -v

