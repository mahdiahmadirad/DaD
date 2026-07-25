# DaD Principles

The [project vision](../PROJECT-VISION.md) defines DaD's durable core
principles. This document translates them into working rules for applying the
methodology.

## Repository context outranks conversational context

Conversations are useful for exploration, but they are incomplete, difficult
to discover, and often unavailable to the next contributor. Any context that
must govern future work belongs in a reviewed repository artifact with a clear
owner and lifecycle.

This does not mean preserving every discussion. Preserve the resulting
decision, requirement, constraint, or task boundary.

## Load context before changing the system

A contributor should understand the active task, applicable decisions,
approved behavior, and affected implementation before making changes.
Implementation started from partial context can be technically sound and still
be wrong for the project.

The required context is the smallest set that can materially affect the work,
not the entire documentation tree.

## Give each concern one authority

Project direction, architectural rationale, behavioral contracts,
implementation scope, code, and verification answer different questions.
Information belongs in the document type responsible for that question.

Other documents link to the authority instead of maintaining competing
versions. When two authoritative sources conflict, work pauses until the
conflict is resolved explicitly.

## Document decisions, not narration

Documentation must help a contributor decide, implement, review, or maintain.
DaD favors concise statements of outcome, rationale, behavior, constraints, and
evidence over diaries of activity or generated summaries.

The amount of documentation should track the durability and risk of the
information. A reversible local choice may need no ADR; a system-wide
constraint usually does.

## Bound work before implementation

Every implementation effort has an explicit outcome, scope, exclusions,
constraints, and acceptance conditions. These boundaries protect the project
from opportunistic expansion and make review meaningful.

Foreseeable work is not automatically current work. Discoveries outside the
boundary are recorded or proposed separately rather than implemented
incidentally.

## Change intent explicitly

Implementation is allowed to reveal missing knowledge. It is not allowed to
turn that discovery into silent policy.

When a new architectural choice is required, decide it through an ADR. When
approved behavior must change, update or supersede the relevant SPEC. When the
desired outcome or scope changes, revise or replace the TASK through review.
Only then does implementation continue under the new context.

## Verification closes the loop

Acceptance conditions define what completion means before work begins.
Verification then provides proportionate evidence that the result satisfies
those conditions and still respects governing decisions and specifications.

Passing tests alone is insufficient when the tests do not cover the intended
outcome. Conversely, documentation alone is not evidence that behavior exists.

## Preserve history without preserving mistakes

Accepted decisions, approved contracts, and completed tasks form useful
project memory. When they cease to govern current work, their lifecycle status
and references make that explicit rather than erasing the path that led to the
current system.

History informs maintenance; current accepted and approved documents govern
new implementation.

## Keep the methodology agent-neutral

DaD describes inputs, authority, workflow, and evidence without depending on a
model, prompt format, editor, or orchestration platform. Integrations may make
the methodology easier to use, but they cannot redefine its document semantics
or approval authority.

Humans and AI agents follow the same task boundaries and quality standard.
Responsibility remains with the people and governance process that accept the
change.

## Practice proportionality

The smallest sufficient governance is the correct amount. Use a TASK for
bounded work, add an ADR only for a consequential architectural decision, and
add a SPEC only when behavior needs an explicit contract.

If maintaining an artifact costs more than the clarity or control it provides,
reduce or remove the process rather than generating more content.
