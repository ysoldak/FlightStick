SIZE   ?= full
TARGET ?= nano-33-ble-s140v6-uf2

ifneq ($(TARGET),nano-33-ble-s140v6-uf2)
FILE = fs_$(TARGET)_$(VERSION).uf2
else
FILE = fs_nano-33-ble_$(VERSION).uf2
endif

.PHONY: clean version build flash

# --- Common targets ---

VERSION := $(shell git describe --tags)

clean:
	@rm -rf build

version:
	@echo "package main" > src/version.go
	@echo "const version = \"$(VERSION)\"" >> src/version.go

build: version
	@mkdir -p build
	tinygo build -target=$(TARGET) -size=$(SIZE) -opt=z -print-allocs=main -o ./build/$(FILE) ./src

flash:
	tinygo flash -target=$(TARGET) -size=$(SIZE) -opt=z -print-allocs=main ./src
