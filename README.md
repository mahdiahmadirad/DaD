# DaD

DaD (Document-aware Development) is an open-source framework for governing
software implementation with durable project context.

It is intended to reduce implementation drift, preserve architectural intent,
and keep software understandable and maintainable when work is performed by
people, AI coding agents, or both.

DaD is agent-neutral. Its conventions must not depend on a particular model,
vendor, editor, or automation platform.

## Status

DaD currently provides its documentation model, methodology, official prompts,
adoption examples, and a cross-platform command-line interface. The repository
uses these artifacts to govern its own incremental development. Versioned CLI
archives are distributed through tag-driven GitHub Releases.

## Start here

- [PROJECT-VISION.md](PROJECT-VISION.md) defines the product vision, audience,
  principles, scope, and non-goals.
- [AGENTS.md](AGENTS.md) defines the rules for humans and coding agents working
  in this repository.
- [TASK-0000](docs/tasks/TASK-0000.md) records the repository bootstrap task.
- [docs/OVERVIEW.md](docs/OVERVIEW.md) explains the DaD methodology.
- [docs/DOCUMENTATION.md](docs/DOCUMENTATION.md) defines document numbering,
  references, lifecycle, and templates.

## Repository structure

```text
.
├── cmd/dad/               # CLI entry point
├── internal/              # CLI implementation
├── .github/workflows/     # Native CI and tag-driven release automation
├── .goreleaser.yaml       # Cross-platform release assembly
├── CHANGELOG.md           # Curated user-facing release history
├── resources.go           # Embedded authoritative framework resources
├── docs/
│   ├── adr/               # Accepted architectural decisions
│   ├── specs/             # Approved, implementable behavior
│   ├── tasks/             # Bounded work and completion evidence
│   └── templates/         # Canonical document templates
├── prompts/               # Official agent-neutral prompts
├── examples/              # Progressive adoption examples and case studies
├── PROJECT-VISION.md      # Durable project direction and boundaries
├── AGENTS.md              # Repository-wide implementation instructions
├── go.mod
├── LICENSE
```

This structure separates intent, decisions, specifications, implementation,
and verification. Directories are introduced only when approved work requires
them.

## CLI

The CLI is implemented in Go as one codebase for Windows, Linux, and macOS. It
operates locally without Git, network access, or an AI service.

Build and verify it with a Go 1.25 or newer toolchain:

```text
go build ./cmd/dad
go test ./...
```

The executable provides:

```text
dad init
dad new
dad list
dad status
dad context
dad check
dad prompt
```

Run `dad --help` for command syntax. The authoritative CLI contracts are
[SPEC-0001](docs/specs/SPEC-0001.md) and
[SPEC-0002](docs/specs/SPEC-0002.md).

## Installation

GitHub Releases is the primary distribution channel. After a release is
published, download the archive for your platform from the
[DaD releases page](https://github.com/mahdiahmadirad/DaD/releases).

The current prerelease is
[v0.2.0-rc.1](https://github.com/mahdiahmadirad/DaD/releases/tag/v0.2.0-rc.1).
Choose the archive matching your platform:

```text
dad_VERSION_linux_amd64.tar.gz
dad_VERSION_linux_arm64.tar.gz
dad_VERSION_darwin_amd64.tar.gz
dad_VERSION_darwin_arm64.tar.gz
dad_VERSION_windows_amd64.zip
```

Extract the archive, place `dad` or `dad.exe` in a directory on your `PATH`,
and verify the installed version:

```text
dad --version
```

Every release includes `dad_VERSION_checksums.txt`. Verify the downloaded
archive with a platform SHA-256 tool, or verify its GitHub build provenance
with the GitHub CLI:

```text
gh attestation verify dad_VERSION_OS_ARCH.EXT \
  -R mahdiahmadirad/DaD
```

Release archives are initially unsigned and not notarized. Windows or macOS
may therefore show reputation or quarantine warnings; review the checksum and
provenance rather than disabling operating-system security globally.

Users with a compatible Go toolchain may instead build a specific tagged
version from source:

```text
go install github.com/mahdiahmadirad/DaD/cmd/dad@vVERSION
```

This is a source build, not a download of the published binary.

## Release development

Release behavior is governed by
[ADR-0002](docs/adr/ADR-0002.md) and
[SPEC-0003](docs/specs/SPEC-0003.md). Maintainers curate
[CHANGELOG.md](CHANGELOG.md) and validate the release configuration without
publishing by running:

```text
goreleaser check
goreleaser release --snapshot --clean
go run ./internal/releasetool verify 0.0.0-snapshot dist
```

A semantic `v*` tag triggers native Windows, Linux, and macOS validation.
Successful validation creates a draft GitHub Release, verifies and attests its
assets, and publishes the draft. Tags and published assets are never replaced;
fixes receive a new version.

## License

DaD is licensed under the [Apache License 2.0](LICENSE).
