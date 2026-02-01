# Release Process

This document describes how to build and publish the Terraform Provider for StoneBranch to GitHub Releases.

## Prerequisites

### Required Tools

- [Go](https://golang.org/doc/install) >= 1.24
- [GoReleaser](https://goreleaser.com/install/) >= 2.0

Install on macOS:

```bash
brew install goreleaser
```

### GitHub Token

Create a GitHub personal access token with `repo` scope:

1. Go to https://github.com/settings/tokens
2. Generate new token (classic)
3. Select `repo` scope
4. Copy the token

Set it as an environment variable:

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
```

Or add to your shell profile (`~/.zshrc` or `~/.bashrc`).

## Versioning

This project uses [Semantic Versioning](https://semver.org/) via git tags.

### Version Format

- Tags must follow the format `v{major}.{minor}.{patch}` (e.g., `v0.3.0`)
- The `v` prefix is stripped when creating release artifacts
- Without a tag, the version defaults to `0.0.0-dev`

### Check Current Version

```bash
make version
```

## Complete Release Workflow

### Step 1: Regenerate Documentation

Ensure documentation is up-to-date with any schema changes:

```bash
make docs
git add docs/
git diff --cached --quiet || git commit -m "Update generated documentation"
```

### Step 2: Ensure Clean Git State

**Important:** GoReleaser requires a clean git state. All changes must be committed before creating a tag.

```bash
git status
git add .
git commit -m "Your commit message"
```

### Step 3: Create Version Tag

Only create the tag after ALL changes are committed:

```bash
make tag V=0.4.0
```

### Step 4: Push to GitHub

Push commits and tag:

```bash
git push origin main
git push origin v0.4.0
```

### Step 5: Publish Release

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
make publish
```

This will:
1. Build binaries for all platforms
2. Create zip archives
3. Generate SHA256 checksums
4. Create a GitHub Release with all artifacts attached

### Step 6: Verify

Check the release at:
https://github.com/OptionMetrics/terraform-provider-stonebranch/releases

## Quick Reference

```bash
# Full release workflow
make docs                           # Regenerate documentation
git status                          # Ensure clean state
git add . && git commit -m "msg"    # Commit changes (if needed)
make tag V=0.4.0                    # Create version tag
git push origin main                # Push commits
git push origin v0.4.0              # Push tag
export GITHUB_TOKEN=ghp_xxx         # Set token (if not in profile)
make publish                        # Build and publish
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make version` | Show current version from git tag |
| `make docs` | Generate provider documentation |
| `make release-snapshot` | Build without tag (for testing) |
| `make release` | Build release artifacts (no publish) |
| `make publish` | Build and publish to GitHub Releases |
| `make tag V=x.y.z` | Create new version tag |

## Build Artifacts

GoReleaser creates the following:

```
terraform-provider-stonebranch_0.4.0_darwin_amd64.zip
terraform-provider-stonebranch_0.4.0_darwin_arm64.zip
terraform-provider-stonebranch_0.4.0_linux_amd64.zip
terraform-provider-stonebranch_0.4.0_linux_arm64.zip
terraform-provider-stonebranch_0.4.0_windows_amd64.zip
terraform-provider-stonebranch_0.4.0_SHA256SUMS
```

## Using the Published Provider

### Step 1: Download and Install

#### macOS (Apple Silicon)

```bash
VERSION=0.4.0
curl -LO "https://github.com/OptionMetrics/terraform-provider-stonebranch/releases/download/v${VERSION}/terraform-provider-stonebranch_${VERSION}_darwin_arm64.zip"
unzip terraform-provider-stonebranch_${VERSION}_darwin_arm64.zip
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/${VERSION}/darwin_arm64
mv terraform-provider-stonebranch_v${VERSION} ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/${VERSION}/darwin_arm64/
rm terraform-provider-stonebranch_${VERSION}_darwin_arm64.zip
```

#### macOS (Intel)

```bash
VERSION=0.4.0
curl -LO "https://github.com/OptionMetrics/terraform-provider-stonebranch/releases/download/v${VERSION}/terraform-provider-stonebranch_${VERSION}_darwin_amd64.zip"
unzip terraform-provider-stonebranch_${VERSION}_darwin_amd64.zip
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/${VERSION}/darwin_amd64
mv terraform-provider-stonebranch_v${VERSION} ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/${VERSION}/darwin_amd64/
rm terraform-provider-stonebranch_${VERSION}_darwin_amd64.zip
```

#### Linux (x86_64)

```bash
VERSION=0.4.0
curl -LO "https://github.com/OptionMetrics/terraform-provider-stonebranch/releases/download/v${VERSION}/terraform-provider-stonebranch_${VERSION}_linux_amd64.zip"
unzip terraform-provider-stonebranch_${VERSION}_linux_amd64.zip
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/${VERSION}/linux_amd64
mv terraform-provider-stonebranch_v${VERSION} ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/${VERSION}/linux_amd64/
rm terraform-provider-stonebranch_${VERSION}_linux_amd64.zip
```

### Step 2: Configure Terraform Project

Create a `main.tf` file:

```hcl
terraform {
  required_providers {
    stonebranch = {
      source  = "registry.terraform.io/stonebranch/stonebranch"
      version = "0.4.0"
    }
  }
}

provider "stonebranch" {
  # API token - can also use STONEBRANCH_API_TOKEN env var
  api_token = var.stonebranch_api_token

  # Base URL - can also use STONEBRANCH_BASE_URL env var
  base_url = "https://your-instance.stonebranch.cloud"
}

variable "stonebranch_api_token" {
  type      = string
  sensitive = true
}
```

### Step 3: Initialize and Use

```bash
# Set your API token
export STONEBRANCH_API_TOKEN="your-bearer-token"

# Initialize Terraform (downloads provider from local plugin directory)
terraform init

# Plan changes
terraform plan

# Apply changes
terraform apply
```

### Example: Create a Unix Task

```hcl
resource "stonebranch_task_unix" "hello" {
  name    = "terraform-hello-world"
  summary = "A simple task managed by Terraform"
  command = "echo 'Hello from Terraform!'"
  agent   = "your-linux-agent"
}
```

### All Releases

Browse all versions at:
https://github.com/OptionMetrics/terraform-provider-stonebranch/releases

## Troubleshooting

### "git is in a dirty state"

GoReleaser requires all changes to be committed before building a release.

```bash
git status
git add .
git commit -m "message"
git tag -d v0.4.0
make tag V=0.4.0
```

### "tag was not made against commit"

The tag points to a different commit than HEAD:

```bash
git tag -d v0.4.0
make tag V=0.4.0
```

### "missing GITHUB_TOKEN"

Set the environment variable:

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
```

### Tag Already Exists on Remote

```bash
git push origin v0.4.0 --force
```

**Warning:** Only force-push tags if the release hasn't been used by others.

## Configuration Files

| File | Purpose |
|------|---------|
| `.goreleaser.yaml` | GoReleaser build configuration |
| `Makefile` | Build automation targets |
| `dist/` | Build output directory (gitignored) |
