PROJECT_NAME := "linky-exporter"
DIST_FOLDER := "dist"
TAG_NAME := $(shell git tag -l --contains HEAD)
ARCH := $(shell go version | awk '{print $4}' | cut -d'/' -f2)

default: build

dist:
	mkdir $(DIST_FOLDER)

build: dist
	go build -o $(DIST_FOLDER)/$(PROJECT_NAME)_$(TAG_NAME)_$(ARCH) cmd/$(PROJECT_NAME)/main.go

clean:
	rm -r $(DIST_FOLDER)
