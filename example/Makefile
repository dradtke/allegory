
# available tags: console
TAGS = console

BUILD_DIR = build
TARGET = main.exe
PROF = cpu.prof

all:
	@mkdir -p $(BUILD_DIR)
	@echo -n "Building..."; go build -o $(BUILD_DIR)/$(TARGET) -tags '$(TAGS)'
	@echo "done."

run:
	@$(BUILD_DIR)/$(TARGET)

profile:
	@rm -f $(BUILD_DIR)/$(PROF)
	@$(BUILD_DIR)/$(TARGET) -cpuprofile=$(BUILD_DIR)/$(PROF)
	go tool pprof $(BUILD_DIR)/$(TARGET) $(BUILD_DIR)/$(PROF)

clean:
	@rm -rf $(BUILD_DIR)
	@rm -f allegro.log

.PHONY: all run profile clean
