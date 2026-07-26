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
uses these artifacts to govern its own incremental development.

## Start here

- [PROJECT-VISION.md](PROJECT-VISION.md) defines the product vision, audience,
  principles, scope, and non-goals.
- [AGENTS.md](AGENTS.md) defines the rules for humans and coding agents working
  in this repository.
- [TASK-0000.md](TASK-0000.md) records the repository bootstrap task.
- [docs/OVERVIEW.md](docs/OVERVIEW.md) explains the DaD methodology.
- [docs/DOCUMENTATION.md](docs/DOCUMENTATION.md) defines document numbering,
  references, lifecycle, and templates.

## Repository structure

```text
.
├── cmd/dad/               # CLI entry point
├── internal/              # CLI implementation
├── resources.go           # Embedded authoritative framework resources
├── docs/
│   ├── adr/               # Accepted architectural decisions
│   ├── specs/             # Approved, implementable behavior
│   └── templates/         # Canonical document templates
├── prompts/               # Official agent-neutral prompts
├── examples/              # Progressive adoption examples and case studies
├── PROJECT-VISION.md      # Durable project direction and boundaries
├── AGENTS.md              # Repository-wide implementation instructions
├── TASK-*.md              # Bounded work and completion evidence
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

## License

DaD is licensed under the [Apache License 2.0](LICENSE).
