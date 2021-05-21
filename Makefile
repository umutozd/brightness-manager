FILES = main.go
OUTPUT = bin/brightness-manager

all:
	go build -o $(OUTPUT) $(FILES)