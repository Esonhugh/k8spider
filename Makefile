BUILD_PLUGIN_DIR = bin/plugins
BUILD_DIR = bin

default: build

# build cf plugin framework
build-cf:
	go	build -o $(BUILD_DIR)/ main/main.go

clean:
	rm -rf $(BUILD_DIR)