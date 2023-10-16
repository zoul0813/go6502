
APP_NAME=go6502
CC=go
BIN=bin/$(APP_NAME)


build:
	$(CC) build -o $(BIN)

run: build
	$(BIN)

log: build
	$(BIN) --debug > run.log