# Medium example: Relay

Relay is a fictional service that accepts shipment events from partner
webhooks. This example shows DaD when one change needs the full chain from
architectural rationale to behavioral contract and bounded implementation.

The governing path is:

1. [PROJECT-VISION.md](PROJECT-VISION.md) sets durable product boundaries.
2. [ADR-0001](docs/adr/ADR-0001.md) accepts durable receipt before processing.
3. [SPEC-0001](docs/specs/SPEC-0001.md) defines the webhook contract.
4. [TASK-0001](TASK-0001.md) authorizes one implementation change.
5. [AGENTS.md](AGENTS.md) tells contributors how to follow that path.

The content is fictional and contains no runnable service. The ADR and SPEC
show states reached after human review; the TASK remains `Ready`, and its
completion evidence is intentionally empty.

This is a medium example because behavior will be maintained across more than
one component and the durability choice has consequences beyond a single
task. Smaller changes should not reproduce this document set automatically.
