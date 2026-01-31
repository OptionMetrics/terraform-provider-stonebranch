# Release Process

This document describes how to build and publish the Terraform Provider for StoneBranch to the JFrog Artifactory registry.

## Prerequisites

### Required Tools

- [Go](https://golang.org/doc/install) >= 1.24
- [GoReleaser](https://goreleaser.com/install/) >= 2.0
- [JFrog CLI](https://jfrog.com/getcli/) (required for publishing)

Install on macOS:

```bash
brew install goreleaser
brew install jfrog-cli
```

### JFrog Artifactory Authentication

The publish process uses JFrog CLI for uploads. Authenticate before publishing:

```bash
jf login
```

This opens a browser for authentication. Once complete, your credentials are stored and used automatically for uploads.

To verify your authentication:

```bash
jf config show
```

## Versioning

This project uses [Semantic Versioning](https://semver.org/) via git tags. The version is automatically derived from the latest git tag.

### Version Format

- Tags must follow the format `v{major}.{minor}.{patch}` (e.g., `v0.3.0`)
- The `v` prefix is stripped when creating release artifacts (e.g., `v0.3.0` → `0.3.0`)
- Without a tag, the version defaults to `0.0.0-dev`

### Check Current Version

```bash
make version
```

## Complete Release Workflow

### Step 1: Ensure Clean Git State

**Important:** GoReleaser requires a clean git state. All changes must be committed before creating a tag.

```bash
# Check for uncommitted changes
git status

# If there are changes, commit them
git add .
git commit -m "Your commit message"
```

### Step 2: Create Version Tag

Only create the tag after ALL changes are committed:

```bash
make tag V=0.4.0
```

This creates an annotated tag `v0.4.0` at the current commit.

### Step 3: Authenticate with JFrog

If you haven't already authenticated (or your session expired):

```bash
jf login
```

### Step 4: Build and Publish

```bash
make publish
```

This will:
1. Build binaries for all platforms using GoReleaser
2. Create zip archives with proper naming
3. Generate SHA256 checksums
4. Upload all artifacts to Artifactory using `jf` CLI

### Step 5: Push to Remote

After successful publish, push commits and tag to GitHub:

```bash
git push origin main
git push origin v0.4.0
```

### Step 6: Verify Upload

```bash
jf rt search "terraform-providers/stonebranch/stonebranch/0.4.0/"
```

## Quick Reference

```bash
# Full release workflow (copy-paste version)
git status                          # Ensure clean state
git add . && git commit -m "msg"    # Commit any changes (if needed)
make tag V=0.4.0                    # Create version tag
jf login                            # Authenticate (if needed)
make publish                        # Build and upload
git push origin main                # Push commits
git push origin v0.4.0              # Push tag
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make version` | Show current version from git tag |
| `make release-snapshot` | Build without tag (for testing) |
| `make release` | Build release artifacts (no upload) |
| `make publish` | Build and upload to Artifactory |
| `make tag V=x.y.z` | Create new version tag |

## Build Artifacts

GoReleaser creates the following in `dist/`:

```
dist/
├── terraform-provider-stonebranch_0.4.0_darwin_amd64.zip
├── terraform-provider-stonebranch_0.4.0_darwin_arm64.zip
├── terraform-provider-stonebranch_0.4.0_linux_amd64.zip
├── terraform-provider-stonebranch_0.4.0_linux_arm64.zip
├── terraform-provider-stonebranch_0.4.0_windows_amd64.zip
└── terraform-provider-stonebranch_0.4.0_SHA256SUMS
```

## Artifactory Location

Artifacts are published using Terraform filesystem mirror structure:

```
terraform-providers/registry.terraform.io/stonebranch/stonebranch/{version}/{os}_{arch}/
    terraform-provider-stonebranch_{version}_{os}_{arch}.zip
```

Example:
```
terraform-providers/registry.terraform.io/stonebranch/stonebranch/0.4.0/darwin_arm64/
    terraform-provider-stonebranch_0.4.0_darwin_arm64.zip
```

## Using the Published Provider

### Option 1: Filesystem Mirror (Recommended)

Download the provider tree and configure Terraform to use it:

```bash
# Create local mirror directory
mkdir -p ~/terraform-providers

# Download all versions (or specific version)
jf rt download "terraform-providers/registry.terraform.io/" ~/terraform-providers/ --flat=false
```

Configure `~/.terraformrc`:

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/Users/yourname/terraform-providers"
    include = ["registry.terraform.io/stonebranch/stonebranch"]
  }

  direct {
    exclude = ["registry.terraform.io/stonebranch/stonebranch"]
  }
}
```

Then in your Terraform config:

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

### Option 2: Download Manually

```bash
# Download specific platform
jf rt download \
  "terraform-providers/registry.terraform.io/stonebranch/stonebranch/0.4.0/darwin_arm64/*.zip" \
  ./

# Unzip and install to plugin cache
unzip terraform-provider-stonebranch_0.4.0_darwin_arm64.zip
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/0.4.0/darwin_arm64
mv terraform-provider-stonebranch_v0.4.0 \
  ~/.terraform.d/plugins/registry.terraform.io/stonebranch/stonebranch/0.4.0/darwin_arm64/
```

## Troubleshooting

### "git is in a dirty state"

GoReleaser requires all changes to be committed before building a release.

```bash
# Check what's dirty
git status

# Commit changes
git add .
git commit -m "message"

# Recreate tag at new commit
git tag -d v0.4.0
make tag V=0.4.0
```

### "tag was not made against commit"

The tag points to a different commit than HEAD. Recreate the tag:

```bash
git tag -d v0.4.0
make tag V=0.4.0
```

### Tag Already Exists on Remote

If you need to update a tag that's already pushed:

```bash
git push origin v0.4.0 --force
```

**Warning:** Only force-push tags if the release hasn't been used by others.

### Authentication Errors (401)

Re-authenticate with JFrog:

```bash
jf login
```

### Repository Not Found (404/405)

- Verify the repository `terraform-providers` exists in Artifactory
- Check you have write permissions

### Build Failures

```bash
# Check GoReleaser configuration
goreleaser check

# Run with debug output
goreleaser release --clean --snapshot --skip=publish --debug
```

## Configuration Files

| File | Purpose |
|------|---------|
| `.goreleaser.yaml` | GoReleaser build configuration |
| `Makefile` | Build automation targets |
| `dist/` | Build output directory (gitignored) |
