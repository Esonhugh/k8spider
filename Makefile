BUILD_PLUGIN_DIR = bin/plugins
BUILD_DIR = bin

MAIN_PROGRAM_NAME = k8spider

default: build

# build
build:
	go	build -o $(BUILD_DIR)/$(MAIN_PROGRAM_NAME) main.go

clean:
	rm -rf $(BUILD_DIR)