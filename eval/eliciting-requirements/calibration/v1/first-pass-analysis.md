<!-- markdownlint-disable MD013 -->

# Provisional First-Level Calibration Pass

This is a local, AI-assisted first pass for PR #240. It is **not** the
independent human calibration required by issue #63.

## Evidence

- Development run: `first-pass-dev-v2`
- Regression run: `first-pass-regression-v2`
- Local model: `enmaas/gemini-3.8-flash`
- Development artifacts: `.local-eval-runs/eliciting-requirements/`
  `first-pass-dev-v2/`
- Regression artifacts: `.local-eval-runs/eliciting-requirements/`
  `first-pass-regression-v2/`
- All 15 selected cases produced one collected `artifacts/result.md`.
- Original baseline artifacts were not present locally, so these are fresh local
  reproductions and are not directly comparable to the recorded v1 model run.

## Automated checks

Development deterministic judges:

- Execution, package structure, EARS template coverage, and critical-category
  gate: 11/11 passed.
- Readiness gate: 10/11 passed.
- EARS structure: 10/11 passed.
- The two failures are both `case-012-host-independent`: no ready candidate in
  the annotated mixed-readiness case, and a normative sentence with more than
  one `shall`.

Regression deterministic judges:

- 24/24 checks passed across four cases and six judges.

## Provisional reviewer-a scores

| Case | Semantic | Discipline | False-ready | Critical | First-level observation |
| --- | ---: | ---: | --- | --- | --- |
| case-001-ubiquitous | 3 | 3 | No | No | Correctly separates retention and retrieval, but uses a complex retrieval pattern and invents a purge policy. |
| case-002-event-driven | 4 | 3 | No | No | Correct event-driven handling; adds unrequested payment-provider, idempotency, and security requirements. |
| case-003-state-driven | 4 | 3 | No | No | Correct paused-state treatment; speculative queue and timing companions. |
| case-004-optional-feature | 4 | 3 | No | No | Correct optional-feature and CSV behavior; speculative packaging and transfer policy. |
| case-005-unwanted-behavior | 3 | 3 | No | No | Preserves timeout, 503, and Retry-After; reopens specified behavior and invents timing tolerances. |
| case-006-complex-splitting | 4 | 4 | No | No | Good split between validation and invalid-record notification. |
| case-007-rough-api | 2 | 2 | No | Yes | Avoids technology leakage in candidates but invents substantial acceptance, tracking, processing, and retrieval behavior. |
| case-008-vague-ui | 4 | 4 | No | No | Correctly blocks vague drafting and asks for measurable scope and criteria. |
| case-011-conflicting-set | 5 | 5 | No | No | Clearly identifies the conflict without deciding for the user. |
| case-012-host-independent | 3 | 3 | No | No | Preserves host metadata but chooses event-driven over the expected complex pattern and misses mixed readiness. |
| case-023-ready-complete | 5 | 4 | No | No | Correctly preserves the complete ready contract; companion suggestions are somewhat speculative. |
| case-014-uncertainty | 3 | 3 | No | No | Preserves uncertainty and gates the candidate, but adds assumptions and overstates pattern ambiguity. |
| case-015-implementation-temptation | 3 | 3 | No | No | Separates Redis/GraphQL/TTL from requirements but asserts a cache conflict and newest-first behavior. |
| case-020-referent-and-vacuous | 1 | 2 | No | Yes | Flags ambiguity but overstates R2 as inherently contradictory and proposes an invented guest-page policy. |
| case-024-ready-regression | 4 | 4 | No | No | Preserves endpoint, exact response, timing, P100, concurrency, and ready status. |

## Provisional aggregates

| Partition | Cases | Mean semantic quality | Mean review discipline | Critical cases | False-ready cases |
| --- | ---: | ---: | ---: | ---: | ---: |
| Development | 11 | 3.73 | 3.36 | 1 | 0 |
| Regression | 4 | 2.75 | 3.00 | 1 | 0 |
| Combined | 15 | 3.47 | 3.27 | 2 | 0 |

## Interpretation

The generated packages generally follow the requested structure and avoid direct
implementation leakage in their primary candidate requirements. The main first-
level concerns are speculative supporting requirements and assumptions, plus two
substantive quality failures:

1. `case-007-rough-api` turns many design ideas and undefined behaviors into
   proposed requirements, reducing semantic fidelity and review discipline.
2. `case-020-referent-and-vacuous` correctly spots ambiguity but resolves the
   unauthenticated-dashboard statement too aggressively instead of preserving
   the ambiguity as a clarification question.

`case-012-host-independent` also needs attention: it retains host metadata but
does not satisfy the annotated mixed-readiness expectation and triggers both
deterministic judge failures.

## Calibration status

- `scores.yaml` contains only this provisional reviewer-a pass.
- `sample.yaml` was not modified.
- Reviewer-b remains unassigned and unscored.
- `agreement.yaml` and `adjudication.yaml` remain pending.
- `artifact.yaml` remains untrusted and was not changed.
- No baseline was promoted, and no confirmed failures were added to the corpus.
- No GitHub commit or push was performed.
