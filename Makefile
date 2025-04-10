PROJECT_ROOT := $(shell sh -c "git rev-parse --path-format=absolute --git-dir | xargs dirname")

APP_NAME := converter
BIN_DIR := $(PROJECT_ROOT)/bin
CMD_DIR := $(PROJECT_ROOT)/cmd/$(APP_NAME)

deps:
	# go get .
	go mod tidy

build: deps
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

clean:
	rm -rvf ${BIN_DIR}
