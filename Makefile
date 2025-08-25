APPNAME=dev.flowingspdg.characterworks.sdPlugin

GOFLAGS=

MAKEFILE_DIR:=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BUILDDIR = $(MAKEFILE_DIR)$(APPNAME)
SRCDIR = $(MAKEFILE_DIR)Source
PIDIR = $(MAKEFILE_DIR)Source/pi
IMAGEDIR = $(MAKEFILE_DIR)Source/images
RELEASEDIR = Release

BINARY = cw_client

RM = rm -rf
ifeq ($(OS),Windows_NT)
    RM = Remove-Item -Recurse -Force
endif

MKDIR = mkdir -p
ifeq ($(OS),Windows_NT)
    MKDIR = New-Item -Force -ItemType Directory
endif

CP = cp -R
ifeq ($(OS),Windows_NT)
    CP = powershell -Command Copy-Item -Recurse -Force
endif

DISTRIBUTION_TOOL = ./DistributionTool.exe
ifeq ($(shell uname),Darwin)
    DISTRIBUTION_TOOL = ./DistributionTool
endif

.DEFAULT_GOAL := build

prepare:
	@$(MKDIR) $(BUILDDIR)
	@$(RM) $(BUILDDIR)/*

build: prepare
	cd $(SRCDIR) && GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o $(BUILDDIR)/$(BINARY) .
	cd $(SRCDIR) && GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o $(BUILDDIR)/$(BINARY).exe .
	$(CP) $(PIDIR) $(BUILDDIR)/inspector
	$(CP) $(SRCDIR)/manifest.json $(BUILDDIR)
	$(CP) $(IMAGEDIR) $(BUILDDIR)

package: build
	@$(RM) ./$(RELEASEDIR)/*
	@$(MKDIR) $(RELEASEDIR)
	$(DISTRIBUTION_TOOL) -b -i $(APPNAME) -o $(RELEASEDIR)

clean:
	@$(RM) $(BUILDDIR)
	@$(RM) $(RELEASEDIR)/* 