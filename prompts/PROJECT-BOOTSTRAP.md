# Project Bootstrap Prompt

## Purpose

Use this prompt to establish the minimum durable DaD context in a new or
existing repository. Supply any known project goals, constraints, and required
scope with the prompt.

## Prompt

You are bootstrapping this repository for Document-aware Development (DaD).
Your goal is to establish the minimum durable project context needed for
subsequent work. Do not implement product features or anticipated framework
phases.

Follow this process:

1. Inspect the repository before making changes. Read existing contributor
   instructions, project documentation, manifests, source layout, tests, and
   version-control status that can materially affect the bootstrap.
2. Preserve existing project intent and contributor work. Do not replace
   established guidance merely to impose a standard layout.
3. Identify contradictions or missing information that would materially
   change the project's purpose, audience, boundaries, or working rules. Ask
   for human direction when those questions cannot be answered from repository
   evidence.
4. Create or update only the minimum bootstrap documents:
   - a concise project entry point;
   - a project vision defining purpose, audience, scope, principles, and
     non-goals;
   - a repository working agreement that applies to humans and AI agents; and
   - a bounded bootstrap TASK with acceptance conditions and completion
     evidence.
5. Design the intended repository structure at the level needed to guide the
   next task. Do not create empty future-phase artifacts or directories solely
   to match the design.
6. Keep each fact in one authoritative home and link to it elsewhere. Keep
   instructions agent-neutral and usable without AI.
7. Verify that the resulting documents agree with one another and with the
   observable repository state.

Do not create ADRs for decisions that have not been made. Do not create SPECs,
implementation, examples, automation, validation, a CLI, or other tooling
unless the supplied bootstrap scope explicitly authorizes them.

At completion, report:

- documents created or changed;
- material assumptions and unresolved questions;
- checks performed; and
- confirmation that no work outside the bootstrap scope was performed.

Stop after the bootstrap task is complete. Do not begin the next phase.
