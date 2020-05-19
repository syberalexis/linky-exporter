empty :=
PROJECT_NAME := linky-exporter
DIST_FOLDER := dist
TAG_NAME := $(shell git tag -l --contains HEAD | head -n1)
GOOS ?= $(shell go version | awk '{print $4}' | cut -d'/' -f1)
GOARCH ?= $(shell go version | awk '{print $4}' | cut -d'/' -f2)
VERSION := $(subst v,$(empty),$(TAG_NAME))

ifeq ($(GOARCH), arm)
    ARCH := armv$(GOARM)
else
    ARCH := $(GOARCH)
endif

default: build

dist:
	mkdir $(DIST_FOLDER)

build: dist
	go build -o $(DIST_FOLDER)/$(PROJECT_NAME)-$(VERSION)-$()-$(ARCH) -ldflags "-X main.version=$(VERSION)" cmd/$(PROJECT_NAME)/main.go

clean:
	rm -rf $(DIST_FOLDER)

version:
	echo $(VERSION)
