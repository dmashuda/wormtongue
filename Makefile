BINARY_NAME=wormtongue
BUILD_DIR=bin
COVERAGE_THRESHOLD=80

.PHONY: build run test vet lint coverage install clean

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

run:
	WORMTONGUE_EXAMPLES=./examples go run . $(ARGS)

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

coverage:
	@go test -coverprofile=coverage.out ./...
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | tr -d '%'); \
	echo "Total coverage: $${COVERAGE}%"; \
	if [ "$$(echo "$${COVERAGE} < $(COVERAGE_THRESHOLD)" | bc)" -eq 1 ]; then \
		echo "FAIL: Coverage $${COVERAGE}% is below the $(COVERAGE_THRESHOLD)% threshold"; \
		rm -f coverage.out; \
		exit 1; \
	fi; \
	rm -f coverage.out

install:
	go install .

clean:
	rm -rf $(BUILD_DIR) coverage.out
