FILES = main.go
OUTPUT = bin/brightness-manager
USER_BIN = /usr/bin/brightness-manager

all:
	go build -o $(OUTPUT) $(FILES)

copy-to-path:
	sudo cp $(OUTPUT) $(USER_BIN)