# DaD

DaD (Document-aware Development) is an open-source framework for governing
software implementation with durable project context.

It is intended to reduce implementation drift, preserve architectural intent,
and keep software understandable and maintainable when work is performed by
people, AI coding agents, or both.

DaD is agent-neutral. Its conventions must not depend on a particular model,
vendor, editor, or automation platform.

## Status

DaD is in its bootstrap phase. This repository currently defines the project's
scope, principles, working rules, and intended structure. It does not yet
implement the framework.

## Start here

- [PROJECT-VISION.md](PROJECT-VISION.md) defines the product vision, audience,
  principles, scope, and non-goals.
- [AGENTS.md](AGENTS.md) defines the rules for humans and coding agents working
  in this repository.
- [TASK-0000.md](TASK-0000.md) records the repository bootstrap task.

## Intended repository structure

The structure below is the target organization, not a list of artifacts that
already exist. Directories are introduced only when an approved task requires
them.

```text
.
├── README.md              # Project entry point
├── PROJECT-VISION.md      # Durable project direction and boundaries
├── AGENTS.md              # Repository-wide implementation instructions
├── TASK-*.md              # Bounded units of work and their outcomes
├── LICENSE
├── docs/
│   ├── adr/               # Accepted architectural decisions
│   └── specs/             # Approved, implementable behavior
├── src/                   # Framework implementation
├── tests/                 # Executable verification
├── tools/                 # Project-owned development utilities
└── examples/              # Focused adoption examples
```

This structure separates intent, decisions, specifications, implementation,
and verification. The exact contents and language-specific layout will be
decided by later tasks rather than assumed during bootstrap.

## License

DaD is licensed under the [Apache License 2.0](LICENSE).
