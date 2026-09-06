# DaD

DaD (Document-aware Development) is an open-source framework for governing
software implementation with durable, versioned project context.

It is intended to reduce implementation drift, preserve architectural intent,
and keep software understandable and maintainable when work is performed by
people, AI coding agents, or both.

DaD is agent-neutral. Its conventions do not depend on a particular model,
vendor, editor, repository host, or automation platform.

## Why DaD?

A contributor can make a locally reasonable change and still move a project
away from its intended direction.

The missing information is often not more source code. It is the context around
the code: why a boundary exists, which alternative was rejected, what behavior
is approved, what the current task excludes, and what evidence is required
before the work is complete. AI-assisted development makes this failure mode
faster and more visible, but it does not create it.

DaD keeps these concerns close to implementation while giving each one a
distinct role:

| Artifact | Question it answers |
| --- | --- |
| Project vision | Why does this project exist, and where are its boundaries? |
| ADR | Why was a consequential technical decision made? |
| SPEC | What behavior or interface has been approved? |
| TASK | What bounded change is authorized now, and how will it be verified? |
| Code | How is the approved behavior implemented? |
| Tests and checks | What repeatable evidence supports the result? |

The goal is not to maximize documentation. Detail should be proportional to the
decision and its risk, and context should be required only when it can guide
implementation or review. The authoritative model is described in
[PROJECT-VISION.md](PROJECT-VISION.md) and
[docs/DOCUMENTATION.md](docs/DOCUMENTATION.md).

## What DaD is—and is not

DaD provides:

- a repository-level documentation and governance model;
- explicit relationships between intent, decisions, specifications, tasks,
  implementation, and evidence;
- official agent-neutral prompts and incremental adoption examples; and
- a local CLI for initializing the structure, creating and listing documents,
  assembling task-specific context, and checking structural consistency.

DaD is not an AI coding agent, an orchestration service, a replacement for
testing or code review, or a guarantee that an implementation is correct. It
does not give generated documentation authority over reviewed decisions or
executable behavior.

## Relationship to related frameworks

DaD shares the intent-first concerns of projects such as
[GitHub Spec Kit](https://github.github.com/spec-kit/),
[OpenSpec](https://github.com/Fission-AI/OpenSpec), and
[BMAD](https://docs.bmad-method.org/), but starts from a narrower question:
**what durable, repository-held context should govern a change before, during,
and after implementation?**

Those projects provide structured specification or broader AI-assisted delivery
workflows. DaD focuses on the roles, authority, lifecycle, traceability, and
proportionality of project artifacts. Its CLI supports that model; it does not
run an end-to-end planning or implementation workflow.

DaD should therefore not be assumed to replace, integrate with, or be
incompatible with any of these projects. They evolve independently, and
coexistence has not yet been verified. The links above point to their own
documentation for current behavior.

## See it in practice

[DaD Sample — Text Analysis API](https://github.com/mahdiahmadirad/DaD-sample)
is a small ASP.NET Core project with one explicit architectural constraint: the
application must remain independent from any specific AI provider.

You can inspect the governing chain without installing anything:

1. Read the sample's
   [project vision](https://github.com/mahdiahmadirad/DaD-sample/blob/main/PROJECT-VISION.md)
   for its purpose and boundaries.
2. Read
   [ADR-0001](https://github.com/mahdiahmadirad/DaD-sample/blob/main/docs/adr/ADR-0001.md)
   for the provider-independence decision.
3. Read
   [SPEC-0001](https://github.com/mahdiahmadirad/DaD-sample/blob/main/docs/specs/SPEC-0001.md)
   for the approved provider and API behavior.
4. Read
   [TASK-0001](https://github.com/mahdiahmadirad/DaD-sample/blob/main/docs/tasks/TASK-0001.md)
   for the bounded implementation scope and completion evidence.
5. Compare those artifacts with the
   [implementation](https://github.com/mahdiahmadirad/DaD-sample/tree/main/src)
   and
   [tests](https://github.com/mahdiahmadirad/DaD-sample/tree/main/tests).

The sample currently demonstrates the governed path from intent to evidence. It
is planned to evolve through deliberate drift and reconciliation, but those
stages are not presented as completed evidence yet.

## Essays

The ideas behind DaD are developed in a bilingual essay series:

1. [When Building Becomes Easier Than Understanding](https://mehdiahmadirad.me/en/articles/building-easier-than-understanding/)
   — [وقتی ساختن آسان‌تر از فهمیدن می‌شود](https://mehdiahmadirad.me/fa/articles/building-easier-than-understanding/)
2. [A Project Should Be Able to Explain Itself](https://mehdiahmadirad.me/en/articles/project-should-explain-itself/)
   — [پروژه باید بتواند خودش را توضیح دهد](https://mehdiahmadirad.me/fa/articles/project-should-explain-itself/)
3. [Building a Project with Document-Aware Development](https://mehdiahmadirad.me/en/articles/building-a-project-with-dad/)
   — [ساختن یک پروژه با Document-Aware Development](https://mehdiahmadirad.me/fa/articles/building-a-project-with-dad/)

## Status

DaD currently provides its documentation model, methodology, official prompts,
adoption examples, and a cross-platform command-line interface. The repository
uses these artifacts to govern its own incremental development. Versioned CLI
archives are distributed through tag-driven GitHub Releases.

The current release is a prerelease. DaD is ready for inspection and
experimentation, but its conventions and CLI contract may still change before
1.0.

## Start here

- [PROJECT-VISION.md](PROJECT-VISION.md) defines the product vision, audience,
  principles, scope, and non-goals.
- [AGENTS.md](AGENTS.md) defines the rules for humans and coding agents working
  in this repository.
- [docs/OVERVIEW.md](docs/OVERVIEW.md) explains the DaD methodology.
- [docs/DOCUMENTATION.md](docs/DOCUMENTATION.md) defines document numbering,
  references, lifecycle, and templates.
- [TASK-0000](docs/tasks/TASK-0000.md) records the repository bootstrap task.

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

GitHub Releases is the primary distribution channel. Download the archive for
your platform from the
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
