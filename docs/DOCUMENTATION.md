# Documentation Model

DaD separates durable intent, decisions, behavior, and implementation work so
that each concern has one authoritative home.

| Document | Purpose | Location |
| --- | --- | --- |
| Project vision | Defines durable direction and boundaries | `PROJECT-VISION.md` |
| ADR | Records why a consequential technical decision was made | `docs/adr/ADR-NNNN.md` |
| SPEC | Defines approved, externally observable or implementation-significant behavior | `docs/specs/SPEC-NNNN.md` |
| TASK | Bounds a unit of work and records its outcome | `docs/tasks/TASK-NNNN.md` |

Templates live in `docs/templates/` and are not project documents. A document
should be created only when its content serves the purpose above; not every
task requires a new ADR or SPEC.

## Numbering rules

1. ADRs, SPECs, and TASKs each use an independent sequence of four-digit,
   zero-padded numbers.
2. The canonical identifier and filename are `TYPE-NNNN` and `TYPE-NNNN.md`.
   The document title follows the identifier inside the file, not in the
   filename.
3. Allocate the next number greater than every number already used for that
   document type. Numbers are never reused or changed, including after a
   document is rejected, cancelled, or superseded.
4. Gaps are allowed. Concurrent contributors must confirm that a number is
   still available before integration and resolve collisions by renumbering
   the document that has not yet been integrated.
5. `TASK-0000` is the repository bootstrap task. New sequences otherwise begin
   at `0001`.
6. Templates and general guidance are unnumbered.

These rules favor stable references over chronological perfection. The version
control history remains the authority for authorship and dates.

## Reference rules

- In Markdown, reference a project document with a relative link whose label
  contains its canonical identifier, for example
  `[SPEC-0004](../specs/SPEC-0004.md)`.
- Use repository-relative paths when a format cannot resolve Markdown links.
  Do not use machine-local absolute paths or vendor-specific URLs as durable
  references.
- Link to the authoritative document instead of copying its normative content.
  A summary is acceptable only when clearly subordinate to the linked source.
- A TASK links to every ADR and SPEC that governs its implementation. A SPEC
  links to ADRs that explain its design constraints. An ADR links to earlier
  decisions it supersedes or materially depends on.
- References must describe the relationship when it is not obvious from
  context. A bare list of identifiers is not a substitute for explaining why
  they matter.
- References normally point from the dependent document to its authority.
  Do not maintain backlinks solely for symmetry.
- When a referenced document is superseded, retain the historical link and add
  the replacement where it governs current work.
- Keep links valid when moving content. Canonical identifiers must remain
  visible even if a repository later changes its directory layout.

Code and tests should reference document identifiers only when the connection
is not evident from normal structure and naming. Documentation references
explain intent; they do not replace readable code or executable verification.

## Document lifecycle

Every ADR, SPEC, and TASK has exactly one status from its lifecycle below.
Status changes are reviewable repository changes. The person or process
authorized to accept work in the adopting project approves transitions; an AI
agent does not grant authority by generating a status change.

### ADR lifecycle

- `Proposed` → `Accepted` or `Rejected`
- `Accepted` → `Superseded`

- **Proposed:** under review and not authoritative.
- **Accepted:** the decision governs applicable work.
- **Rejected:** considered but not adopted.
- **Superseded:** replaced by a later accepted ADR, linked in both records.

After acceptance, preserve the decision and rationale as a historical record.
Correct minor errors in place, but use a new ADR to change the decision.

### SPEC lifecycle

- `Draft` → `Approved` or `Withdrawn`
- `Approved` → `Superseded` or `Retired`

- **Draft:** under development and not authoritative.
- **Approved:** governs applicable behavior.
- **Withdrawn:** abandoned before approval.
- **Superseded:** replaced wholly by a later approved SPEC.
- **Retired:** no longer governs current behavior and has no replacement.

An approved SPEC may be clarified in place only when observable meaning does
not change. A behavioral change requires review; use a new SPEC when preserving
the earlier contract or migration history matters.

### TASK lifecycle

- `Proposed` → `Ready` → `In Progress` → `Complete`
- `Proposed`, `Ready`, or `In Progress` → `Cancelled`
- `In Progress` ↔ `Blocked`

- **Proposed:** scope and acceptance conditions are being prepared.
- **Ready:** approved and sufficiently defined to begin.
- **In Progress:** active implementation or verification is underway.
- **Blocked:** work started but cannot currently proceed; the blocker is
  recorded.
- **Complete:** acceptance conditions are verified and completion evidence is
  recorded.
- **Cancelled:** intentionally closed without completing the outcome; the
  reason is recorded.

A completed or cancelled TASK is not reopened. Follow-up work receives a new
TASK and references the earlier one.

## Templates

- [ADR template](templates/ADR-TEMPLATE.md)
- [SPEC template](templates/SPEC-TEMPLATE.md)
- [TASK template](templates/TASK-TEMPLATE.md)
