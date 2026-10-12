# Requirements Elicitation Package

## Understanding

**Requested outcome:** When a customer submits a payment that the payment
service judges to be valid, the payment service must return an authorization
result to the caller within 2 seconds.

**Proposed system boundary:** The payment service — specifically its inbound
payment submission interface (the boundary at which a customer or calling
system delivers a payment request) and its outbound response (the boundary at
which the payment service returns an authorization result).

**Source classification:**

| Source fragment | Kind |
|---|---|
| "When a customer submits a valid payment" | Trigger / behavioral obligation scope |
| "the payment service shall return an authorization result" | Behavioral obligation |
| "within 2 seconds" | Timing constraint (normative) |

The source uses `shall`, which is treated as a normative obligation. The
phrasing maps directly to the Event-driven EARS pattern. Several terms are
undefined or ambiguous and block full readiness.

---

## Clarifying Questions

**Q1 — Definition of "valid payment":** What conditions must a payment request
satisfy to be considered valid? For example: presence of required fields,
passing format checks, passing fraud screening, successful authentication, or
some combination. This matters because the trigger fires only on a valid
payment; the behavior for an invalid submission is currently unspecified.

**Q2 — Definition of "authorization result":** What is the observable content
or structure of an authorization result? At minimum: does it include an
authorization code, an approval/decline indicator, or both? Is a partial
result (e.g., a pending status) considered a valid authorization result for
this requirement?

**Q3 — Measurement point for the 2-second deadline:** Is the 2 seconds
measured from the moment the payment service receives the complete payment
request, to the moment the last byte of the authorization result is sent at
the service boundary? Or is it measured end-to-end from the customer's
perspective? The measurement point is required to make the timing contract
verifiable.

**Q4 — Behavior when the deadline is exceeded:** If the payment service
cannot return an authorization result within 2 seconds, what must it do?
Options include: return a timeout error response, cancel the authorization
attempt, or continue processing and return a late result. The source is silent
on this.

**Q5 — Behavior on invalid payment submission:** What shall the payment
service return (or do) when a customer submits a payment that is not valid?
This is a gap in the source — the 2-second obligation applies only to valid
submissions.

**Q6 — Concurrency and workload scope:** Does the 2-second deadline apply at
all load levels, or only up to a defined concurrent request rate or transaction
volume? An unconditional latency requirement may be infeasible under
pathological load; knowing the intended workload scope makes the timing
constraint verifiable.

---

## Candidate Requirements

### Candidate 1

- **Status:** needs clarification
- **Template:** Event-driven
- **Requirement:** When a customer submits a valid payment, the payment service
  shall return an authorization result within 2 seconds.

- **Source intent and rationale:** The source requests a bounded-latency
  guarantee on the normal path of a payment authorization. The 2-second limit
  implies a real-time or near-real-time user experience goal for payment
  confirmation.

- **Implementability contract**
  - System boundary: The payment service inbound submission interface and
    outbound response interface.
  - Responsible system: The payment service.
  - Actor: Customer (or a calling system acting on behalf of a customer).
  - Affected interfaces: Payment submission interface (inbound); authorization
    result interface (outbound).
  - Inputs and data definitions: "Valid payment" — **unresolved** (see Q1).
    "Authorization result" — **unresolved** (see Q2).
  - Preconditions, state, feature scope, and trigger: Trigger is the discrete
    event of a customer submitting a payment that satisfies the validity
    criteria. No precondition beyond validity is stated.
  - Observable normal, alternate, and failure response: Normal response is an
    authorization result returned within 2 seconds. Behavior when the deadline
    is exceeded is **unresolved** (see Q4). Behavior on invalid submission is
    **unresolved** (see Q5).
  - Quantities, units, limits, timing, ordering, and tolerances: 2 seconds.
    Measurement point is **unresolved** (see Q3). Applicable workload scope is
    **unresolved** (see Q6).
  - Domain terms, dependencies, and unresolved decisions: "Valid payment,"
    "authorization result," measurement point, timeout behavior, and invalid
    submission behavior are all unresolved.

- **Verification contract**
  - Setup, state, input data, and stimulus: A well-formed, valid payment
    request is submitted to the payment service at its inbound interface.
    Definition of "valid" must be resolved before the test can be specified.
  - Observable evidence or oracle: The payment service emits an authorization
    result at its outbound response interface. The elapsed time between
    receiving the request and emitting the result is recorded.
  - Expected result: An authorization result is returned; the elapsed time does
    not exceed 2 seconds.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass if
    an authorization result is observable at the service boundary within 2
    seconds of receiving a valid payment submission. The exact measurement
    points are **unresolved** (see Q3).
  - Boundary, negative, and failure cases: Behavior at exactly 2 seconds;
    behavior when the deadline is exceeded; behavior on invalid input — all
    **unresolved** (see Q4, Q5).
  - Required data, instrumentation, assessment, or external evidence:
    Timestamps at the inbound and outbound service boundary; definition of
    authorization result structure; agreed workload parameters.
  - Observable at the named boundary alone: yes (timing and response content
    are both observable at the payment service interface, once terms are
    defined).
  - If no, internal evidence needed and rationale: N/A.

---

## Suggested Supporting Requirements

**SR-1 — Failure-path companion (timeout behavior):**
Related to Candidate 1 (timeout on valid payment submission).
Missing behavior: What the payment service shall return if the 2-second
deadline is exceeded before an authorization result is available.
Suggested draft (provisional, not normative until accepted): "If the payment
service cannot return an authorization result within 2 seconds of receiving a
valid payment submission, then the payment service shall return a timeout
error response to the caller."
Reason: Without this, a verifier cannot distinguish a late response from no
response, and an implementer cannot choose a safe default.

**SR-2 — Failure-path companion (invalid payment response):**
Related to Candidate 1 (behavior on invalid payment submission).
Missing behavior: What the payment service shall return when the submitted
payment is not valid.
Suggested draft (provisional): "When a customer submits a payment that does
not satisfy the validity criteria, the payment service shall return a rejection
response that identifies the reason for rejection."
Reason: The source scopes the 2-second obligation to valid payments only;
behavior on invalid submissions is a material gap.

**SR-3 — Boundary companion (concurrent workload scope):**
Related to Candidate 1 (workload scope of the timing constraint).
Missing behavior: The maximum concurrent load or transaction volume under which
the 2-second guarantee applies.
Suggested draft (provisional): "The payment service shall return an
authorization result within 2 seconds for up to [N] concurrent valid payment
submissions." ([N] is a placeholder requiring a product or capacity decision.)
Reason: An unconditional latency guarantee without a workload bound is not
meaningfully verifiable and may not be implementable.

**SR-4 — Operational companion (auditability):**
Related to Candidate 1 (observability of authorization results).
Missing behavior: Whether authorization results and their associated timing
must be recorded for audit or dispute resolution.
Suggested draft (provisional): "The payment service shall record the timestamp
of each received valid payment submission and the timestamp of each returned
authorization result."
Reason: A 2-second latency requirement implies the ability to measure and
audit response time in production; without retained timestamps this may not be
operationally verifiable.

---

## Consistency Findings

**Finding 1 — Gap: No failure path for deadline breach**
Candidate 1 defines the normal-path obligation (return within 2 seconds) but
is silent on what occurs if the deadline is not met. This leaves both
implementers and verifiers with no normative guidance on timeout behavior.
Resolution options: (a) add SR-1 as an accepted requirement; (b) confirm that
late results are acceptable and remove the deadline; (c) confirm that the
calling system must handle absence of a response within 2 seconds without a
service-side obligation.

**Finding 2 — Gap: Invalid payment behavior undefined**
The trigger condition ("valid payment") implies a complementary case
(invalid/unrecognized payment), for which no requirement exists. This is a
behavioral gap, not a conflict.
Resolution options: (a) accept SR-2; (b) explicitly state that invalid
submissions are out of scope for this requirement set.

**Finding 3 — Potential vacuousness without workload scope**
An unconditional 2-second latency requirement may be infeasible or untestable
without a defined workload ceiling, making it potentially vacuous in
production. This is not a conflict among candidates but a concern about the
single candidate's verifiability.
Resolution options: (a) accept SR-3 to bound the workload; (b) confirm the
requirement is intended to apply unconditionally and document the implementer's
obligation to provision accordingly.

No conflicts, duplicates, or precedence questions exist among the current
candidates because there is only one candidate. The above findings concern
completeness and verifiability.

---

## Assumptions

**A1** — "Customer" is interpreted as any principal (human user or calling
system) that submits a payment to the payment service's submission interface.
If customers and calling systems have distinct behavioral obligations, this
assumption must be revisited.

**A2** — "Submit" is interpreted as the discrete event of delivering a
complete payment request to the payment service at its inbound interface. If
submission is a multi-step process (e.g., initiate, confirm, finalize), the
trigger point must be specified.

**A3** — The 2-second limit is interpreted as a maximum elapsed time, not a
target or a mean. No statistical tolerance (e.g., 99th percentile) is assumed.
If the intent is a percentile-based constraint, the requirement must be
revised.

**A4** — "Return" is interpreted as the payment service emitting the
authorization result at its outbound interface to the caller. Network
transmission time beyond the service boundary is not included in the 2-second
obligation unless Q3 resolves otherwise.

---

## Rationale and Context

The source sentence establishes a real-time latency guarantee for the
authorization path of a payment service. The 2-second bound is a common
industry target for interactive payment flows, suggesting a user-experience
rationale (customers should not wait more than 2 seconds for a payment
decision). However, the source does not state this rationale explicitly.

The requirement is structurally well-formed — it names a system, a trigger,
and a response with a timing constraint — but several terms are undefined and
failure paths are absent. These gaps do not change the intent of the source
request; they are information required to make the requirement implementable
and verifiable.

---

## Review Status

- **Candidate 1:** `needs clarification`
  Blocking items: Q1 (definition of "valid payment"), Q2 (definition of
  "authorization result"), Q3 (measurement point for 2-second deadline), Q4
  (behavior when deadline is exceeded), Q5 (behavior on invalid payment), Q6
  (workload scope). The verification contract cannot specify setup, oracle, or
  pass/fail criteria until these are resolved.

- **SR-1 through SR-4:** Suggested supporting requirements, not yet accepted.
  Status will be assigned after the user confirms scope.

No candidates are `ready for review`. All candidates require resolution of the
clarifying questions before the readiness gate can be re-evaluated.
