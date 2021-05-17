BIN_DIR	:= ./bin
TARGET	:= scp-action

$(BIN_DIR)/$(TARGET): main.go
	@mkdir -p $(@D)
	go build -o $@ $^

.PHONY: all clean

all: $(BIN_DIR)/$(TARGET)

clean:
	-rm -rf $(BIN_DIR)
