BINARY_NAME=wormtongue
BUILD_DIR=bin

.PHONY: build run test vet install clean

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

run:
	WORMTONGUE_EXAMPLES=./examples go run . $(ARGS)

test:
	go test ./...

vet:
	go vet ./...

install:
	go install .

clean:
	rm -rf $(BUILD_DIR)
