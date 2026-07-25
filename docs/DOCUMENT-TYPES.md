# DaD Document Types

DaD document types separate questions that are often mixed together. Each type
has a distinct responsibility and should contain only the detail needed to
fulfil it.

Numbering, reference rules, statuses, and templates are defined in
[DOCUMENTATION.md](DOCUMENTATION.md).

## Project vision

**Answers:** Why does the project exist, for whom, and within what durable
boundaries?

The project vision defines purpose, audience, core principles, scope, and
non-goals. It is the highest-level product context and should change rarely.
It does not prescribe individual features, architectural choices, or task
plans.

DaD's vision is [PROJECT-VISION.md](../PROJECT-VISION.md).

## Repository working agreement

**Answers:** How must contributors operate in this repository?

The working agreement contains repository-wide instructions for discovering
context, making changes, verifying work, and respecting scope. It applies to
humans and AI agents. More local instructions may refine it for part of a
repository without silently contradicting higher-level rules.

A working agreement is not product behavior, architectural rationale, or a
task backlog. DaD uses [AGENTS.md](../AGENTS.md) as its portable filename;
tool-specific adapter files may point to the same authority when necessary.

## Architecture Decision Record

**Answers:** Why was a consequential technical choice made, what was chosen,
and what consequences follow?

Create an ADR when a decision:

- constrains multiple changes or components;
- is costly to reverse;
- resolves a meaningful tradeoff or recurring disagreement; or
- must remain understandable after the original discussion is gone.

Do not create an ADR for routine implementation detail, a requirement, a work
plan, or a diary of discussion. An ADR records one decision and its rationale;
it does not specify every behavior produced by that decision.

Accepted ADRs govern applicable architecture. Changing an accepted decision
requires a new ADR that supersedes it.

## Specification

**Answers:** What behavior or contract must an implementation provide?

Create a SPEC when behavior needs a durable, implementation-independent source
of truth, especially for:

- public or cross-component interfaces;
- data formats, protocols, and state transitions;
- invariants, compatibility expectations, and error behavior; or
- complex behavior that multiple tasks will implement or maintain.

A SPEC states observable and testable requirements without becoming a task
plan or architectural history. It may reference ADRs that constrain its
design. Approved SPECs govern applicable behavior.

Avoid a SPEC when acceptance conditions in one TASK completely and durably
describe a small local change.

## Task

**Answers:** What outcome is currently authorized, what are its boundaries, and
how will completion be shown?

Every implementation effort has one active TASK. A TASK defines outcome,
in-scope and out-of-scope work, constraints, governing references, acceptance
conditions, and completion evidence.

A TASK coordinates a bounded unit of change. It does not become the permanent
home for architectural rationale or behavior that must govern future tasks.
Move those concerns to an ADR or SPEC and reference them.

Completed and cancelled TASKs preserve implementation history. They do not
serve as a backlog for unrelated future work.

## Explanatory documentation

**Answers:** How can a reader understand, adopt, or operate the system?

Guides, overviews, and tutorials explain authoritative project content for a
particular reader. They may summarize and link, but they do not silently create
new architectural decisions, behavioral requirements, or task scope.

When explanatory documentation conflicts with its source, the source remains
authoritative and the explanation must be corrected.

## Code and verification

Code is the implementation of approved behavior. Tests and other checks are
evidence about that implementation. They are essential parts of the governance
chain, but they do not replace the documents that explain intent and
authority.

If code, tests, and governing documents disagree, treat the disagreement as
drift to resolve. Do not assume that whichever artifact changed most recently
is automatically correct.

## Choosing a document type

| Need | Use |
| --- | --- |
| Define durable purpose or boundaries | Project vision |
| Define contributor operating rules | Working agreement |
| Preserve why a consequential technical choice governs | ADR |
| Define durable behavior or a contract | SPEC |
| Authorize and bound implementation work | TASK |
| Teach or orient a reader | Explanatory documentation |
| Implement behavior | Code |
| Demonstrate behavior or quality | Tests and other checks |

If one change raises several questions, use the necessary document types and
link them. Do not combine their responsibilities merely to reduce the file
count.

## Relationship rules

The normal direction of authority is:

1. vision and working agreements bound the project and its contributors;
2. accepted ADRs constrain architecture;
3. approved SPECs define behavior within those constraints;
4. ready TASKs authorize implementation under the applicable context;
5. code and tests implement and verify the result.

This order is not a license for a lower layer to reinterpret a higher one.
References make the applicable path explicit, and conflicts return to review
rather than being resolved implicitly during implementation.
