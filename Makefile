BINARY_NAME=gostutter
GOPATH ?= $(shell go env GOPATH)
BIN_DIR := $(GOPATH)/bin

all: build

build:
	go mod tidy
	go build -o ${BINARY_NAME} ./cmd/gostutter
run:
	build
	./${BINARY_NAME}
clean:
	rm ${BINARY_NAME}

vet:
	go vet ...

lint:
	${BIN_DIR}/golangci-lint run

install:
	go install -o ./bin/${BINARY_NAME} ./cmd/gostutter