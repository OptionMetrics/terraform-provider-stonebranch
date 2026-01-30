.PHONY: build install test testunit testacc testcov clean fmt lint dev

BINARY_NAME=terraform-provider-stonebranch
INSTALL_PATH=~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/0.1.0/darwin_arm64

build:
	go build -o $(BINARY_NAME)

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/

# Run all unit tests (no API calls required)
test:
	go test -v ./...

# Run only client unit tests
testunit:
	go test -v ./internal/client/...

# Run acceptance tests (requires TF_ACC=1 and real API credentials)
testacc:
	TF_ACC=1 go test -v ./internal/provider/... -timeout 120m

# Run acceptance tests for Unix task resource only
testacc-unix:
	TF_ACC=1 go test -v ./internal/provider/... -run='TestAccTaskUnix' -timeout 30m

# Run tests with coverage report
testcov:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Run terraform with dev overrides
dev:
	TF_CLI_CONFIG_FILE=./examples/dev.tfrc terraform -chdir=examples/provider $(ARGS)
