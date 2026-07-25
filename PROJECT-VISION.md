# Project Vision

## Purpose

DaD (Document-aware Development) is a framework for governing software
implementation through explicit, versioned project context.

Software teams make decisions across tasks, conversations, reviews, code, and
documentation. When that context is incomplete or disconnected from
implementation, contributors can produce locally reasonable changes that
conflict with the system's intended direction. AI-assisted development makes
this failure mode faster and more visible, but it does not create it.

DaD aims to make relevant intent discoverable before implementation, keep work
bounded while it is performed, and leave an understandable record after it is
complete. This repository must apply those practices to its own development
and serve as the reference implementation.

## Scope

DaD will define a practical way to:

- capture durable project intent and constraints;
- express bounded implementation work;
- connect architectural decisions and approved behavior to implementation;
- give contributors the context needed for the work in front of them;
- detect or expose divergence between intent and implementation;
- preserve traceability without making documentation an end in itself; and
- support review and maintenance by people who did not perform the original
  work.

The framework may eventually include conventions, schemas, validation, and
supporting tools. Their form is deliberately left to later approved phases.

## Target audience

DaD is for software teams that want implementation to remain aligned with
explicit project intent, including:

- maintainers governing long-lived or evolving systems;
- engineers collaborating with one or more AI coding agents;
- teams using different agents, editors, and automation platforms;
- open-source projects accepting contributions from people with limited prior
  context; and
- teams working without AI that still need clear decisions, bounded tasks, and
  maintainable implementation history.

The primary user is the contributor preparing, implementing, reviewing, or
maintaining a change. Framework authors and tool integrators are secondary
users.

## Core principles

### Intent precedes implementation

The reason, boundaries, and acceptance conditions for a change should be clear
before code is written. Detail should be proportional to the decision and the
risk.

### Context must be actionable

Documentation should change how work is performed or reviewed. Context that
cannot guide a decision should not become mandatory process.

### One concern, one authoritative home

Vision, decisions, specifications, tasks, code, and tests serve different
purposes. Each fact should have a clear source of authority, with links used
instead of competing copies.

### Work is bounded and incremental

Changes should be small enough to review against explicit acceptance
conditions. Future work must not be implemented merely because it is
foreseeable.

### Traceability is a means, not the product

A maintainer should be able to understand why behavior exists and what governs
it. Traceability should support that understanding with the least durable
process needed.

### Verification belongs with implementation

Claims about behavior should be backed by proportionate, repeatable evidence.
Documentation and implementation are incomplete when they materially
contradict each other.

### Humans remain responsible

AI agents can plan, implement, and review, but project authority remains in
versioned repository artifacts and human governance. Generated output is held
to the same standards as any other contribution.

### Agent neutrality

Core semantics must be portable across Codex, Claude Code, Goose, Cursor,
Cline, future agents, and human-only workflows. Vendor-specific integrations
may adapt the framework but must not define it.

### The repository practices what it defines

DaD's own development is its first test. Rules that cannot guide this
repository clearly and economically should not be imposed on adopters.

## Non-goals

DaD is not:

- an AI coding agent, model, editor, or orchestration service;
- a replacement for version control, issue tracking, testing, code review, or
  engineering judgment;
- a method for generating large amounts of documentation;
- a guarantee that an implementation is correct, secure, or valuable;
- a mandate for one programming language, architecture, delivery process, or
  repository host;
- a system that gives generated text authority over reviewed project
  decisions and executable behavior;
- a requirement that teams use AI; or
- an attempt to specify every future feature before evidence demands it.

## Product standard

DaD succeeds when a contributor can determine what governs a change, implement
within its boundaries, verify the result, and leave the repository easier for
the next person or agent to understand.

The framework should prefer a small number of strong, composable conventions
over a large prescriptive process. Adoption should be incremental, and useful
behavior should not depend on a single tool.
