MAKEFLAGS += --silent

BUILD_DIR=build
RELEASE_DIR=$(BUILD_DIR)/release
TMP_DIR=$(BUILD_DIR)/tmp
VERSION := $(shell cat VERSION)
GO=mewn

DEFAULT_SECRET := $(shell grep "var secretKey" main.go | cut -c24-85)
SECRET := $(shell cat secret)

ARCH := amd64
GOARCH := amd64

build: clean prepare main.go~ cp gobuild control
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
	go get github.com/leaanthony/mewn/cmd/mewn > /dev/null 2>&1

cp:
	cp -R deb/* $(TMP_DIR)

gobuild:
	GOOS=linux GOARCH=$(GOARCH) $(GO) build -o $(TMP_DIR)/usr/bin/delauncher main.go> /dev/null 2>&1
	mv main.go~ main.go

control:
	$(eval size=$(shell du -sbk $(TMP_DIR)/ | grep -o '[0-9]*'))
	@sed -i "s/{{version}}/$(VERSION)/g;s/{{size}}/$(size)/g;s/{{architecture}}/$(ARCH)/g" "$(TMP_DIR)/DEBIAN/control"
