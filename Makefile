MAKEFLAGS += --silent

BUILD_DIR=build
RELEASE_DIR=$(BUILD_DIR)/release
TMP_DIR=$(BUILD_DIR)/tmp
VERSION := $(shell cat VERSION)
PLATFORM := $(shell uname -m)
GO=go

DEFAULT_SECRET := $(shell grep "var secretKey" main.go | cut -c24-85)
SECRET := $(shell cat secret)

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

package: clean prepare main.go~ cp compile control
	@echo Building package...
	fakeroot dpkg-deb -b -z9 $(TMP_DIR) $(RELEASE_DIR)

main.go~:
	cp main.go main.go~
	@sed -i "s/$(DEFAULT_SECRET)/$(SECRET)/" main.go

clean:
	rm -rf $(TMP_DIR) $(RELEASE_DIR)

prepare:
	@echo Prepare...
	mkdir -p $(TMP_DIR) $(RELEASE_DIR)

cp:
	cp -R deb/* $(TMP_DIR)

compile:
	go mod vendor
	GOOS=linux GOARCH=$(GOARCH) $(GO) build -o $(TMP_DIR)/usr/bin/delauncher main.go
	mv main.go~ main.go

control:
	$(eval size=$(shell du -sbk $(TMP_DIR)/ | grep -o '[0-9]*'))
	@sed -i "s/{{version}}/$(VERSION)/g;s/{{size}}/$(size)/g;s/{{architecture}}/$(ARCH)/g" "$(TMP_DIR)/DEBIAN/control"
