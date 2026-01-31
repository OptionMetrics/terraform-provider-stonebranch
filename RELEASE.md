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

### Step 1: Ensure Clean Git State

**Important:** GoReleaser requires a clean git state. All changes must be committed before creating a tag.

```bash
git status
git add .
git commit -m "Your commit message"
```

### Step 2: Create Version Tag

Only create the tag after ALL changes are committed:

```bash
make tag V=0.4.0
```

### Step 3: Push to GitHub

Push commits and tag:

```bash
git push origin main
git push origin v0.4.0
```

### Step 4: Publish Release

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
make publish
```

This will:
1. Build binaries for all platforms
2. Create zip archives
3. Generate SHA256 checksums
4. Create a GitHub Release with all artifacts attached

### Step 5: Verify

Check the release at:
https://github.com/OptionMetrics/terraform-provider-stonebranch/releases

## Quick Reference

```bash
# Full release workflow
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

### Download from GitHub Releases

Download the appropriate zip for your platform from:
https://github.com/OptionMetrics/terraform-provider-stonebranch/releases

### Install Manually

```bash
# Download (example for macOS ARM)
curl -LO https://github.com/OptionMetrics/terraform-provider-stonebranch/releases/download/v0.4.0/terraform-provider-stonebranch_0.4.0_darwin_arm64.zip

# Unzip
unzip terraform-provider-stonebranch_0.4.0_darwin_arm64.zip

# Install to plugin directory
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/0.4.0/darwin_arm64
mv terraform-provider-stonebranch_v0.4.0 ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/0.4.0/darwin_arm64/
```

### Configure Terraform

```hcl
terraform {
  required_providers {
    stonebranch = {
      source  = "registry.terraform.io/stonebranch/stonebranch"
      version = "0.4.0"
    }
  }
}
```

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
