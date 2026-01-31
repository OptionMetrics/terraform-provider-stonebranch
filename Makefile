.PHONY: build install test testunit testacc testcov clean fmt lint dev release release-snapshot publish tag version

BINARY_NAME=terraform-provider-stonebranch
INSTALL_PATH=~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/0.1.0/darwin_arm64

# Version from git tag (strips leading 'v'), fallback to 0.0.0-dev
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0-dev")

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
	rm -rf dist/
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Run terraform with dev overrides
dev:
	TF_CLI_CONFIG_FILE=./examples/dev.tfrc terraform -chdir=examples/provider $(ARGS)

# Build release binaries using goreleaser (dry-run, no publish)
release:
	goreleaser release --clean --skip=publish

# Build snapshot release (no tag required, for testing)
release-snapshot:
	goreleaser release --clean --snapshot --skip=publish

# Build and publish to Artifactory
# Uses jf CLI (authenticate first with 'jf login')
publish: release
	@echo "Publishing to Artifactory..."
	@jf rt upload "dist/*.zip" "terraform-providers/stonebranch/stonebranch/$(VERSION)/" --flat
	@jf rt upload "dist/*SHA256SUMS" "terraform-providers/stonebranch/stonebranch/$(VERSION)/" --flat
	@echo "Published version $(VERSION) to Artifactory"

# Create a new version tag (usage: make tag V=0.3.0)
tag:
	@if [ -z "$(V)" ]; then echo "Usage: make tag V=0.3.0"; exit 1; fi
	@echo "Creating tag v$(V)..."
	git tag -a "v$(V)" -m "Release v$(V)"
	@echo "Tag v$(V) created. Push with: git push origin v$(V)"

# Show current version
version:
	@echo $(VERSION)
