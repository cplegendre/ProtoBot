# Requirements Elicitation Package

## Understanding

**Requested outcome:** When a client submits an order that meets the defined validity criteria, the order service must durably identify that order and communicate the identifier back to the client within the same submission interaction.

**Proposed system boundary:** The order service, at its order-submission interface (the endpoint or API surface that receives order submissions from clients).

**Actors:** A client (an external caller submitting the order).

**Source language note:** The source uses `shall` twice in one sentence joined by `and` ("shall assign … and return"). These are two separately testable obligations and are split into Candidate 1 and Candidate 2. The source also provides an inline definition of "valid order" — this is a constraint/input definition, not a separate behavioral obligation; it is captured in the implementability contracts and as a boundary companion suggestion.

---

## Clarifying Questions

1. **Identifier uniqueness and format:** Must the assigned order identifier be unique across all orders (e.g., globally unique, scoped to the service, or scoped to the client)? Is there a required format, length, or encoding (e.g., UUID, numeric sequence)?

2. **Identifier persistence:** Does the assigned identifier represent a durable record — i.e., must the order service persist the order and its identifier, or is the identifier only required in the response?

3. **Invalid order behavior:** The source defines what constitutes a valid order but does not specify what the service shall do when an invalid order is submitted (missing product identifier, zero or negative quantity, empty product list). What shall the service return in that case?

4. **Duplicate or idempotent submission:** If a client submits the same order more than once (same products and quantities, or same client-supplied correlation token), shall the service assign a new identifier each time, or return the existing identifier?

5. **"Positive quantity" definition:** Does "positive quantity" mean a whole-number (integer) value greater than zero, or are fractional quantities permitted?

6. **Product identifier validity:** Must the order service verify that each product identifier refers to a known product, or is presence of a non-empty identifier sufficient for validity?

---

## Candidate Requirements

### Candidate 1

- **Status:** candidate
- **Template:** Event-driven
- **Requirement:** When a client submits a valid order, the order service shall assign a unique order identifier to that order.
- **Source intent and rationale:** The source states the order service shall assign an order identifier upon receipt of a valid order. Assignment is separated from return (Candidate 2) because the two obligations are independently testable and could be implemented or verified separately.
- **Implementability contract**
  - System boundary: The order service, at its order-submission interface.
  - Responsible system: The order service.
  - Actor: A client submitting an order.
  - Affected interfaces: The order-submission interface (inbound).
  - Inputs and data definitions: A valid order — an order containing at least one product identifier and a positive quantity for each product identifier present. "Positive quantity" is assumed to mean an integer greater than zero (see Clarifying Question 5). "Product identifier" is assumed to be a non-empty value; whether it must resolve to a known product is unresolved (see Clarifying Question 6).
  - Preconditions, state, feature scope, and trigger: Trigger — the client submits an order that satisfies the validity criteria.
  - Observable normal, alternate, and failure response: Normal — the order service assigns a unique order identifier. Alternate/failure — behavior for invalid orders is unresolved and addressed in Candidate 3 (suggested).
  - Quantities, units, limits, timing, ordering, and tolerances: Identifier uniqueness scope is unresolved (see Clarifying Question 1). No timing constraint specified in the source.
  - Domain terms, dependencies, and unresolved decisions: "Unique" scope (Q1); "positive quantity" — integer vs. fractional (Q5); product identifier validity check (Q6); persistence of identifier (Q2).
- **Verification contract**
  - Setup, state, input data, and stimulus: Submit a well-formed order to the order service containing at least one product identifier and a quantity greater than zero for each.
  - Observable evidence or oracle: The order service produces an order identifier that was not previously assigned to a different order.
  - Expected result: A distinct order identifier is assigned to this order submission.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — an identifier is produced and differs from all previously assigned identifiers (within the uniqueness scope once Q1 is resolved). Fail — no identifier is produced, or the identifier duplicates a previously assigned one.
  - Boundary, negative, and failure cases: (a) Order with exactly one product and quantity of 1 — must produce an identifier. (b) Uniqueness across multiple sequential submissions — each must receive a different identifier.
  - Required data, instrumentation, assessment, or external evidence: Access to the set of previously assigned identifiers (or a reliable equality check between two independent submissions).
  - Observable at the named boundary alone: yes

---

### Candidate 2

- **Status:** candidate
- **Template:** Event-driven
- **Requirement:** When a client submits a valid order, the order service shall return the assigned order identifier in the submission response.
- **Source intent and rationale:** The source explicitly states the identifier shall be "returned in the submission response." This is the client-visible delivery of the identifier and is distinct from the assignment act (Candidate 1): a system could assign without returning, or return without a durable assignment — both paths must be independently required and verified.
- **Implementability contract**
  - System boundary: The order service, at its order-submission interface (outbound response).
  - Responsible system: The order service.
  - Actor: A client submitting an order.
  - Affected interfaces: The order-submission interface (outbound response to the client).
  - Inputs and data definitions: Same validity criteria as Candidate 1. The identifier returned must be the same identifier assigned per Candidate 1.
  - Preconditions, state, feature scope, and trigger: Trigger — the client submits a valid order. Precondition — the order service has assigned an identifier (Candidate 1).
  - Observable normal, alternate, and failure response: Normal — the submission response includes the assigned order identifier. Failure path — behavior for invalid orders is unresolved (see Candidate 3 suggestion and Clarifying Question 3).
  - Quantities, units, limits, timing, ordering, and tolerances: The identifier must appear in the same response as the submission acknowledgement (synchronous, same interaction). Format is unresolved (Q1).
  - Domain terms, dependencies, and unresolved decisions: "Submission response" — assumed to be the direct, synchronous reply to the client's submission call. If the interaction is asynchronous, this requirement needs revision.
- **Verification contract**
  - Setup, state, input data, and stimulus: Submit a valid order to the order service and capture the response.
  - Observable evidence or oracle: The submission response body (or equivalent response artefact) contains an order identifier field with a non-empty value matching the identifier assigned per Candidate 1.
  - Expected result: The response contains the assigned order identifier.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — the response includes a non-empty identifier that matches the one assigned (verifiable by subsequently retrieving the order by that identifier, if a retrieval interface exists). Fail — no identifier in the response, or the identifier does not match the assigned one.
  - Boundary, negative, and failure cases: (a) Minimum valid order (one product, quantity of 1) — identifier must be present in response. (b) Invalid order submission — must not return a valid identifier (behavior subject to Q3 resolution).
  - Required data, instrumentation, assessment, or external evidence: The response payload at the order-submission interface. If identifier–order binding must be verified, a secondary query (e.g., order lookup by identifier) provides additional evidence.
  - Observable at the named boundary alone: yes

---

## Suggested Supporting Requirements

### SR-1 — Failure-path companion (relates to Candidate 1 and Candidate 2)

**Missing behavior:** The source defines what a valid order is but does not specify what the order service shall do when a client submits an invalid order (e.g., no product identifiers, a quantity of zero or below, or a negative quantity).

**Suggested candidate (draft):** `If a client submits an order that does not contain at least one product identifier with a positive quantity for each, then the order service shall reject the submission and return an error response indicating that the order is invalid.`

**Reason:** Without this companion, a verifier cannot determine whether invalid submissions produce no response, an error, or (incorrectly) a valid identifier. An implementer would have to guess. This is a required failure-path for completeness.

**Status:** needs clarification — exact error response content and format are unspecified; confirmation that this behavior is in scope is needed.

---

### SR-2 — Boundary companion (relates to Candidate 1)

**Missing behavior:** The source does not specify whether the assigned order identifier must be unique, and if so, within what scope (global, per-service instance, per-client, per-session).

**Suggested candidate (draft):** `The order service shall assign an order identifier that is unique among all order identifiers previously assigned by the order service.`

**Reason:** Uniqueness is implied by the purpose of an "identifier" but is not stated. Without a uniqueness requirement, two independent implementations may assign conflicting identifiers, making orders indistinguishable. This companion makes the uniqueness scope normative and verifiable.

**Status:** needs clarification — uniqueness scope must be confirmed (Q1).

---

### SR-3 — Optional consideration (relates to Candidate 1 and Candidate 2)

**Missing behavior:** The source does not address what happens if the same client submits the same order more than once (idempotency).

**Suggested candidate (draft; not normative until accepted):** `When a client resubmits an order that is identical to a previously accepted order, the order service shall [assign a new identifier | return the existing identifier] for the resubmitted order.`

**Reason:** Without an idempotency policy, a client retry on network failure may create duplicate orders. Whether this is in scope depends on the system's retry and deduplication model (Q4).

**Status:** optional consideration — ask whether idempotency or deduplication is in scope.

---

### SR-4 — Operational companion (relates to Candidate 2)

**Missing behavior:** The source does not state whether the identifier assignment and response are synchronous within the same call or may be deferred (e.g., asynchronous acknowledgement with identifier delivered later).

**Suggested candidate (draft):** `When a client submits a valid order, the order service shall return the assigned order identifier in the synchronous HTTP response to the submission request.`

**Reason:** The source says "return it in the submission response," which implies synchrony, but does not preclude an asynchronous model. If the interaction pattern is asynchronous, the candidate requirement text and its verification contract must be revised.

**Status:** needs clarification — confirmation of synchronous vs. asynchronous interaction model needed.

---

## Consistency Findings

| Finding | Kind | Affected Candidates | Detail |
|---|---|---|---|
| Assignment and return are coupled | Intentional split | C1, C2 | The source joined assignment and return with "and." They are split to allow independent verification. No conflict; both must be satisfied. |
| Validity definition used in both | Overlap (intentional) | C1, C2 | Both candidates share the same trigger condition (valid order submission). This is correct: the same event triggers both obligations. No conflict. |
| Invalid-order behavior undefined | Gap | C1, C2 | Neither candidate specifies the failure path. SR-1 addresses this gap. |
| Uniqueness of identifier unstated | Gap | C1 | The term "order identifier" implies uniqueness but does not state it normatively. SR-2 addresses this gap. |

No conflicts found among the candidate requirements. The primary finding is a gap in failure-path behavior (invalid order submission) and an unstated uniqueness constraint.

---

## Assumptions

1. **Synchronous interaction:** "Return it in the submission response" is interpreted as a synchronous response to the same submission call. If the order service operates asynchronously, Candidate 2 must be re-drafted.

2. **Positive quantity means integer > 0:** "Positive quantity" is interpreted as a whole-number integer greater than zero. Fractional quantities are assumed out of scope unless confirmed otherwise (Q5).

3. **Non-empty product identifier is sufficient for validity:** Validity of a product identifier is assumed to require only that it is present and non-empty — not that it resolves to a known product in a catalog. This assumption materially affects what constitutes a valid order (Q6).

4. **"Assign" implies persistence:** "Assign an order identifier" is interpreted to mean the identifier is durably associated with the order, not merely generated and discarded. If the service is stateless, this assumption does not hold (Q2).

---

## Rationale and Context

The source request establishes two behavioral obligations — identifier assignment and identifier return — and provides an inline validity definition. The primary purpose is to give clients a stable handle for the submitted order, enabling subsequent operations (e.g., order tracking, cancellation, fulfillment queries). Separating assignment (C1) from return (C2) makes both obligations independently verifiable and avoids conflating internal state management with the client-visible interface contract.

The inline validity definition ("at least one product identifier and a positive quantity for each product") functions as a precondition on the trigger, not a separate system behavior. It is captured in the implementability contracts of both candidates and expanded in SR-1 (failure-path companion) to cover the invalid-order case that the source does not address.

---

## Review Status

| Candidate | Status | Blocking items |
|---|---|---|
| Candidate 1 — Assign identifier | candidate | Q1 (uniqueness scope), Q2 (persistence), Q5 (quantity type), Q6 (product ID validity) should be resolved to reach `ready for review` |
| Candidate 2 — Return identifier | candidate | Q1 (format), Q3 (invalid-order response), synchrony assumption should be confirmed to reach `ready for review` |
| SR-1 — Invalid order rejection | needs clarification | Error response content and scope unconfirmed |
| SR-2 — Identifier uniqueness | needs clarification | Uniqueness scope unconfirmed (Q1) |
| SR-3 — Idempotency | optional consideration | Not yet in scope; Q4 must be answered |
| SR-4 — Synchronous response | needs clarification | Interaction model (sync vs. async) unconfirmed |

No candidate is marked `ready for review` because material uniqueness scope, failure-path behavior, and synchrony model remain unresolved. Both primary candidates are at `candidate` status: their implementability and verification contracts are complete modulo the open questions, and no contract item says "open" or "unresolved" without an accompanying clarifying question. The open questions are explicitly listed and would not prevent a single reasonable implementer from choosing a consistent behavior — but two reasonable implementers could differ on uniqueness scope (Q1) and invalid-order handling, which is why `ready for review` is withheld pending answers.
