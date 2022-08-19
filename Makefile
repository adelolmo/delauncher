#MAKEFLAGS += --silent

BIN_DIR=/usr/bin
BIN=delauncher
BUILD_DIR=build
RELEASE_DIR=$(CURDIR)/..
VERSION := $(shell cat VERSION)
PLATFORM := $(shell uname -m)

ARCH :=
	ifeq ($(PLATFORM),x86_64)
		ARCH = amd64
	endif
	ifeq ($(PLATFORM),aarch64)
		ARCH = arm64
	endif
	ifeq ($(PLATFORM),armv7l)
		ARCH = armhf
	endif
GOARCH :=
	ifeq ($(ARCH),amd64)
		GOARCH = amd64
	endif
	ifeq ($(ARCH),i386)
		GOARCH = 386
	endif
	ifeq ($(ARCH),armhf)
		GOARCH = arm
	endif
	ifeq ($(ARCH),arm64)
		GOARCH = arm64
	endif

ifeq ($(GOARCH),)
	$(error Invalid ARCH: $(ARCH))
endif

all: build

$(BUILD_DIR)/DEBIAN:
	@echo Prapare package...
	cp -R deb/* $(BUILD_DIR)
	$(eval size=$(shell du -sbk $(BUILD_DIR) | grep -o '[0-9]*'))
	@sed -i "s/{{version}}/$(VERSION)/g;s/{{size}}/$(size)/g;s/{{architecture}}/$(ARCH)/g" "$(BUILD_DIR)/DEBIAN/control"

.PHONY: debian
debian: clean $(BUILD_DIR)/DEBIAN
	@echo Building package...
	mkdir $(BUILD_DIR)$(BIN_DIR)
	cp $(BIN) $(BUILD_DIR)$(BIN_DIR)
	chmod --quiet 0555 $(BUILD_DIR)/DEBIAN/p* || true
	fakeroot dpkg-deb -b -z9 $(BUILD_DIR) $(RELEASE_DIR)

.PHONY: clean
clean:
	@echo Clean...
	rm -rf $(BUILD_DIR)
	mkdir $(BUILD_DIR)

.PHONY: build
build:
	GOOS=linux GOARCH=$(GOARCH) go build -o $(BIN) .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor: tidy
	go mod vendor
