# Release Process

This document describes how to build and publish the Terraform Provider for StoneBranch to the JFrog Artifactory registry.

## Prerequisites

### Required Tools

- [Go](https://golang.org/doc/install) >= 1.24
- [GoReleaser](https://goreleaser.com/install/) >= 2.0
- [JFrog CLI](https://jfrog.com/getcli/) (optional, for manual uploads)

Install GoReleaser:

```bash
brew install goreleaser
```

### JFrog Artifactory Access

You need credentials for the JFrog Artifactory instance at `https://optionmetrics.jfrog.io`.

Configure JFrog CLI (for manual operations):

```bash
jf config add optionmetrics \
  --url=https://optionmetrics.jfrog.io \
  --user=your-username \
  --interactive
```

## Versioning

This project uses [Semantic Versioning](https://semver.org/) via git tags. The version is automatically derived from the latest git tag.

```bash
# Check current version
make version

# Create a new version tag
make tag V=0.3.0

# Push the tag to remote
git push origin v0.3.0
```

### Version Format

- Tags must follow the format `v{major}.{minor}.{patch}` (e.g., `v0.3.0`)
- The `v` prefix is stripped when creating release artifacts
- Without a tag, the version defaults to `0.0.0-dev`

## Building Releases

### Snapshot Build (Testing)

Build without publishing, no tag required:

```bash
make release-snapshot
```

This creates artifacts in `dist/` with a `-dev` suffix.

### Release Build (No Publish)

Build release artifacts from the current tag:

```bash
make release
```

### Build Artifacts

GoReleaser creates the following in `dist/`:

```
dist/
├── terraform-provider-stonebranch_0.3.0_darwin_amd64.zip
├── terraform-provider-stonebranch_0.3.0_darwin_arm64.zip
├── terraform-provider-stonebranch_0.3.0_linux_amd64.zip
├── terraform-provider-stonebranch_0.3.0_linux_arm64.zip
├── terraform-provider-stonebranch_0.3.0_windows_amd64.zip
└── terraform-provider-stonebranch_0.3.0_SHA256SUMS
```

## Publishing to Artifactory

### Authentication

Set environment variables before publishing:

```bash
# Option 1: Username and password/token
export JFROG_USER=your-username
export JFROG_PASSWORD=your-api-key

# Option 2: Access token
export ARTIFACTORY_TOKEN=your-access-token
```

You can add these to a `.env` file (gitignored) or your shell profile.

### Publish Command

Build and publish to Artifactory:

```bash
make publish
```

This will:
1. Build binaries for all platforms (darwin, linux, windows)
2. Create zip archives with proper naming
3. Generate SHA256 checksums
4. Upload all artifacts to Artifactory

### Artifactory Location

Artifacts are published to:

```
https://optionmetrics.jfrog.io/artifactory/terraform-providers/stonebranch/stonebranch/{version}/
```

### Verify Upload

```bash
# Using JFrog CLI
jf rt search "terraform-providers/stonebranch/stonebranch/*"

# Or check a specific version
jf rt search "terraform-providers/stonebranch/stonebranch/0.3.0/"
```

## Complete Release Workflow

```bash
# 1. Ensure all changes are committed
git status

# 2. Run tests
make testacc

# 3. Create and push version tag
make tag V=0.3.0
git push origin v0.3.0

# 4. Set credentials
export JFROG_USER=your-username
export JFROG_PASSWORD=your-api-key

# 5. Build and publish
make publish

# 6. Verify
jf rt search "terraform-providers/stonebranch/stonebranch/0.3.0/"
```

## Using the Published Provider

### Configure Terraform

Add to `~/.terraformrc` or project `.terraformrc`:

```hcl
provider_installation {
  direct {
    exclude = ["registry.terraform.io/stonebranch/stonebranch"]
  }

  filesystem_mirror {
    path    = "/path/to/local/providers"
    include = ["registry.terraform.io/stonebranch/stonebranch"]
  }
}
```

Or if Artifactory is configured as a Terraform registry:

```hcl
provider_installation {
  network_mirror {
    url = "https://optionmetrics.jfrog.io/artifactory/api/terraform/terraform-providers/"
  }
}
```

### Download Manually

```bash
# Download a specific version
curl -u $JFROG_USER:$JFROG_PASSWORD -O \
  "https://optionmetrics.jfrog.io/artifactory/terraform-providers/stonebranch/stonebranch/0.3.0/terraform-provider-stonebranch_0.3.0_darwin_arm64.zip"
```

## Configuration Files

| File | Purpose |
|------|---------|
| `.goreleaser.yaml` | GoReleaser build and publish configuration |
| `Makefile` | Build automation targets |
| `dist/` | Build output directory (gitignored) |

## Troubleshooting

### Authentication Errors (401)

- Verify `JFROG_USER` and `JFROG_PASSWORD` are set correctly
- Check that your API key/token hasn't expired
- Ensure you have write permissions to the `terraform-providers` repository

### Repository Not Found (404/405)

- Verify the repository `terraform-providers` exists in Artifactory
- Check you have the correct repository type (Generic or Terraform)

### Build Failures

```bash
# Check GoReleaser configuration
goreleaser check

# Run with debug output
goreleaser release --clean --snapshot --skip=publish --debug
```

### Version Issues

```bash
# List all tags
git tag -l

# Ensure tag is pushed
git push origin --tags

# Check current version
make version
```
