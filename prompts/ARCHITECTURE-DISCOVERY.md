# Architecture Discovery Prompt

## Purpose

Use this prompt to build an evidence-based understanding of an existing
system's architecture before making decisions or implementation changes.
Supply the repository area or question that should bound the discovery.

## Prompt

Perform an architecture discovery of the supplied repository scope. This is a
read-only investigation unless an active TASK explicitly authorizes a named
documentation artifact.

Before analysis:

1. Read the repository working agreement, project vision, active TASK if one
   exists, and all accepted ADRs or approved SPECs relevant to the scope.
2. Inspect the version-control status and preserve unrelated contributor work.
3. Define the discovery boundary. Do not expand into unrelated subsystems.

Examine repository evidence sufficient to understand:

- system boundaries and external actors;
- runtime components and their responsibilities;
- dependency direction and important module boundaries;
- primary control flows and data flows;
- public, cross-component, and external interfaces;
- persistence, configuration, deployment, and operational assumptions;
- cross-cutting concerns such as security, failure handling, observability,
  concurrency, and compatibility where relevant;
- tests or checks that reveal intended behavior; and
- areas where implementation and documentation disagree.

For every material claim, distinguish:

- **Observed:** directly supported by repository evidence;
- **Inferred:** the most plausible explanation, with the supporting evidence
  and uncertainty stated; or
- **Unknown:** not determinable from the available repository context.

Do not infer historical rationale from current code. Do not convert an
observed pattern into an accepted architectural decision. Do not create or
change ADRs, SPECs, or implementation without explicit authorization.

Produce a concise discovery report containing:

1. scope and sources inspected;
2. architecture summary;
3. components and responsibilities;
4. key dependency, control, and data flows;
5. interfaces and system boundaries;
6. observed constraints and quality attributes;
7. documentation or implementation conflicts;
8. unknowns, risks, and confidence limits; and
9. proposed follow-up decisions or documentation work, clearly marked as
   proposals.

Cite repository paths and identifiers near the claims they support. Prefer a
small number of useful relationships over an exhaustive file inventory.

Stop after reporting the discovery. Do not implement recommendations.
