# DaD Overview

DaD (Document-aware Development) is a methodology for keeping software
implementation aligned with explicit project intent. It gives contributors a
small, versioned set of documents that answers:

- Why does this project and change exist?
- Which decisions and behaviors govern the work?
- What is the current task allowed to change?
- How will completion be demonstrated?

DaD applies whether implementation is performed by humans, AI coding agents,
or both. Its authority comes from reviewed repository content, not from a
particular tool or conversation.

## The problem

Implementation drift occurs when locally reasonable work diverges from the
system's intended direction. Common causes include missing context, decisions
held only in conversations, ambiguous behavior, expanding task scope, and
verification that is disconnected from the original outcome.

DaD addresses this by making the context needed for implementation durable,
specific, and discoverable. It does not attempt to document everything.

## The governance chain

DaD organizes project knowledge by responsibility:

1. The [project vision](../PROJECT-VISION.md) defines durable purpose and
   boundaries.
2. The [working agreement](../AGENTS.md) defines how contributors operate in
   the repository.
3. Architecture Decision Records explain why consequential technical choices
   govern the system.
4. Specifications define approved behavior and contracts.
5. Tasks authorize bounded implementation work and define completion.
6. Code implements the approved behavior.
7. Tests and recorded completion evidence show what was verified.

This is a chain of context, not a requirement to create every document for
every change. A task references only the accepted decisions and approved
specifications that materially govern it. Small work may need only a clear
task; consequential work may need the full chain.

## The operating loop

DaD work follows a short feedback loop:

1. **Frame:** define the outcome, boundaries, constraints, and acceptance
   conditions.
2. **Load context:** read the governing project documents, decisions,
   specifications, and affected implementation.
3. **Resolve gaps:** obtain a decision or specification when implementation
   would otherwise invent durable policy or behavior.
4. **Implement:** make the smallest coherent change within the approved task.
5. **Verify:** gather proportionate, repeatable evidence against the acceptance
   conditions.
6. **Reconcile:** ensure documents, implementation, and evidence agree.
7. **Close:** record the outcome and stop at the task boundary.

When implementation reveals that the governing context is wrong or incomplete,
the loop returns to the appropriate document. The contributor does not silently
encode a new decision in code.

The detailed process is defined in [WORKFLOW.md](WORKFLOW.md).

## Authority and responsibility

Different artifacts have authority over different questions. Vision governs
direction, ADRs govern architectural decisions, SPECs govern approved
behavior, and TASKs govern the scope of current work. Code and tests provide
the implemented reality and verification evidence, but their existence does
not silently override an accepted decision or approved specification.

People remain accountable for approving direction, decisions, behavior, and
work. AI agents may prepare, analyze, implement, and verify changes, but they
operate under the same repository rules and do not grant authority to their
own output.

## Minimal adoption

A project can adopt DaD incrementally:

1. State the project vision and repository working agreement.
2. Require bounded tasks with explicit acceptance conditions.
3. Record consequential architectural decisions as ADRs.
4. Add SPECs where behavior needs an implementation-independent contract.
5. Require completion evidence and reconcile drift when it is found.

The methodology is useful before automation exists. Tools may later assist
discovery or verification, but the document model and workflow must remain
understandable and usable without them.

## Further reading

- [PRINCIPLES.md](PRINCIPLES.md) explains the rules behind the methodology.
- [WORKFLOW.md](WORKFLOW.md) defines the end-to-end working process.
- [DOCUMENT-TYPES.md](DOCUMENT-TYPES.md) explains the responsibility of each
  document type.
- [DOCUMENTATION.md](DOCUMENTATION.md) defines numbering, references,
  lifecycle, and templates.
