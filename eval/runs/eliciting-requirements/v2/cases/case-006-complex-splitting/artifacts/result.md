# Requirements Elicitation Package

## Understanding

**Requested outcome:** When an operator submits a batch that is currently in validation mode, the batch service must examine every record in the batch and, if any record is found to be invalid, inform the operator of that fact.

**Proposed system boundary:** The batch service is the responsible system. The operator is the human actor. The batch is the data object acted upon. The submission event is the boundary trigger.

**Information sorted by kind:**

| Kind | Source text | Treatment |
| --- | --- | --- |
| State (precondition) | "while a batch is in validation mode" | Scopes the obligation; maps to a `While` precondition |
| Trigger (event) | "when the operator submits it" | A discrete boundary event; maps to `When` trigger |
| Behavioral obligation 1 | "shall validate every record" | Candidate EARS requirement |
| Behavioral obligation 2 | "notify the operator if any record is invalid" | Second obligation, conditionally qualified — split into a separate candidate |
| Unwanted condition | "if any record is invalid" | Qualifies the notification; maps to an `If … then` pattern for Candidate 2 |

**Splitting rationale:** The source joins two responses with `and`. Validation and notification can be verified independently — one concerns whether the service inspects records, the other concerns whether the operator receives a message. The Final Output Invariant requires these to be separate candidates.

The notification response is scoped to an unwanted condition ("if any record is invalid"). The appropriate pattern for that scoped failure response is an `If … then` unwanted-behavior candidate. The precondition (validation mode) and trigger (submission) still apply as context for that candidate.

---

## Clarifying Questions

1. **Definition of "invalid record":** What makes a record invalid? Is there a defined validation rule set, a schema, a business rule, or an external reference? Two implementers may use different criteria without a definition. *(Blocks Candidate 1 and 2 from `ready for review`)*

2. **Scope of "every record":** Does "every record" mean all records currently in the batch at submission time, or only records added since the last validation run? Are deleted or draft records included? *(Blocks Candidate 1)*

3. **Notification content and channel:** What must the notification contain — record identifiers, error codes, a count, or a human-readable message? Through what channel does the operator receive it (UI message, email, system log entry, API response)? *(Blocks Candidate 2)*

4. **Notification granularity:** Does the batch service notify the operator once (summarizing all invalid records) or once per invalid record? *(Blocks Candidate 2)*

5. **Behavior on partial completion:** If validation of one record causes an error (e.g., a dependency failure), does the batch service halt validation, skip that record and continue, or report a system error separately from validation failures? *(Blocks Candidate 1)*

6. **"Validation mode" definition and lifecycle:** How does a batch enter and exit validation mode? Can the operator submit a batch that is not in validation mode, and if so what happens? *(Scoping concern for both candidates)*

7. **Batch identity:** Is "it" (the batch being submitted) uniquely identified by a batch identifier? Is the operator always the same actor who put the batch in validation mode? *(Scoping concern for both candidates)*

---

## Candidate Requirements

### Candidate 1

- **Status:** needs clarification
- **Template:** Complex
- **Requirement:** While a batch is in validation mode, when the operator submits the batch, the batch service shall validate every record in the batch.
- **Source intent and rationale:** The operator submits a batch for processing. Before the batch proceeds, the system must check all records. This is the inspection obligation — it is independently testable from notification.

- **Implementability contract**
  - System boundary: Batch service
  - Responsible system: Batch service
  - Actor: Operator
  - Affected interfaces: The batch service's submission interface (the endpoint or API call that accepts a submission event from the operator)
  - Inputs and data definitions: The submitted batch; the set of records it contains. **Open:** "record" and "invalid record" are undefined — their schema, validation rules, or reference source must be specified.
  - Preconditions, state, feature scope, and trigger: Precondition — the batch is in validation mode. Trigger — the operator submits the batch. **Open:** "validation mode" entry, exit, and persistence are undefined.
  - Observable normal, alternate, and failure response: Normal — the batch service inspects every record. **Open:** behavior when validation of an individual record raises an unexpected error (halt, skip, report separately) is undefined.
  - Quantities, units, limits, timing, ordering, and tolerances: **Open:** no timing constraint is stated. "Every record" implies full coverage, but boundary cases (empty batch, batch with one record) are not addressed.
  - Domain terms, dependencies, and unresolved decisions: "validation mode," "record," "invalid record," "batch," "submits" — all require agreed definitions. See Questions 1, 2, 5, 6, 7.

- **Verification contract**
  - Setup, state, input data, and stimulus: Place a batch containing at least one record into validation mode; have the operator submit it via the submission interface.
  - Observable evidence or oracle: Evidence that the batch service has examined every record — **open:** the observable artifact (audit log entry per record, validation result payload, status field, etc.) is undefined until Q1–Q2 are answered.
  - Expected result: All records in the batch have been evaluated against the applicable validation rules.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — every record receives a validation result. Fail — any record is skipped or the service returns before completing all records. **Open:** "every" cannot be verified without knowing the authoritative record set (Q2).
  - Boundary, negative, and failure cases: Empty batch; batch with a single record; batch with all records valid; batch with all records invalid; record that causes a processing error during validation (Q5).
  - Required data, instrumentation, assessment, or external evidence: Observable validation result or audit evidence per record; known record count for the batch.
  - Observable at the named boundary alone: no
  - If no, internal evidence needed and rationale: The batch service must expose per-record validation results or an audit trace at its output boundary; internal processing completeness is not directly visible from a single submission response without that evidence.

---

### Candidate 2

- **Status:** needs clarification
- **Template:** Complex (with embedded unwanted condition)
- **Requirement:** While a batch is in validation mode, when the operator submits the batch, if any record in the batch is invalid, then the batch service shall notify the operator.
- **Source intent and rationale:** The operator must be informed when validation identifies at least one invalid record. This is the notification obligation. It is triggered only when an unwanted condition exists — at least one record is invalid — and is therefore split from the validation obligation in Candidate 1.

  > **Note on pattern:** This candidate combines a state precondition and event trigger (from the complex form) with an unwanted-condition response (from the `If … then` form). The unwanted condition is nested within the complex scope. The full clause order is: `While <state>, when <trigger>, if <unwanted condition>, then the <system> shall <response>.`

- **Implementability contract**
  - System boundary: Batch service
  - Responsible system: Batch service
  - Actor: Operator (recipient of notification)
  - Affected interfaces: The interface through which the operator receives output from the batch service (UI, API response, message channel — **open:** Q3).
  - Inputs and data definitions: Result of validation from Candidate 1; the set of records determined to be invalid. **Open:** "invalid" is undefined (Q1); the notification content is undefined (Q3).
  - Preconditions, state, feature scope, and trigger: Precondition — the batch is in validation mode. Trigger — the operator submits the batch. Unwanted condition — at least one record in the batch is invalid (as determined after Candidate 1's validation).
  - Observable normal, alternate, and failure response: Normal (unwanted condition present) — the batch service notifies the operator. Alternate (no invalid records) — no notification is required. **Open:** whether a "no invalid records found" confirmation is also sent is not stated. **Open:** notification cardinality (one summary vs. one per invalid record) is undefined (Q4).
  - Quantities, units, limits, timing, ordering, and tolerances: **Open:** timing of notification relative to completion of validation is not specified. "Any record" implies the condition holds for one or more invalid records.
  - Domain terms, dependencies, and unresolved decisions: "invalid record," "notify," "operator" — require agreed definitions. Notification depends on completion of Candidate 1. See Questions 1, 3, 4, 6, 7.

- **Verification contract**
  - Setup, state, input data, and stimulus: Place a batch containing at least one known-invalid record in validation mode; have the operator submit it.
  - Observable evidence or oracle: A notification is delivered to the operator through the defined channel (open: channel undefined, Q3). The notification references the invalid record(s) (open: content undefined, Q3).
  - Expected result: The operator receives a notification after submission when at least one record is invalid.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — a notification is delivered to the operator after submission and before the batch service's submission response is considered complete. Fail — no notification is delivered, or the notification does not identify invalid records. **Open:** timing threshold and observable delivery mechanism are undefined.
  - Boundary, negative, and failure cases: Batch with exactly one invalid record; batch with all records invalid; batch with all records valid (no notification expected); operator not available to receive notification (delivery failure path).
  - Required data, instrumentation, assessment, or external evidence: Observable notification artifact at the operator interface; confirmation of delivery when the channel requires it.
  - Observable at the named boundary alone: no
  - If no, internal evidence needed and rationale: Whether the notification was generated and delivered requires evidence at the operator-facing boundary (e.g., a message received, a UI alert displayed, or an API response containing validation errors). If the operator interface is asynchronous, delivery confirmation may require additional instrumentation.

---

## Suggested Supporting Requirements

1. **Failure-path companion (Candidate 1):** While a batch is in validation mode, when the operator submits the batch, if validation of one or more records cannot be completed due to a service error, then the batch service shall report the failure to the operator and preserve the batch in its current state.
   - *Source behavior:* Candidate 1 requires validating every record. A service error during validation is a material failure path not addressed in the source.

2. **Boundary companion (Candidate 1):** While a batch is in validation mode, when the operator submits a batch that contains no records, the batch service shall [accept the submission without validation / reject the submission — ask user to decide].
   - *Source behavior:* "Every record" implies a non-empty set. An empty batch is a boundary case whose expected behavior is undefined.

3. **Required companion (Candidate 2):** When the batch service notifies the operator of invalid records, the batch service shall include the identifier and reason for each invalid record in the notification.
   - *Source behavior:* The notification in Candidate 2 is vacuous unless it carries enough information for the operator to act. This companion defines minimum notification content; it depends on Q1 and Q3.

4. **Operational companion (both candidates):** While a batch is in validation mode, when the operator submits the batch, the batch service shall complete validation of all records within [T seconds — value TBD] of receiving the submission.
   - *Source behavior:* No timing constraint is stated. A large batch could make the submission operation unbounded; a latency bound may be required for operational use.

5. **Optional consideration:** Where audit logging is included, the batch service shall record the validation result for each record, including the batch identifier, record identifier, validation outcome, and timestamp.
   - *Source behavior:* Not stated in the source, but batch validation in regulated or high-volume contexts often requires an audit trail. Ask whether auditability is in scope.

---

## Consistency Findings

| Finding | Kind | Detail | Resolution options |
| --- | --- | --- | --- |
| C1 & C2 ordering dependency | Gap | Candidate 2 (notification) depends on the output of Candidate 1 (validation). The source does not state whether notification must wait until all records are validated or may be sent incrementally. | Clarify ordering: add an explicit sequencing clause, or confirm that notification is issued only after validation of every record is complete (Q4). |
| "If any record is invalid" scope | Precedence question | "Any" means one or more. The condition is met when at least one invalid record exists, but the source does not say whether the batch service must notify once (aggregate) or per record. Candidate 2 is ambiguous on cardinality. | Ask Q4 and update Candidate 2 to specify cardinality. |
| Validation mode scope | Gap | Neither candidate addresses what happens when a batch is not in validation mode and the operator submits it. An implicit assumption is that submission outside validation mode is handled elsewhere, but a conflict could exist if another requirement mandates different behavior. | Confirm that submission outside validation mode is out of scope for these candidates, or add a scoping note to both candidates. |
| Empty batch | Gap (boundary) | Neither candidate defines the response when the batch contains no records. "Every record" and "any record is invalid" are vacuously satisfied for an empty set, which could lead to inconsistent implementation. | Resolve via the boundary companion suggestion (item 2 above). |

---

## Assumptions

1. **Assumption A1 (provisional):** "Validation mode" is a well-defined, distinguishable state of the batch that the system can determine at submission time. If this state is implementation-specific or ambiguous, both candidates require revision.

2. **Assumption A2 (provisional):** "Submits" refers to a discrete, atomic boundary event — a single operator action that signals the batch is ready for processing. If submission is a multi-step or asynchronous process, both candidates require revision.

3. **Assumption A3 (provisional):** "The operator" who receives the notification in Candidate 2 is the same operator who performed the submission. If notifications may be routed to a different actor (e.g., a batch administrator), Candidate 2 requires revision.

4. **Assumption A4 (provisional):** The phrase "notify the operator" in the source implies a positive, observable delivery of information to the operator, not merely logging to a system record. If system logging satisfies the intent, the verification contract for Candidate 2 must be revised.

---

## Rationale and Context

The source request expresses a data-quality gate: an operator should not unknowingly progress a batch containing invalid records. The validation step (Candidate 1) and the notification step (Candidate 2) serve this gate together. Splitting them preserves independent negotiability — a team may decide to implement validation with a different notification channel or timing without coupling both obligations.

The complex EARS pattern is appropriate for both candidates because both a state (validation mode) and a discrete trigger (submission) are genuinely required to scope each obligation. Neither condition alone is sufficient: the operator submitting a batch outside validation mode, or a batch being in validation mode without a submission event, should not trigger these behaviors.

The unwanted condition in Candidate 2 ("if any record is invalid") is nested within the complex scope rather than standing alone, because the notification is only meaningful in the context of a submission event on a batch in validation mode. A standalone `If … then` requirement without that context could be triggered outside the intended scope.

---

## Review Status

| Candidate | Status | Blocking items |
| --- | --- | --- |
| Candidate 1 | needs clarification | Q1 (invalid record definition), Q2 (record scope), Q5 (partial-completion behavior), Q6 (validation mode definition) |
| Candidate 2 | needs clarification | Q1 (invalid record definition), Q3 (notification content and channel), Q4 (notification cardinality), Q6 (validation mode definition) |

Neither candidate is `ready for review`. Both implementability and verification contracts contain open items that would allow two reasonable implementers or verifiers to choose different behavior. Resolution of the clarifying questions is required before either candidate can be advanced.
