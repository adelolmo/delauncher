MAKEFLAGS += --silent

BUILD_DIR=build
VERSION := $(shell cat VERSION)
GO=mewn

DEFAULT_SECRET := $(shell grep "var secretKey" main.go | cut -c24-85)
SECRET := $(shell cat secret)

ARCH := amd64
GOARCH := amd64

build: clean prepare main.go~
	@echo Building package...
	GOOS=linux GOARCH=$(GOARCH) $(GO) build -o $(BUILD_DIR)/tmp/usr/bin/delauncher main.go> /dev/null 2>&1
	mv main.go~ main.go
	$(eval size := $(shell du -cs $(BUILD_DIR)/tmp | sed '1!d' | grep -oe "^[0-9]*"))
	@sed -i "s/{{version}}/$(VERSION)/g;s/{{size}}/$(size)/g;s/{{architecture}}/$(ARCH)/g" $(BUILD_DIR)/tmp/DEBIAN/control
	fakeroot dpkg-deb -b -z9 $(BUILD_DIR)/tmp $(BUILD_DIR)/release

main.go~:
	cp main.go main.go~
	@sed -i "s/$(DEFAULT_SECRET)/$(SECRET)/" main.go

clean:
	rm -rf $(BUILD_DIR)/tmp $(BUILD_DIR)/release

prepare:
	@echo Prepare...
	mkdir -p $(BUILD_DIR)/tmp $(BUILD_DIR)/release
	cp -R deb/* $(BUILD_DIR)/tmp
	go get github.com/leaanthony/mewn/cmd/mewn > /dev/null 2>&1

