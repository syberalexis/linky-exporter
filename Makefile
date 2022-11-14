empty :=
PROJECT_NAME := linky-exporter
DIST_FOLDER := dist
TAG_NAME := $(shell git tag -l --contains HEAD | head -n1)
GOOSARCH := $(shell go version | cut -d' ' -f4)
GOOS ?= $(shell echo ${GOOSARCH} | cut -d'/' -f1)
GOARCH ?= $(shell echo ${GOOSARCH} | cut -d'/' -f2)
VERSION ?= $(subst v,$(empty),$(TAG_NAME))

ifeq ($(GOARCH), arm)
    ARCH := armv$(GOARM)
else
    ARCH := $(GOARCH)
endif

default: build

dist:
	mkdir $(DIST_FOLDER)

build: dist
	go mod vendor
	go build -o $(DIST_FOLDER)/$(PROJECT_NAME)-$(VERSION)-$(GOOS)-$(ARCH) -ldflags "-X main.version=$(VERSION)" cmd/$(PROJECT_NAME)/main.go

clean:
	rm -rf $(DIST_FOLDER)

version:
	echo $(VERSION)

info:
	echo "PROJECT_NAME = $(PROJECT_NAME)"
	echo "DIST_FOLDER = $(DIST_FOLDER)"
	echo "TAG_NAME = $(TAG_NAME)"
	echo "GOOSARCH = $(GOOSARCH)"
	echo "GOOS = $(GOOS)"
	echo "GOARCH = $(GOARCH)"
	echo "VERSION = $(VERSION)"

docker:
	docker build -t syberalexis/linky-exporter:latest --build-arg VERSION=$(VERSION) .
ifneq ($(VERSION),)
	docker tag syberalexis/linky-exporter syberalexis/linky-exporter:$(VERSION)
endif
