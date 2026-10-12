# Requirements Elicitation Package

**Host metadata**
- ID: `BILL-REQ-42`
- Tags: `[billing, external-api]`
- Trace: `interface://billing/v2/charge`

---

## Understanding

**Requested outcome:** The Billing API shall evaluate an incoming charge request and return exactly one of three decision values (approved, declined, or indeterminate) as a JSON response body, within a 500 ms deadline, for all valid requests under specified concurrency.

**Proposed system boundary:** The Billing API service boundary — specifically the HTTP request/response layer exposed at `interface://billing/v2/charge`. The timing scope is precisely delimited: clock starts when the complete request body has been received and parsed, and ends when the final response byte is emitted at the service boundary.

**Source material sorted by kind:**

| Kind | Source text |
|---|---|
| Behavioral obligation | Return a JSON response body containing only one field named `decision` |
| Behavioral obligation | Value of `decision` shall be `approved`, `declined`, or `indeterminate` |
| Behavioral obligation | `indeterminate` is terminal for the request |
| Behavioral obligation | Response within 500 ms for every valid request under ≤100 concurrent valid requests |
| Constraint (scope) | Valid request = positive decimal amount ≤ 10,000 in the request currency, with currency in {USD, EUR, GBP} |
| Constraint (exclusion) | Invalid amounts, unsupported currencies, and behavior after a missed deadline are out of scope |
| Timing definition | Clock starts: complete request body received and parsed; ends: final response byte emitted at service boundary |
| Design hint | "only one field" in the JSON body — externally observable contract |
| Domain term needing confirmation | "terminal" — interpreted as: no further action by the Billing API on this request; the decision is final and not retried or escalated within this API |
| Assumption | "positive decimal amount" means strictly greater than zero and expressed as a decimal number |

**Independent obligations identified (requiring separate candidates):**

1. The structure and content of the decision response body (what the field is, what values are allowed).
2. The terminal semantics of `indeterminate` (a behavioral constraint on what the API does after returning that value).
3. The latency obligation (500 ms deadline under ≤100 concurrent requests).

---

## Clarifying Questions

**Q1 — Currency field in the request body:**
The source describes the request currency as part of the valid-request scope, but does not name the field. Is the currency communicated as a request body field (e.g., `"currency": "USD"`), a path parameter, a query parameter, or a header? This affects what a verifier observes at the boundary.

**Q2 — Meaning of "terminal":**
The source states that `indeterminate` is terminal for the request. Does this mean:
(a) The Billing API closes or discards the connection after emitting the response and will not process a retry of the same request, or
(b) The caller must not retry the same logical charge on the same connection, but the API has no enforcement obligation, or
(c) Something else?
This affects whether a companion requirement for retry/idempotency is in scope.

**Q3 — Exactly-one-field constraint:**
The source says the JSON body contains "only one field." Does this mean:
(a) The response body is a JSON object with exactly the key `decision` and no other keys at any time, or
(b) Additional metadata keys may appear in future versions, and this is a minimum-field guarantee?
The word "only" as written reads as a strict constraint, but verifiers need to know whether extra keys cause a test failure.

**Q4 — Amount type precision:**
"Positive decimal amount" — is there a minimum positive value (e.g., 0.01), a maximum number of decimal places, or a representation constraint (e.g., IEEE 754 double vs. arbitrary-precision string)? Two implementers could differ on whether `0.001` or `9999.999999` is a valid amount.

**Q5 — Concurrency scope for latency:**
The 500 ms deadline applies "under a load of up to 100 concurrent valid requests." Does this mean:
(a) The deadline applies at any point up to and including exactly 100 concurrent requests in flight simultaneously, or
(b) The deadline applies to the steady-state throughput of 100 requests per second (arrival rate)?
These are different workload models and affect how verification is performed.

**Q6 — Latency percentile:**
The source says "within a maximum of 500 ms for every valid request." Does "every" mean 100th percentile (no single response may exceed 500 ms), or is a percentile (e.g., p99) intended? "Every" reads as an absolute ceiling, but this is operationally significant.

---

## Candidate Requirements

### Candidate 1 — Decision response structure

- **Status:** `ready for review`
- **Template:** Event-driven
- **Requirement:** When the Billing API receives a valid charge request at the `interface://billing/v2/charge` service boundary, the Billing API shall return a JSON response body containing exactly one field, named `decision`, with a value of `approved`, `declined`, or `indeterminate`.
- **Source intent and rationale:** The caller must be able to parse the outcome of a charge attempt from a single, stable field. Restricting the body to one field eliminates ambiguity about which field is authoritative. *(See Q3 regarding whether "only one field" is a strict exclusion of additional keys.)*
- **Implementability contract**
  - System boundary: `interface://billing/v2/charge` HTTP response boundary
  - Responsible system: Billing API
  - Actor: Any client submitting a valid charge request
  - Affected interfaces: HTTP response body at `interface://billing/v2/charge`
  - Inputs and data definitions: Valid charge request — positive decimal amount ≤ 10,000 in the request currency; currency in {USD, EUR, GBP}. *(Q1: currency field location unresolved but does not block this candidate. Q4: amount precision unresolved but scoped to the valid-request definition, which is excluded from this candidate.)*
  - Preconditions, state, feature scope, and trigger: Receipt of a valid charge request at the service boundary.
  - Observable normal, alternate, and failure response: JSON object with exactly one key `decision` set to one of the three enumerated string values. *(Q3: whether extra keys cause failure is open but does not block the normative content of this candidate.)*
  - Quantities, units, limits, timing, ordering, and tolerances: Exactly one field; exactly three permissible values for that field (enumerated).
  - Domain terms, dependencies, and unresolved decisions: "Valid charge request" defined in source. "Exactly one field" treated as a strict exclusion of additional keys per source text; Q3 raises whether this is intended as an invariant or a minimum.
- **Verification contract**
  - Setup, state, input data, and stimulus: Submit a well-formed charge request with a positive decimal amount ≤ 10,000 and a currency of USD, EUR, or GBP to `interface://billing/v2/charge`.
  - Observable evidence or oracle: The HTTP response body, inspected at the service boundary.
  - Expected result: A JSON object with exactly the key `decision` and a value that is one of the strings `approved`, `declined`, or `indeterminate`.
  - Pass/fail criteria: Pass — body is a JSON object, key count is exactly 1, key name is `decision`, value is one of the three permitted strings. Fail — any additional key present, key absent, key misspelled, value outside the three-string enum, or body is not valid JSON.
  - Boundary, negative, and failure cases: (a) Value is exactly `approved`; (b) value is exactly `declined`; (c) value is exactly `indeterminate`; (d) no fourth value is accepted.
  - Required data, instrumentation, assessment, or external evidence: HTTP client able to read and parse the response body at the network boundary.
  - Observable at the named boundary alone: **Yes**
  - If no, internal evidence needed and rationale: N/A

---

### Candidate 2 — Indeterminate is terminal

- **Status:** `needs clarification`
- **Template:** Unresolved (see Q2)
- **Requirement:** Drafting blocked pending resolution of Q2. Provisional sentence pending clarification:

  > *When the Billing API returns `indeterminate` for a charge request at the `interface://billing/v2/charge` service boundary, the Billing API shall [observable terminal action — to be determined].*

- **Source intent and rationale:** The source marks `indeterminate` as terminal for the request, distinguishing it from `declined` (which may be retryable) and `approved`. The observable behavior the API must produce after returning `indeterminate` is not yet defined: does the API close the connection, reject a retry, log the decision, or simply make no further commitment? Until Q2 is answered, two implementers could produce different observable behavior.
- **Implementability contract**
  - System boundary: `interface://billing/v2/charge`
  - Responsible system: Billing API
  - Actor: Client that received an `indeterminate` response
  - Affected interfaces: `interface://billing/v2/charge` — behavior after response emission
  - Inputs and data definitions: The previously emitted `indeterminate` decision
  - Preconditions, state, feature scope, and trigger: The Billing API has just emitted a response with `decision: indeterminate`
  - Observable normal, alternate, and failure response: **Open — depends on Q2**
  - Quantities, units, limits, timing, ordering, and tolerances: N/A until Q2 resolved
  - Domain terms, dependencies, and unresolved decisions: **"Terminal" is undefined as an observable API behavior** — Q2 is blocking.
- **Verification contract**
  - Setup, state, input data, and stimulus: **Open — depends on Q2**
  - Observable evidence or oracle: **Open**
  - Expected result: **Open**
  - Pass/fail criteria: **Open**
  - Boundary, negative, and failure cases: **Open**
  - Required data, instrumentation, assessment, or external evidence: **Open**
  - Observable at the named boundary alone: **Open**
  - If no, internal evidence needed and rationale: **Open**

---

### Candidate 3 — Latency deadline

- **Status:** `needs clarification`
- **Template:** Complex (While ≤100 concurrent valid requests are in flight, when a valid charge request is received…) — *blocked on Q5 and Q6*
- **Requirement:** Provisional sentence:

  > *While the Billing API is processing up to 100 concurrent valid charge requests, when a valid charge request is received at the `interface://billing/v2/charge` service boundary, the Billing API shall emit the final response byte within 500 ms of completing receipt and parsing of the complete request body.*

- **Source intent and rationale:** The source establishes a strict 500 ms wall-clock ceiling per valid request, with a defined clock start (complete body received and parsed) and a defined clock end (final response byte at the service boundary). The workload scope (≤100 concurrent) must be confirmed (Q5), and the percentile interpretation of "every" must be confirmed (Q6) before verification can be specified without ambiguity.
- **Implementability contract**
  - System boundary: `interface://billing/v2/charge`
  - Responsible system: Billing API
  - Actor: Any client submitting a valid charge request
  - Affected interfaces: `interface://billing/v2/charge` response emission point
  - Inputs and data definitions: Valid charge request as defined in Candidate 1
  - Preconditions, state, feature scope, and trigger: While ≤100 concurrent valid charge requests are in flight; when a valid charge request body has been completely received and parsed.
  - Observable normal, alternate, and failure response: Final response byte emitted at the service boundary. *(What happens if the deadline is missed is explicitly out of scope per the source.)*
  - Quantities, units, limits, timing, ordering, and tolerances: 500 ms maximum. Clock starts: complete request body received and parsed. Clock ends: final response byte emitted at service boundary. Concurrency scope: ≤100 concurrent valid requests. **Q5: concurrency vs. arrival rate open. Q6: "every" (100th percentile) vs. a lower percentile open.**
  - Domain terms, dependencies, and unresolved decisions: "Concurrent valid requests" — concurrency model unclear (Q5). "Every valid request" — percentile unclear (Q6).
- **Verification contract**
  - Setup, state, input data, and stimulus: Generate up to 100 simultaneous valid charge requests directed at `interface://billing/v2/charge`. Measure elapsed time from the event "complete request body received and parsed by the Billing API" to the event "final response byte emitted at the service boundary" for each request.
  - Observable evidence or oracle: Per-request elapsed time measured at the service boundary, correlated with the request-body-parsed timestamp.
  - Expected result: Every measured elapsed time is ≤ 500 ms.
  - Pass/fail criteria: Pass — all elapsed times ≤ 500 ms under the specified load. Fail — any single response exceeds 500 ms. *(Q6: if a percentile threshold is intended instead of an absolute ceiling, the pass/fail rule changes.)*
  - Boundary, negative, and failure cases: (a) Exactly 1 concurrent request; (b) exactly 100 concurrent requests; (c) the request that arrives when 99 are already in flight (boundary). *(Q5: if "up to 100" means arrival rate, boundary cases differ.)*
  - Required data, instrumentation, assessment, or external evidence: Instrumentation at the service boundary to capture (i) the timestamp when the complete request body is parsed by the Billing API and (ii) the timestamp when the final byte of the response is emitted. Both timestamps must be observable at or derivable from the service boundary.
  - Observable at the named boundary alone: **Partially.** The response-byte timestamp is observable at the network boundary. The request-body-parsed timestamp may require instrumentation inside the Billing API or at the transport layer, because "parsed" is an internal event. External instrumentation that captures the moment the body is fully received at the boundary is a proxy, but may not capture parsing latency.
  - If no, internal evidence needed and rationale: The "clock starts when the complete request body has been received **and parsed**" definition means parsing time is inside the clock. Verifiers need either (a) an internal instrumentation hook that emits the parsed timestamp, or (b) agreement that the receipt timestamp (last byte of request body received at the boundary) is an acceptable proxy for "received and parsed."

---

## Suggested Supporting Requirements

**SR-1 — Currency field location (Required companion)**
*Relates to Candidate 1.* The source defines valid currency as a member of {USD, EUR, GBP} but does not specify where in the request the currency is conveyed. A companion requirement should define the field name, location (request body, path parameter, header), and required format. Without it, two implementers could parse currency differently, making the valid/invalid boundary ambiguous at the service level.

**SR-2 — Amount field definition (Required companion)**
*Relates to Candidate 1.* The source specifies "positive decimal amount no greater than 10,000" but does not name the request field, its data type (string vs. number), number of decimal places, or minimum granularity. A companion requirement should specify the field name, wire format, precision, and minimum accepted value. *(See Q4.)*

**SR-3 — HTTP response status code (Interface companion)**
*Relates to Candidate 1.* The source specifies the JSON body but is silent about the HTTP status code returned with a valid charge decision. If the status code is part of the interface contract, it should be specified. An unspecified status code could cause callers to misinterpret a 200 vs. 202 vs. other code as success or failure independent of the `decision` value.

**SR-4 — Content-Type response header (Interface companion)**
*Relates to Candidate 1.* The source specifies a JSON response body but does not specify the `Content-Type` header value. A companion requirement should state that the Billing API shall return `Content-Type: application/json` (or the applicable media type) with every decision response, so that clients can safely parse the body.

**SR-5 — Retry and idempotency behavior for `declined` (Boundary companion)**
*Relates to Candidate 2.* The source distinguishes `indeterminate` as terminal but is silent about whether `declined` is retryable or whether the charge endpoint is idempotent. If a client submits the same charge request twice, the API's behavior should be specified to avoid double-charging or inconsistent decisions.

**SR-6 — Behavior after a missed deadline (Failure-path companion)**
*Relates to Candidate 3.* The source explicitly places missed-deadline behavior out of scope and notes it is "specified separately." This companion is a placeholder to ensure that the separately specified behavior is linked to this requirement set and that the gap is tracked.

**SR-7 — Clock-start observability instrumentation (Operational companion)**
*Relates to Candidate 3.* The timing clause requires the clock to start "when the complete request body has been received and parsed." This internal event is not directly observable at the network boundary. A companion requirement or operational obligation should specify that the Billing API shall expose a machine-readable timestamp (e.g., a response header, log entry, or telemetry event) marking when the request body was fully parsed, so that the 500 ms obligation can be independently verified.

---

## Consistency Findings

**Finding 1 — Gap: valid-request definition is shared but not owned**
Candidates 1 and 3 both rely on the definition of "valid charge request" (positive decimal ≤ 10,000, currency in {USD, EUR, GBP}). The source states that invalid amounts and unsupported currencies are "specified separately." There is a gap: the valid-request boundary is referenced here but defined only by inline prose. If the separately specified requirements define a different boundary, Candidates 1 and 3 could be satisfied under different inputs. The valid-request definition should be a separately identified, traceable requirement that Candidates 1 and 3 reference explicitly.

**Finding 2 — Precedence question: indeterminate vs. declined semantics**
The source distinguishes `indeterminate` as terminal but does not say whether `declined` is also terminal, retryable, or subject to any other constraint. Candidate 2 covers only `indeterminate`. If a future requirement defines retry behavior for `declined`, it must be checked for conflict with the terminal semantics of `indeterminate`. The precedence or difference between these two non-`approved` values should be made explicit.

**Finding 3 — Potential overlap: Candidate 1 and Candidate 3 share the same trigger**
Both Candidate 1 (response body structure) and Candidate 3 (latency) are triggered by receiving a valid charge request. They are correctly separated as independent obligations (structure vs. timing). No conflict exists, but a verifier should confirm that both are evaluated against the same valid-request population.

**Finding 4 — Gap: no content-type or HTTP status specified**
The source is silent on HTTP status code and `Content-Type` header. If either is later added by a separate requirement, it must be checked against Candidate 1's response body contract. *(See SR-3 and SR-4.)*

---

## Assumptions

1. **"Positive decimal amount"** is interpreted as a real number strictly greater than zero and representable as a decimal (not binary-only). The maximum precision is not assumed; Q4 asks about it.
2. **"Terminal"** is interpreted colloquially as "the Billing API makes no further attempt or escalation on this request after returning `indeterminate`." The observable API behavior that enforces or signals this is unknown; Q2 asks about it.
3. **"Concurrent valid requests"** is interpreted as simultaneously in-flight requests — i.e., requests whose response has not yet been emitted while another request is being processed. Q5 asks whether this is correct or whether an arrival-rate model is intended.
4. **"Every valid request"** in the latency clause is taken at face value as a 100th-percentile (absolute ceiling) constraint. Q6 asks whether a lower percentile was intended.
5. The phrase "only one field" in the JSON body is treated as a strict constraint — the body shall contain exactly the key `decision` and no other keys. Q3 asks whether this is intended as an invariant.
6. The HTTP method and full URL path for the charge endpoint are implied by the trace `interface://billing/v2/charge` but are not specified in this requirement. They are assumed to be defined by the interface contract and not reproduced here.
7. The separately specified requirements for invalid amounts, unsupported currencies, and missed-deadline behavior exist and are consistent with this package. They are not assumed to be in scope here.

---

## Rationale and Context

The source requirement is unusually precise for a natural-language obligation. It explicitly defines:
- The valid-input scope (amount range and currency enumeration)
- The response shape (single field, enumerated values)
- The timing model (clock start and end events)
- The concurrency scope for the latency obligation
- The out-of-scope behaviors

This precision reduces elicitation risk substantially. The main open items are:
- Where currency is communicated in the request (SR-1, Q1)
- The observable meaning of "terminal" for `indeterminate` (Candidate 2, Q2)
- Whether the latency constraint is an absolute ceiling or a percentile (Candidate 3, Q6)
- The concurrency model for the latency workload (Q5)

The response-body structure (Candidate 1) is close to `ready for review` subject only to the Q3 clarification about whether extra keys in the response are a hard failure.

---

## Review Status

| Candidate | Status | Blocking item |
|---|---|---|
| Candidate 1 — Decision response structure | `ready for review` | None blocking; Q3 is a refinement, not a gate |
| Candidate 2 — Indeterminate is terminal | `needs clarification` | Q2: observable meaning of "terminal" undefined |
| Candidate 3 — Latency deadline | `needs clarification` | Q5: concurrency model; Q6: percentile interpretation; timing instrumentation observability |

**No candidates are approved.** This package is the output of elicitation; approval is a host responsibility.
