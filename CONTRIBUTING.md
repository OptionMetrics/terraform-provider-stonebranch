# Contributing to Terraform Provider for StoneBranch

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing.

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Create a feature branch from `main`
4. Make your changes
5. Submit a pull request

## Development Setup

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.24
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- A StoneBranch Universal Controller instance (for acceptance tests)

### Building

```bash
make build
```

### Running Tests

```bash
# Unit tests (no credentials needed)
make test

# Acceptance tests (requires StoneBranch instance)
export STONEBRANCH_API_TOKEN="your-token"
export STONEBRANCH_BASE_URL="https://your-instance.stonebranch.cloud"
make testacc
```

### Generating Documentation

After modifying resource schemas, regenerate documentation:

```bash
make docs
```

## Pull Request Guidelines

- Keep PRs focused on a single change
- Include tests for new resources or bug fixes
- Run `make fmt` before submitting
- Run `make test` to ensure unit tests pass
- Update example files in `examples/` if adding new resources
- Regenerate docs with `make docs` if schemas changed

## Adding a New Resource

1. Create the resource file in `internal/provider/resources/`
2. Register it in `internal/provider/provider.go`
3. Add acceptance tests in `internal/provider/resources/*_test.go`
4. Add example configuration in `examples/resources/stonebranch_<resource>/resource.tf`
5. Run `make docs` to generate documentation
6. Update CHANGELOG.md

### Naming Conventions

- Avoid Go platform-specific file suffixes (`_windows.go`, `_linux.go`, `_darwin.go`) as Go treats these as build constraints
- Use `taskwindows.go` instead of `task_windows.go`

## Reporting Issues

- Use GitHub Issues for bug reports and feature requests
- Include Terraform version, provider version, and relevant configuration
- For bugs, include the full error output and debug logs if possible (`TF_LOG=DEBUG`)

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
