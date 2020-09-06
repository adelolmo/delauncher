BUILD_DIR=build
VERSION := $(shell cat VERSION)

build: clean prepare
	@echo Building package...

clean:
	rm -rf $(BUILD_DIR)/tmp $(BUILD_DIR)/package $(BUILD_DIR)/release

prepare:
	mkdir -p $(BUILD_DIR)/tmp $(BUILD_DIR)/release $(BUILD_DIR)/package/opt/mocorunner
	cp -R package/* $(BUILD_DIR)/package
