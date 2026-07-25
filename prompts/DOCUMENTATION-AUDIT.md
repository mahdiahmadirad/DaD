# Documentation Audit Prompt

## Purpose

Use this prompt to evaluate whether project documentation is authoritative,
connected, current, and useful for implementation. Supply the repository scope
and any known change or risk that should focus the audit.

## Prompt

Audit the documentation in the supplied scope against DaD's document model and
the observable repository state. The audit is read-only unless an active TASK
explicitly authorizes fixes.

Begin by reading the repository working agreement, project vision,
documentation rules, and relevant accepted ADRs, approved SPECs, active or
completed TASKs, implementation, and tests. Inspect version-control status and
do not treat uncommitted contributor work as established project policy.

Evaluate:

- whether each document has a clear responsibility and authoritative home;
- lifecycle status, numbering, naming, and location;
- references to governing and superseding documents;
- broken, missing, circular, or unexplained relationships;
- contradictions between vision, ADRs, SPECs, TASKs, code, and tests;
- behavior or constraints implemented without durable authority;
- stale explanatory documentation and completion evidence;
- duplicated normative content likely to drift; and
- missing documentation only where its absence creates a concrete
  implementation or maintenance risk.

Do not score documentation by volume, formatting preference, or template
completeness alone. Do not recommend a new ADR, SPEC, or TASK unless the
information belongs in that document type and has an identifiable consumer.

Report findings in priority order. For each finding include:

- severity: `Critical`, `High`, `Medium`, or `Low`;
- the affected files or identifiers;
- direct evidence;
- the authority or DaD rule involved;
- the practical risk; and
- the smallest appropriate corrective action.

Separate confirmed defects from uncertain observations. Include a short
summary of the scope inspected, checks performed, limitations, and areas with
no material findings.

Do not edit files, create governance documents, or implement fixes unless the
supplied task explicitly requires it. Stop after the audit report.
