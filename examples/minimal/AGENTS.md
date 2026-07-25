# Pulse Working Agreement

These instructions apply to humans and AI coding agents working on Pulse.

Before changing the project:

1. Read `PROJECT-VISION.md`.
2. Read the requested `TASK-NNNN.md`.
3. Inspect the affected code and tests.

Implement only a TASK in `Ready` or `In Progress` status. Do not add behavior
outside its scope. If a missing architectural decision or durable behavioral
contract is discovered, stop and request the appropriate ADR or SPEC rather
than inventing it in code.

Verify every acceptance condition, record only checks actually performed, and
mark the TASK `Complete` only when all conditions pass. Stop after the
requested TASK.
