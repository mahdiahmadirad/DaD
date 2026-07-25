# TranStudio case study

## Scope and evidence

TranStudio is a local-first, AI-assisted workbench for translating long-form
scholarly documents into Persian while preserving structure, terminology,
traceability, and human review history.

This case study applies the DaD lens to a TranStudio repository snapshot
inspected on 2026-07-26. It is observational: TranStudio has its own
documentation conventions, and this example does not claim that every current
DaD rule was used during its development.

The review examined these representative sources:

- `README.md`, `AGENTS.md`, and `IMPLEMENTATION-PROTOCOL.md`;
- `docs/DOCUMENTATION-CONVENTIONS.md` and the documentation manifest;
- `ADR-0010`, governing the translation execution engine;
- `SPEC-0004`, defining that engine's implementation behavior;
- `TASK-0003`, covering resilient translation jobs and crash recovery;
- the TASK-0003 completion report;
- translation-engine, persistence, application, and UI modules; and
- recovery, performance, boundary, and documentation-consistency tests.

The TranStudio repository was not modified, and its checks were not rerun for
this case study. Test results below are evidence recorded by TranStudio's
completion report.

## Why the project needs durable context

TranStudio's difficult behavior lies at boundaries:

- model calls are slow and can fail outside database transactions;
- a Streamlit rerun must not own durable job state;
- crashes must not lose successful translations or duplicate accepted
  revisions;
- human-approved text must never be overwritten;
- provider, prompt, profile, and model identity must remain traceable; and
- concurrency must not violate SQLite connection or ownership rules.

These constraints span persistence, configuration, translation execution,
application services, UI, and recovery. A task scoped only as "make
translation jobs resumable" would leave implementers to rediscover or invent
critical policy.

## The observed governance chain

TranStudio separates the change across documents and executable evidence.

### Working agreement

`AGENTS.md` defines the repository entry point, required reading, architectural
prohibitions, session workflow, and stop behavior.
`IMPLEMENTATION-PROTOCOL.md` adds one-phase-per-session execution rules,
verification expectations, and a conflict protocol.

Together they make contributor behavior explicit before implementation begins.

### Architectural decision

`ADR-0010` places translation execution behind a model-provider abstraction
and separates orchestration from provider-specific behavior. Its later
resilient-job design establishes a project-local worker, durable state,
ownership, and recovery boundaries.

This is architectural context: it explains constraints that affect several
components and cannot safely be chosen independently by each implementation
task.

### Behavioral specification

`SPEC-0004` turns the decision into implementable behavior. It defines durable
job creation, atomic item acceptance, worker leases, bounded scheduling,
Pause, Resume, Cancel, Retry, state transitions, and the application-facing
service boundary.

The specification makes several subtle requirements reviewable before code:

- no database transaction remains open during a model call;
- executor slots use independent SQLite connections;
- the UI does not launch workers or call the execution engine directly;
- Resume and Retry failed remain distinct operations; and
- exactly one accepted result is guaranteed per job item, not exactly one
  external model invocation.

### Bounded task

`TASK-0003` authorizes resilient translation jobs and crash recovery. It
records global constraints, dependencies, acceptance criteria, verification,
and stop conditions across 28 named subtasks.

Each implementation session names one subtask and stops afterward. This turns
a large cross-component outcome into reviewable units while keeping the shared
constraints visible.

### Implementation and evidence

The resulting repository contains dedicated translation-engine modules for
scheduling, execution, job services, workers, leases, cancellation, recovery,
events, failures, and progress. Persistence migrations and repositories hold
durable job and item state, while application and UI boundaries keep Streamlit
non-authoritative.

The completion report records:

- 298 passing tests in the full suite;
- 156 passing persistence, configuration, and translation-engine tests;
- 76 passing application and UI tests;
- 10 passing recovery and performance tests; and
- 5 passing documentation-consistency and UI-boundary checks.

It also records a MyPy environment blocker rather than claiming a successful
type check. That distinction is important DaD behavior: completion evidence
states what was actually verified and preserves known limitations.

## How drift was constrained

Several practices made the intended architecture difficult to bypass
accidentally:

1. The task names non-negotiable constraints, including additive migrations,
   independent SQLite connections, and preservation of approved text.
2. The specification defines state transitions and ownership rules before
   implementation.
3. Subtasks pair acceptance criteria with targeted verification and explicit
   stop conditions.
4. Boundary tests reject forbidden UI imports instead of relying only on
   reviewer memory.
5. Crash tests exercise stale leases, duplicate launches, and failures around
   atomic commit.
6. The completion report audits architecture boundaries and records the
   remaining environment limitation.

The value is not the number of documents. It is that rationale, behavior,
scope, and evidence have different homes and can be reviewed against one
another.

## Lessons for DaD

### Load relevant context, not all context

TranStudio has a substantial documentation pack: its manifest records 12 ADRs,
10 SPECs, four TASKs, and several reports at the inspected snapshot. A fixed
reading sequence helps onboarding, but implementation work benefits from a
task-specific context path so unrelated documents do not obscure the governing
ones.

DaD therefore requires the active TASK to link the accepted ADRs and approved
SPECs that materially affect it.

### Authority is concern-specific

TranStudio documents both a global source priority and concern ownership. The
case shows why ownership is the more precise model: an ADR answers an
architecture question, a SPEC answers a behavior question, and a TASK answers
a scope and execution question.

DaD does not resolve every disagreement by choosing one globally higher
document. It identifies which artifact owns the disputed concern and pauses
when authorities genuinely conflict.

### Completion evidence must stay connected to the task

TranStudio's dedicated completion report provides unusually strong evidence,
including limitations and a subtask ledger. The maintenance risk is separation:
a report and its TASK can drift unless their relationship remains explicit.

DaD keeps completion evidence in the TASK by default. A separate report is
appropriate when detail warrants it, but the TASK should link it and retain a
concise outcome.

### Stable identifiers matter more than filename style

TranStudio filenames include descriptive slugs, while current DaD conventions
use canonical `TYPE-NNNN.md` filenames. Both approaches rely on stable
identifiers such as `ADR-0010` and `SPEC-0004`.

The transferable lesson is to preserve identifiers and valid references rather
than renumbering history for cosmetic consistency.

### Tooling supports governance but does not define it

TranStudio includes a documentation manifest and tests for document counts,
core files, version alignment, and selected references. Those checks expose
some mechanical drift, while the documents remain understandable without the
checks.

This matches DaD's intended boundary: tools can verify conventions, but the
methodology cannot depend on a particular validator.

## Result

The resilient-jobs work demonstrates the DaD operating loop at meaningful
scale:

1. architectural constraints were made explicit;
2. approved behavior was specified;
3. implementation was divided into bounded work;
4. code and tests supplied evidence; and
5. completion stopped before the next UI phase.

The case also shows where DaD can simplify a mature documentation-driven
project: use concern-specific authority, load only relevant context, keep
completion evidence connected to its task, and treat tooling as support rather
than the source of governance.
