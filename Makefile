SIZE   ?= full
TARGET ?= nano-33-ble-s140v6-uf2

VERSION := $(shell git describe --tags)
LD_FLAGS := -ldflags="-X 'main.Version=$(VERSION)'"

ifneq ($(TARGET),nano-33-ble-s140v6-uf2)
FILE = fs_$(TARGET)_$(VERSION).uf2
else
FILE = fs_nano-33-ble_$(VERSION).uf2
endif

.PHONY: clean build flash

# --- Common targets ---

clean:
	@rm -rf build

build:
	@mkdir -p build
	tinygo build $(LD_FLAGS) -target=$(TARGET) -size=$(SIZE) -opt=z -print-allocs=main -o ./build/$(FILE) ./src

flash:
	tinygo flash $(LD_FLAGS) -target=$(TARGET) -size=$(SIZE) -opt=z -print-allocs=main ./src

# --- Arduino Nano 33 BLE bootloader targets ---

UF2_BOOTLOADER_HEX=./build/arduino_nano_33_ble_bootloader-0.7.0_s140_6.1.1.hex

$(UF2_BOOTLOADER_HEX):
	@curl -L -o $(UF2_BOOTLOADER_HEX) https://github.com/adafruit/Adafruit_nRF52_Bootloader/releases/download/0.7.0/arduino_nano_33_ble_bootloader-0.7.0_s140_6.1.1.hex

flash-uf2-bootloader: $(UF2_BOOTLOADER_HEX)
	nrfjprog -f nrf52 --eraseall
	nrfjprog -f nrf52 --program $(UF2_BOOTLOADER_HEX)

flash-uf2-bootloader-dap: $(UF2_BOOTLOADER_HEX)
	openocd -f interface/cmsis-dap.cfg -f target/nrf52.cfg -c "transport select swd" -c "program $(UF2_BOOTLOADER_HEX) verify reset exit"
