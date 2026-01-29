.PHONY: build install test clean

BINARY_NAME=terraform-provider-stonebranch
INSTALL_PATH=~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/0.1.0/darwin_arm64

build:
	go build -o $(BINARY_NAME)

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY_NAME) $(INSTALL_PATH)/

test:
	go test -v ./...

clean:
	rm -f $(BINARY_NAME)
	go clean

fmt:
	go fmt ./...

lint:
	golangci-lint run

# Run terraform with dev overrides
dev:
	TF_CLI_CONFIG_FILE=./examples/dev.tfrc terraform -chdir=examples/provider $(ARGS)
