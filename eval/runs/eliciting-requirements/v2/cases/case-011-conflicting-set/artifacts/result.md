# Requirements Elicitation Package

## Understanding

**Requested outcome:** Define how the order service responds when a client submits a duplicate order.

**Proposed system boundary:** The order service, at the boundary where it receives an order submission request from a client and returns a response.

**Source material:** Two requirements, R1 and R2, with the same scope, trigger, and system, but prescribing mutually exclusive responses. The source explicitly acknowledges the conflict and does not state which behavior has priority.

**What each artifact is:**

- **R1 [orders]** — A normative behavioral obligation: on duplicate submission, return the existing order's identifier.
- **R2 [orders]** — A normative behavioral obligation: on duplicate submission, create and return a new order with a new identifier.

**Conflict:** R1 and R2 cannot both be satisfied for the same duplicate submission event. They are not complementary obligations applying to different conditions; they apply to the same condition and mandate opposite responses. This is a direct conflict. No candidate can be marked `candidate` or `ready for review` until the conflict is resolved by a decision-maker.

---

## Clarifying Questions

**Q1 (Blocking — conflict resolution).** R1 and R2 prescribe opposite responses to the same event. Which behavior is intended?

- Option A: Return the identifier of the existing order (idempotent behavior — treat a duplicate as a repeat of the original request).
- Option B: Create a new order with a new identifier (non-idempotent behavior — treat every submission as a distinct order regardless of content).
- Option C: Choose different behavior based on a condition not yet captured in the source (e.g., client opt-in flag, order state, time window). If so, what is the differentiating condition?

This question must be answered before either candidate can advance.

**Q2 (Definition — "duplicate order").** What makes an order a duplicate? Specifically:

- Which fields are compared (e.g., client identifier, product identifier, quantity, requested delivery date, idempotency key)?
- Is a submitted order a duplicate only if it is identical in all fields, or only in selected key fields?
- Does the duplicate determination depend on the state of the existing order (e.g., only orders in an open or pending state count as duplicates)?

**Q3 (Scope — detection window).** Over what time window or scope is duplicate detection applied? Is it per client session, per calendar day, indefinitely, or until the existing order reaches a terminal state?

**Q4 (Response detail).** Whichever behavior is chosen:

- What identifier format does the service return (e.g., a UUID, an opaque token, an order number)?
- Is any additional information returned alongside the identifier (e.g., order status, creation timestamp)?

**Q5 (Failure paths).** If the service cannot determine whether an order is a duplicate (e.g., due to a data-store error or timeout), what should it return to the client?

---

## Candidate Requirements

### Candidate R1-revised

- **Status:** needs clarification
- **Template:** Event-driven
- **Requirement:** *(Drafting blocked pending resolution of Q1.)* When a client submits a duplicate order, the order service shall return the identifier of the existing order.
- **Source intent and rationale:** R1 intends idempotent handling of duplicate submissions — the client receives the same identifier it would have received had it received the original response, preventing duplicate order creation.
- **Implementability contract**
  - System boundary: The order service, at the order-submission API boundary.
  - Responsible system: Order service.
  - Actor: Client.
  - Affected interfaces: Order submission endpoint (inbound request); response returned to client (outbound).
  - Inputs and data definitions: An order submission request whose content matches an existing order. "Duplicate" is undefined — see Q2 and Q3.
  - Preconditions, state, feature scope, and trigger: A matching order already exists (definition of "matching" unresolved — Q2, Q3). Trigger: client submits an order submission request.
  - Observable normal, alternate, and failure response: Return the identifier of the existing matching order. Failure path: unresolved — Q5.
  - Quantities, units, limits, timing, ordering, and tolerances: Identifier format and response content unresolved — Q4.
  - Domain terms, dependencies, and unresolved decisions: "Duplicate order" undefined (Q2, Q3). Conflict with R2 unresolved (Q1). Failure behavior unresolved (Q5).
- **Verification contract**
  - Setup, state, input data, and stimulus: An existing order is present in the service. A client submits a second order whose fields match the duplicate criteria.
  - Observable evidence or oracle: The HTTP response (or equivalent) returned by the order service to the client.
  - Expected result: The response contains the identifier of the pre-existing order, not a new identifier.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — response identifier matches the identifier of the pre-existing order. Fail — a new identifier is returned, or no identifier is returned. Specific identifier format and response schema unresolved (Q4).
  - Boundary, negative, and failure cases: First submission of an order (non-duplicate) — out of scope for this candidate. Service error during duplicate detection — unresolved (Q5).
  - Required data, instrumentation, assessment, or external evidence: Access to the identifier assigned to the pre-existing order to compare against the response.
  - Observable at the named boundary alone: yes
  - If no, internal evidence needed and rationale: N/A

---

### Candidate R2-revised

- **Status:** needs clarification
- **Template:** Event-driven
- **Requirement:** *(Drafting blocked pending resolution of Q1.)* When a client submits a duplicate order, the order service shall create a new order and shall return the identifier of the new order.
- **Source intent and rationale:** R2 intends non-idempotent handling — each submission, regardless of matching content, produces a distinct order record. This may be appropriate when each submission represents a distinct business intent (e.g., reordering the same product).
- **Note on splitting:** The source sentence joins two independently testable obligations — creation of the new order, and return of the new identifier. These have been kept together here because they are causally bound (the returned identifier must be the one created), but a host may choose to split them.
- **Implementability contract**
  - System boundary: The order service, at the order-submission API boundary.
  - Responsible system: Order service.
  - Actor: Client.
  - Affected interfaces: Order submission endpoint (inbound); response returned to client (outbound); order data store (side effect — creation of new record).
  - Inputs and data definitions: An order submission request whose content matches an existing order. "Duplicate" is undefined — see Q2 and Q3.
  - Preconditions, state, feature scope, and trigger: A matching order already exists (definition of "matching" unresolved — Q2, Q3). Trigger: client submits an order submission request.
  - Observable normal, alternate, and failure response: Create a new order record with a new identifier; return that identifier to the client. Failure path: unresolved — Q5.
  - Quantities, units, limits, timing, ordering, and tolerances: Identifier format and response content unresolved — Q4.
  - Domain terms, dependencies, and unresolved decisions: "Duplicate order" undefined (Q2, Q3). Conflict with R1 unresolved (Q1). Failure behavior unresolved (Q5).
- **Verification contract**
  - Setup, state, input data, and stimulus: An existing order is present in the service. A client submits a second order whose fields match the duplicate criteria.
  - Observable evidence or oracle: The HTTP response (or equivalent) returned by the order service; the set of orders present in the service after the submission.
  - Expected result: The response contains a new identifier distinct from the pre-existing order's identifier. A new order record exists in the service with that identifier.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — response identifier is new and distinct; a new order record is observable. Fail — the response identifier matches the pre-existing order, or no new order record is created.
  - Boundary, negative, and failure cases: First submission of an order (non-duplicate) — out of scope for this candidate. Service error during creation — unresolved (Q5).
  - Required data, instrumentation, assessment, or external evidence: Ability to inspect the set of orders in the service after submission to confirm a new record was created.
  - Observable at the named boundary alone: no
  - If no, internal evidence needed and rationale: Confirming that a new order was created (not merely that a new identifier was returned) requires observing the order store or a subsequent retrieval API call. The response alone cannot distinguish "new identifier for a new order" from "new identifier assigned to a pre-existing record."

---

## Suggested Supporting Requirements

1. **Failure-path companion (relates to both R1-revised and R2-revised).**
   Missing behavior: what the order service returns to the client when it cannot determine whether the submitted order is a duplicate (e.g., data-store unavailability, timeout, lock contention).
   Reason: without this, two implementers may choose different observable responses (error code vs. success with degraded behavior vs. retry instruction), making verification impossible.
   Suggested draft: `If the order service cannot complete duplicate detection for a submitted order, then the order service shall return an error response to the client that distinguishes the failure from a successful submission.`

2. **Required companion — definition of "duplicate order" (relates to both candidates).**
   Missing behavior: the fields and scope used to classify a submission as a duplicate.
   Reason: without a definition, no implementer or verifier can determine which submissions trigger either R1-revised or R2-revised.
   Suggested draft: a data definition requirement specifying the fields compared and the time or state scope (to be drafted once Q2 and Q3 are answered).

3. **Boundary companion — first submission (non-duplicate) (relates to both candidates).**
   Missing behavior: the response when a submitted order does not match any existing order.
   Reason: the source addresses only the duplicate case; a normal-path requirement is needed for a complete specification of the order-submission operation.
   Suggested draft: `When a client submits an order that is not a duplicate, the order service shall create a new order and shall return the identifier of the new order.`

4. **Optional consideration — concurrent duplicate submissions.**
   If two clients submit identical orders simultaneously, race conditions may cause both to appear non-duplicate at detection time, resulting in two new orders regardless of which behavior (R1 or R2) is chosen. Ask whether the service must guarantee exactly-once creation under concurrency.

---

## Consistency Findings

### Finding 1 — Direct conflict (R1, R2)

| Attribute | R1 | R2 |
|---|---|---|
| System | Order service | Order service |
| Trigger | Client submits a duplicate order | Client submits a duplicate order |
| Response | Return identifier of existing order | Create new order; return new identifier |

**Consequence:** Both requirements apply to identical conditions and mandate opposite observable responses. They cannot both be satisfied for any single duplicate submission event. This is the most severe finding class: a direct conflict.

**Resolution options:**
- **Narrow:** Introduce a differentiating condition (e.g., a client-supplied idempotency key present → apply R1; absent → apply R2). The differentiating condition must be added as an explicit precondition or trigger qualifier to each candidate.
- **Choose one:** Discard or supersede the other. Requires a decision from the requirement owner.
- **Precedence rule:** State explicitly which requirement takes priority and under what authority. A general precedence rule without a differentiating condition still leaves one requirement permanently inactive (vacuous).
- **Ask the source:** The source explicitly states that priority is not defined. A decision from the requirement owner is required.

**Status:** Blocking. Neither candidate can advance past `needs clarification` until this finding is resolved.

---

## Assumptions

1. **System boundary is the external API.** The order service is a distinct system with a defined client-facing API. The normative obligations apply at that boundary, not at internal service components.
2. **"Identifier" is a single, externally visible value.** Both R1 and R2 refer to an identifier that the client can observe in the response. Whether that identifier is a UUID, order number, or another format is unresolved (Q4).
3. **"Duplicate order" is a meaningful, deterministic classification.** The service has access to enough information to classify an incoming submission as a duplicate. The criteria are currently undefined (Q2, Q3).
4. **R1 and R2 share the same scope.** Both requirements apply to the same order service and the same order-submission operation. They are not scoped to different environments, deployments, or feature flags unless a differentiating condition is introduced to resolve Finding 1.

---

## Rationale and Context

The source explicitly presents R1 and R2 as conflicting and acknowledges that priority is not stated. This is an unresolved design decision, not a drafting imprecision. The two behaviors represent genuinely different product policies:

- **R1 (idempotent)** supports at-least-once delivery patterns: clients can safely retry a submission and always receive the same order identifier. This reduces duplicate order creation but requires the service to track and match prior submissions.
- **R2 (non-idempotent)** treats every submission as a distinct business act. This may be correct when reordering the same items is a legitimate and common workflow, but it shifts the burden of deduplication to the client.

Neither policy is universally correct. The decision belongs to the requirement owner and must be made before implementation begins.

---

## Review Status

| Candidate | Status | Blocking reason |
|---|---|---|
| R1-revised | needs clarification | Direct conflict with R2 (Finding 1, Q1); "duplicate order" undefined (Q2, Q3); failure path undefined (Q5) |
| R2-revised | needs clarification | Direct conflict with R1 (Finding 1, Q1); "duplicate order" undefined (Q2, Q3); failure path undefined (Q5) |

**No candidate is ready for review.** The conflict identified in Finding 1 is the primary blocker. Questions Q2, Q3, Q4, and Q5 are secondary blockers that must also be resolved before either candidate can reach `ready for review`.
