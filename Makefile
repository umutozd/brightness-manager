FILES = main.go
OUTPUT = bin/brightness-manager
USER_BIN = /usr/bin/brightness-manager

all: $(OUTPUT)

$(OUTPUT): $(FILES)
	mkdir -p $(dir $@)
	go build -o $@ $^

copy-to-path:
	sudo cp $(OUTPUT) $(USER_BIN)
