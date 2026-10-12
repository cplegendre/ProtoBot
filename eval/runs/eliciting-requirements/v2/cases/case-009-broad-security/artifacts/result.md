# Requirements Elicitation Package

## Understanding

**Source request:** "The system must not contain security vulnerabilities."

**Requested outcome:** The requester wants assurance that the system does not expose exploitable weaknesses that could be used to compromise confidentiality, integrity, availability, or other security properties.

**Proposed system boundary:** Unspecified. "The system" names no component, service, interface, or deployment scope. This is a goal statement, not a bounded behavioral obligation.

**Information sort:**

| Kind | Source text | Treatment |
|---|---|---|
| Goal | "The system must not contain security vulnerabilities." | Preserve as motivation; do not promote directly to a normative requirement without decomposition. |
| Constraint | Implied: no exploitable weaknesses at release | Scope, classes, severity, and assessment method are undefined. |
| Assumption | "Must" is treated as a normative intent, not a shall-strength obligation as stated. | Marked provisional. |
| Design idea | None present. | — |
| Behavioral obligation | None yet stated with observable boundary or response. | Candidates require clarification on threat scope, vulnerability classes, assessment method, severity threshold, and release gate. |

**Why this cannot be directly promoted to a ready requirement:**
The source sentence names no responsible system, no affected interface or actor, no defined vulnerability class or threat scope, no assessment method, and no measurable pass/fail criterion. Two reasonable implementers or verifiers reading the sentence would choose entirely different behaviors. Per the readiness gate, the sentence fails both contracts.

---

## Clarifying Questions

The following questions are the highest-value blockers. Each unanswered question keeps one or more candidates at `needs clarification`.

**Q1 — System boundary and responsible system:**
Which system, component, or set of components is in scope? Examples: the public-facing web API, the authentication service, the data-processing pipeline, the client application, the entire deployed product. A single requirement can name only one responsible system.

**Q2 — Threat scope and vulnerability classes:**
Which classes of vulnerability are in scope? Examples: injection flaws (SQL, command, LDAP), broken authentication, insecure deserialization, cross-site scripting, sensitive data exposure, insecure direct object references, known vulnerable dependencies, server-side request forgery, supply-chain weaknesses. Is the scope limited to a recognized taxonomy (e.g., OWASP Top 10, CWE Top 25, a specific CVE severity range)?

**Q3 — Severity threshold:**
What severity level triggers a release-blocking finding? Examples: any CVSS score ≥ 7.0 (High/Critical), any Critical finding regardless of score, any finding in a defined class regardless of CVSS. Who decides the threshold?

**Q4 — Assessment method and evidence:**
How is the absence of vulnerabilities demonstrated? Examples: static analysis (SAST) tool scan with no findings above threshold, dynamic analysis (DAST) of a running instance, dependency-vulnerability scan (SCA) against a known-vulnerability database, penetration-test report from a named assurance level, code review by a credentialed reviewer, formal audit. Which methods are required and which are optional?

**Q5 — Assessment frequency and lifecycle gate:**
When must the assessment be performed and evidence provided? Examples: before every production release, on every pull request merged to the main branch, continuously in CI/CD, on a defined periodic schedule, after every significant feature change. Does a finding discovered after release require a remediation deadline?

**Q6 — Remediation and residual-risk behavior:**
When a vulnerability is found, what must the system or development process do? Examples: block the release until remediated, record the finding in a tracking system, obtain a documented exception with an owner and deadline, apply a compensating control and document it. What happens when a finding cannot be remediated before a defined deadline?

**Q7 — Scope of "contain":**
Does "contain" mean the system binary/artifact, its runtime behavior, its dependencies, its configuration, its deployment environment, or all of the above?

**Q8 — Actor and authorization model:**
Is there an actor (e.g., an external attacker, an authenticated user, a privileged internal process) against whom the system must be protected? Some vulnerability requirements are scoped to specific attacker capabilities.

---

## Candidate Requirements

### Candidate 1

- **Status:** needs clarification
- **Template:** unresolved
- **Requirement:** Blocked — pattern cannot be selected without guessing. The source provides no named responsible system, no defined vulnerability class, no assessment method, no severity threshold, and no observable response that can be verified at a system boundary.
- **Source intent and rationale:** The requester wants the system to be free of exploitable security weaknesses. This is a valid and important security assurance goal. It does not yet describe a testable behavioral obligation.
- **Implementability contract**
  - System boundary: **open** — Q1 unanswered.
  - Responsible system: **open** — Q1 unanswered.
  - Actor: **open** — Q8 unanswered; no attacker profile or user role is defined.
  - Affected interfaces: **open** — no interface (API, UI, data store, dependency chain) is named.
  - Inputs and data definitions: **open** — no input, payload, or artifact is scoped.
  - Preconditions, state, feature scope, and trigger: **open** — no lifecycle gate, state, or triggering event is defined.
  - Observable normal, alternate, and failure response: **open** — Q2, Q3, Q4, Q6 unanswered; "does not contain vulnerabilities" is not observable at a boundary.
  - Quantities, units, limits, timing, ordering, and tolerances: **open** — Q3 (severity threshold), Q5 (frequency and deadline) unanswered.
  - Domain terms, dependencies, and unresolved decisions: "security vulnerability" is undefined; "system" is undefined; "contain" is undefined in scope (Q7).
- **Verification contract**
  - Setup, state, input data, and stimulus: **open** — no assessment method defined (Q4); no lifecycle gate defined (Q5).
  - Observable evidence or oracle: **open** — no defined tool, report format, or evidence artifact.
  - Expected result: **open** — no threshold or pass criterion defined (Q3).
  - Pass/fail criteria, including scope, thresholds, units, and timing: **open** — severity threshold (Q3) and assessment scope (Q2, Q7) unanswered.
  - Boundary, negative, and failure cases: **open** — remediation and residual-risk path undefined (Q6).
  - Required data, instrumentation, assessment, or external evidence: **open** — Q4 unanswered.
  - Observable at the named boundary alone: **no** — security assurance typically requires internal evidence (scan results, reports) that cannot be observed at an external functional boundary without defined instrumentation.
  - If no, internal evidence needed and rationale: Evidence would include tool-generated reports (SAST, DAST, SCA), penetration-test findings, or code-review records — none of which are defined.

---

*Once Q1–Q8 are answered, the single broad source statement will decompose into multiple narrower, independently verifiable requirements. Illustrative decomposition directions are shown in Suggested Supporting Requirements below.*

---

## Suggested Supporting Requirements

These suggestions illustrate how the source goal decomposes once clarifying answers are received. None is an approved requirement. Each relates to the source statement and identifies the missing behavior and reason.

**SR-1 — Dependency vulnerability scan (Required companion)**
*Related to:* Candidate 1 (scope of "contain" covering third-party dependencies).
*Missing behavior:* A requirement that the build process scans declared dependencies against a known-vulnerability database and blocks release when a dependency has an unmitigated finding at or above the agreed severity threshold.
*Reason:* "The system" routinely includes third-party libraries. Without a dependency-scan obligation, the source goal cannot be met even if the first-party code is clean.

**SR-2 — Static analysis gate (Required companion)**
*Related to:* Candidate 1 (first-party code).
*Missing behavior:* A requirement that static analysis is run on source code and that the pipeline does not promote a build when findings above the agreed threshold are present.
*Reason:* Static analysis is a common, reproducible method for detecting injection, memory-safety, and logic-error classes. Without it, first-party code assurance is undefined.

**SR-3 — Known-vulnerability remediation deadline (Failure-path companion)**
*Related to:* Candidate 1 (what happens when a vulnerability is found after release).
*Missing behavior:* A requirement that, when a vulnerability of or above the agreed severity is confirmed in a released version, the responsible system or team shall produce a remediated release or a documented compensating control within a defined period.
*Reason:* A pre-release gate alone does not address newly disclosed vulnerabilities in shipped software. Without a post-release remediation obligation, the goal has no lifecycle coverage.

**SR-4 — Authentication and authorization boundary (Interface companion)**
*Related to:* Candidate 1 (attacker-model scope, Q8).
*Missing behavior:* A requirement that the system enforces access control at each named external interface, denying requests that lack valid credentials or authorization.
*Reason:* Authentication and authorization failures are among the most common vulnerability classes. If the source goal covers this class, a separate behavioral requirement at the interface boundary is needed.

**SR-5 — Security assessment evidence record (Operational companion)**
*Related to:* Candidate 1 (assessment evidence, Q4 and Q5).
*Missing behavior:* A requirement that the assessment evidence (scan reports, penetration-test findings, review records) is retained for a defined period and available for audit.
*Reason:* Without a retention obligation, compliance with the source goal cannot be retrospectively verified.

**SR-6 — Exception and residual-risk process (Boundary companion)**
*Related to:* Candidate 1 (when a finding cannot be remediated before a deadline).
*Missing behavior:* A requirement that accepted residual risk is documented with an owner, a compensating control description, and a review date.
*Reason:* No realistic system achieves zero findings continuously. Without a documented exception process, the source goal creates an unachievable absolute obligation with no defined handling for edge cases.

**SR-7 — Scope definition (Optional consideration)**
*Related to:* Candidate 1 (Q7 — what "contain" means).
*Missing behavior:* An explicit scope statement listing which artifacts, environments, and configurations are covered (e.g., production deployment only, container images, infrastructure-as-code, CI/CD pipeline configuration).
*Reason:* Without explicit scope, each team may interpret coverage differently, creating gaps and disputes. Whether this belongs in a requirement or a separate document depends on how the host manages scope.

---

## Consistency Findings

**Finding 1 — Gap: No measurable boundary between goal and obligation**
The source statement is a goal. No behavioral obligation at a named system boundary has been stated. Until Q1–Q5 are answered and the goal is decomposed into obligations with measurable responses, the candidate set contains only one unresolved item covering all security assurance. This is a significant gap: verification is impossible.

**Finding 2 — Gap: Pre-release and post-release lifecycle not covered**
Even when clarified, a single pre-release assessment requirement leaves post-release vulnerability disclosure unaddressed. SR-3 flags this gap. The candidate set will need both a release-gate requirement and a remediation-obligation requirement to be internally consistent.

**Finding 3 — Precedence question: Absolute prohibition vs. risk-acceptance process**
The source phrase "must not contain" implies an absolute prohibition. SR-6 notes that an exception process is typical for practical systems. If both an absolute prohibition and an exception process are included in the final requirement set, a precedence rule is needed: under what authority and conditions may residual risk be accepted, and does acceptance satisfy the prohibition?

No conflicts or duplicates are present in the current (single-candidate) set. Additional findings will emerge when the candidate set is expanded after answers to Q1–Q8 are received.

---

## Assumptions

**A1 — "Must" is treated as normative intent.**
The source uses "must not" rather than "shall not." Per the operating rules, `must` in source language is interpreted as normative intent and translated to `shall` in EARS candidates. This assumption is provisional; the requester should confirm that the obligation is normative.

**A2 — "The system" refers to a software system under active development.**
The source is interpreted as applying to software artifacts and their runtime behavior, not to physical infrastructure or organizational process alone. If physical or organizational scope is intended, Q1 should be answered to reflect that.

**A3 — No specific taxonomy is assumed.**
No particular vulnerability taxonomy (OWASP, CWE, CVE, NIST NVD) has been assumed as the intended scope. Q2 must be answered before a taxonomy is referenced in a normative sentence.

---

## Rationale and Context

The source statement "The system must not contain security vulnerabilities" is a common, well-intentioned security policy goal. It fails both the implementability and verification contracts for the following reasons, each grounded in the readiness-gate rules:

1. **No named responsible system.** A requirement must name exactly one responsible system. "The system" is a placeholder.
2. **"Security vulnerability" is undefined.** The term spans dozens of distinct classes with different detection methods, severity scales, and mitigations. Without a defined scope, two verifiers cannot agree on what constitutes a passing result.
3. **No observable response.** "Does not contain" is a state of an artifact, not an observable system response at a boundary. A verifiable requirement must describe what the system produces, returns, denies, records, or signals that a verifier can observe.
4. **No assessment method.** Whether the goal is met cannot be determined without knowing how it is to be assessed (tool, method, evidence artifact, evaluation standard).
5. **No pass/fail criterion.** Without a severity threshold, no evaluator can declare pass or fail.

The skill reference explicitly names this sentence as a canonical example of a `needs clarification` case: *"The system must not contain security vulnerabilities. Ask for threat scope, vulnerability classes, assessment method, severity threshold, and release boundary. Decompose it into observable behaviors and scoped assurance criteria."*

This package faithfully applies that guidance.

---

## Review Status

| Item | Status | Blocking question(s) |
|---|---|---|
| Candidate 1 | **needs clarification** | Q1, Q2, Q3, Q4, Q5, Q6, Q7, Q8 |
| SR-1 through SR-6 | Suggestions — not candidates | Pending answers to clarifying questions |
| SR-7 | Optional consideration — not a candidate | Whether scope belongs in a requirement or a separate document |
| Consistency Finding 1 | Gap — blocks all verification | Resolved when Candidate 1 is decomposed |
| Consistency Finding 2 | Gap — post-release lifecycle | Resolved when remediation obligation is added |
| Consistency Finding 3 | Precedence question | Resolved when absolute-prohibition vs. exception-process intent is confirmed |

**Overall package status: needs clarification.**
No candidate requirement is ready for review. The source statement is a valid and important security assurance goal; it requires decomposition into bounded, observable, and measurable behavioral obligations before any candidate can advance.
