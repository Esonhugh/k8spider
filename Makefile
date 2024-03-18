BUILD_PLUGIN_DIR = bin/plugins
BUILD_DIR = bin

MAIN_PROGRAM_NAME = k8spider

default: build build-static

# build
build:
	go	build -o $(BUILD_DIR)/$(MAIN_PROGRAM_NAME) main.go

build-static:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(MAIN_PROGRAM_NAME)-linux-static main.go

clean:
	rm -rf $(BUILD_DIR)