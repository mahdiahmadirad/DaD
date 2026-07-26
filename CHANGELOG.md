# Changelog

All notable user-facing changes to DaD are recorded in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and release versions follow [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- The agent-neutral DaD documentation model, workflow, templates, official
  prompts, and progressive adoption examples.
- The cross-platform `dad` CLI for repository initialization, governed
  document creation, inventory, context resolution, structural checks, and
  official prompt access.
- Tag-driven CI and release assembly for versioned Windows, Linux, and macOS
  archives with SHA-256 checksums and build provenance.

### Changed

- Development builds now report `dev`; release and versioned module builds
  obtain their semantic version from build metadata instead of a maintained
  source constant.
