BUILD_PLUGIN_DIR = bin/plugins
BUILD_DIR = bin

MAIN_PROGRAM_NAME = k8spider

default: build build-static check-size

# build
build:
	go	build -o $(BUILD_DIR)/$(MAIN_PROGRAM_NAME) main.go

test: build
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) all --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) axfr --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) dns --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) dnssd --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) dnssd ptr --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) dnssd srv --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) metrics --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) neighbor --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) neighbor pod --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) neighbor svc --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) whereisdns --help
	$(BUILD_DIR)/$(MAIN_PROGRAM_NAME) wild --help

build-static:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(MAIN_PROGRAM_NAME)-linux-static main.go
	upx --lzma --brute  $(BUILD_DIR)/$(MAIN_PROGRAM_NAME)-linux-static

check-size:
	ls -alh $(BUILD_DIR)/$(MAIN_PROGRAM_NAME)*

clean:
	rm -rf $(BUILD_DIR)