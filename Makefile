MAKEFLAGS += --silent

BIN_DIR=/usr/bin
BIN=delauncher
BUILD_DIR=build-debian
RELEASE_DIR := $(realpath $(CURDIR)/..)

ASSETS_DIR = usr/share/icons/hicolor/16x16/status
APP_ICON_DIR = usr/share/icons/hicolor/128x128/apps

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

.PHONY: all
all: build

.PHONY: debian
debian: clean $(BUILD_DIR)/DEBIAN
	@echo Building package...
	cp $(BIN) $(BUILD_DIR)$(BIN_DIR)
	chmod --quiet 0555 $(BUILD_DIR)/DEBIAN/p* || true
	fakeroot dpkg-deb -b -z9 $(BUILD_DIR) $(RELEASE_DIR)

.PHONY: clean
clean:
	@echo Clean...
	rm -rf $(BUILD_DIR)

$(BUILD_DIR)/DEBIAN: $(BUILD_DIR)
	@echo Prapare package...
	cp -R deb/DEBIAN $(BUILD_DIR)
	$(MAKE) install DESTDIR=$(BUILD_DIR)
	$(eval SIZE := $(shell du -sbk $(BUILD_DIR) | grep -o '[0-9]*'))
	@sed -i "s/{{version}}/$(VERSION)/g;s/{{size}}/$(SIZE)/g;s/{{architecture}}/$(ARCH)/g" "$(BUILD_DIR)/DEBIAN/control"

$(BUILD_DIR):
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

.PHONY: install
install:
	install -Dm755 $(BIN) $(DESTDIR)$(BIN_DIR)/$(BIN)
	install -Dm644 deb/usr/share/applications/delauncher.desktop $(DESTDIR)/usr/share/applications/delauncher.desktop
	install -Dm644 deb/$(ASSETS_DIR)/delauncher-error.png $(DESTDIR)/$(ASSETS_DIR)/delauncher-error.png
	install -Dm644 deb/$(ASSETS_DIR)/delauncher-success.png $(DESTDIR)/$(ASSETS_DIR)/delauncher-success.png
	install -Dm644 deb/$(APP_ICON_DIR)/delauncher.png $(DESTDIR)/$(APP_ICON_DIR)/delauncher.png

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(BIN_DIR)/$(BIN)
	rm -f $(DESTDIR)/usr/share/applications/delauncher.desktop
	rm -f $(DESTDIR)/$(ASSETS_DIR)/delauncher-error.png
	rm -f $(DESTDIR)/$(ASSETS_DIR)/delauncher-success.png
	rm -f $(DESTDIR)/$(APP_ICON_DIR)/delauncher.png
