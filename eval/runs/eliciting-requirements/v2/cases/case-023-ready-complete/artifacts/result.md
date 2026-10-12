# Requirements Elicitation Package

## Understanding

The request covers a single unauthenticated HTTP endpoint exposed by the Health API. Any client may call it without credentials or a request body. The source specifies:

- The endpoint: `GET /health`
- The response status, Content-Type header, and exact body bytes
- A latency upper bound (100 ms, P100) measured at a named process boundary
- The concurrency scope within which the latency bound applies (≤ 50 concurrent health requests in flight)
- An explicit exclusion of non-GET methods, request bodies, error responses, and behavior above 50 concurrent requests

The source is unusually precise. It names the measurement boundary ("the Health API process boundary"), specifies P100 (every request must meet the deadline), defines "concurrent" in terms of an observable interval (complete-header receipt to final-byte emission), and explicitly delegates out-of-scope behaviors to separate requirements.

The source text joins three independently testable obligations in one prose block: (1) status code and headers, (2) response body, and (3) latency under concurrency. These are split into separate candidates below because they can be verified and negotiated independently.

---

## Clarifying Questions

No material question remains unresolved. All items below are observations or confirmatory notes, not blocking ambiguities.

1. **Content-Type parameters:** The source says "Content-Type application/json with no parameters." This is interpreted as the literal header value `application/json` with no `;charset=…` or other parameters appended. If an existing framework unconditionally appends `; charset=utf-8`, confirm whether that constitutes a failure.

2. **Byte sequence encoding:** The source specifies the exact byte sequence `{"status":"ok"}` (15 bytes, UTF-8, no BOM, no trailing newline). If the implementation may emit a trailing newline or whitespace, confirm whether that is acceptable.

3. **P100 clarification:** The source explicitly states P100 ("every request must meet it"). This is treated as a hard ceiling, not a percentile target. No statistical relaxation is assumed.

4. **Concurrency scope boundary:** The source says the latency requirement applies "while no more than 50 such requests are in flight at the process boundary." Behavior above 50 concurrent requests is out of scope here. This is accepted as stated.

5. **HTTP version:** The source does not specify HTTP/1.1 or HTTP/2. The requirement is drafted version-agnostically. If a specific version matters, it should be added to the implementability contract.

---

## Candidate Requirements

### Candidate 1

- **Status:** ready for review
- **Template:** Complex (While … When … shall …)
- **Requirement:** While no more than 50 health requests are concurrently in flight at the Health API process boundary, when a client sends a GET /health request with no request body, the Health API shall return HTTP status 200.
- **Source intent and rationale:** The source obligates a 200 response to any GET /health request with no body. The concurrency scope (≤ 50 in-flight) is a precondition that limits when this obligation is stated to hold; out-of-scope behavior is explicitly deferred. The complex pattern is chosen because both the state (≤ 50 in-flight) and the event (receipt of a complete GET /health request with no body) are necessary to scope the obligation.
- **Implementability contract**
  - System boundary: The Health API process boundary — the point at which complete request headers are received and the final response byte is emitted.
  - Responsible system: Health API.
  - Actor: Any client (unauthenticated).
  - Affected interfaces: HTTP `GET /health` endpoint.
  - Inputs and data definitions: An HTTP GET request to the path `/health` with no request body. "No request body" means Content-Length: 0 or absent, and no Transfer-Encoding body payload.
  - Preconditions, state, feature scope, and trigger: Precondition — the count of health requests concurrently in flight (complete-header receipt to final-byte emission) at the process boundary is ≤ 50. Trigger — receipt of a complete GET /health request with no request body.
  - Observable normal, alternate, and failure response: The HTTP response status line contains status code 200. No alternate or failure path is in scope for this candidate.
  - Quantities, units, limits, timing, ordering, and tolerances: Concurrency limit: 50 in-flight requests (inclusive upper bound).
  - Domain terms, dependencies, and unresolved decisions: "In flight" — a health request is in flight from the moment the process boundary receives its complete request headers until the moment the final response byte is emitted at that same boundary. "Process boundary" — the named measurement and scoping boundary for all obligations in this package.
- **Verification contract**
  - Setup, state, input data, and stimulus: Configure the Health API. Establish a steady state with ≤ 50 concurrent GET /health requests in flight. Send one additional GET /health request with no request body.
  - Observable evidence or oracle: The HTTP response received by the client.
  - Expected result: The response status code is 200.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — status code equals 200. Fail — any other status code.
  - Boundary, negative, and failure cases: Lower boundary: zero other concurrent requests in flight (solo request). Upper boundary: exactly 49 other concurrent requests in flight (total = 50 including this one). Requests beyond 50 in-flight are out of scope.
  - Required data, instrumentation, assessment, or external evidence: A client capable of issuing HTTP GET requests and reading the response status code. A mechanism to control and observe concurrency at the process boundary.
  - Observable at the named boundary alone: yes

---

### Candidate 2

- **Status:** ready for review
- **Template:** Complex (While … When … shall …)
- **Requirement:** While no more than 50 health requests are concurrently in flight at the Health API process boundary, when a client sends a GET /health request with no request body, the Health API shall return a Content-Type response header with the value `application/json` and no parameters.
- **Source intent and rationale:** The source specifies the Content-Type as `application/json` "with no parameters," explicitly excluding semicolon-delimited additions such as `; charset=utf-8`. This is a distinct observable header obligation, split from status and body for independent verifiability.
- **Implementability contract**
  - System boundary: Health API process boundary.
  - Responsible system: Health API.
  - Actor: Any client (unauthenticated).
  - Affected interfaces: HTTP `GET /health` response headers.
  - Inputs and data definitions: Same as Candidate 1.
  - Preconditions, state, feature scope, and trigger: Same as Candidate 1.
  - Observable normal, alternate, and failure response: The HTTP response includes a `Content-Type` header whose value is exactly `application/json` with no parameters (no semicolon, no charset, no boundary, no other parameter).
  - Quantities, units, limits, timing, ordering, and tolerances: None beyond the exact header value.
  - Domain terms, dependencies, and unresolved decisions: "No parameters" — the header value field contains no semicolon-separated token after `application/json`. Confirmatory note: if a framework unconditionally appends `; charset=utf-8`, that constitutes a failure under this requirement.
- **Verification contract**
  - Setup, state, input data, and stimulus: Same as Candidate 1. Inspect the response headers.
  - Observable evidence or oracle: The value of the `Content-Type` response header.
  - Expected result: `Content-Type: application/json` with no parameters appended.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — header value is exactly `application/json`. Fail — header is absent, has a different media type, or has any appended parameters.
  - Boundary, negative, and failure cases: Header absent (fail). Value `application/json; charset=utf-8` (fail). Value `application/json` (pass).
  - Required data, instrumentation, assessment, or external evidence: A client or proxy capable of reading raw response headers.
  - Observable at the named boundary alone: yes

---

### Candidate 3

- **Status:** ready for review
- **Template:** Complex (While … When … shall …)
- **Requirement:** While no more than 50 health requests are concurrently in flight at the Health API process boundary, when a client sends a GET /health request with no request body, the Health API shall return a response body consisting of exactly the byte sequence `{"status":"ok"}`.
- **Source intent and rationale:** The source specifies an exact byte sequence. "Exactly" excludes leading/trailing whitespace, newlines, BOMs, or any variation in key order, spacing, or quoting. This is split from the status and Content-Type obligations because it is independently testable and negotiable.
- **Implementability contract**
  - System boundary: Health API process boundary.
  - Responsible system: Health API.
  - Actor: Any client (unauthenticated).
  - Affected interfaces: HTTP `GET /health` response body.
  - Inputs and data definitions: Same as Candidate 1.
  - Preconditions, state, feature scope, and trigger: Same as Candidate 1.
  - Observable normal, alternate, and failure response: The response body is exactly the 15 UTF-8 bytes: `{`, `"`, `s`, `t`, `a`, `t`, `u`, `s`, `"`, `:`, `"`, `o`, `k`, `"`, `}`. No BOM, no trailing newline, no leading whitespace, no alternative key or value.
  - Quantities, units, limits, timing, ordering, and tolerances: Body length: 15 bytes. Encoding: UTF-8 (all characters are ASCII, so encoding is unambiguous).
  - Domain terms, dependencies, and unresolved decisions: "Exactly the byte sequence" — a byte-for-byte comparison, not a semantic JSON equivalence. The source uses this phrase deliberately.
- **Verification contract**
  - Setup, state, input data, and stimulus: Same as Candidate 1. Capture the raw response body bytes.
  - Observable evidence or oracle: A byte-level comparison of the response body against the reference sequence `{"status":"ok"}` (15 bytes).
  - Expected result: The response body is byte-for-byte identical to `{"status":"ok"}`.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — body length is 15 bytes and all 15 bytes match the reference. Fail — any deviation in length, byte value, or position.
  - Boundary, negative, and failure cases: Body with trailing newline (`{"status":"ok"}\n`, 16 bytes) — fail. Body with spaces (`{ "status": "ok" }`) — fail. Semantically equivalent JSON with different key order — fail.
  - Required data, instrumentation, assessment, or external evidence: A client capable of capturing raw response body bytes. A byte-comparison tool or assertion.
  - Observable at the named boundary alone: yes

---

### Candidate 4

- **Status:** ready for review
- **Template:** Complex (While … When … shall …)
- **Requirement:** While no more than 50 health requests are concurrently in flight at the Health API process boundary, when a client sends a GET /health request with no request body, the Health API shall emit the final response byte within 100 milliseconds of receiving the complete request headers, measured at the Health API process boundary.
- **Source intent and rationale:** The source establishes a hard P100 latency ceiling of 100 ms. "P100: every request must meet it" means there is no statistical relaxation; every individual request must satisfy the bound. The measurement window is explicitly defined: from receipt of the complete request headers to emission of the final response byte, both at the Health API process boundary. The concurrency scope (≤ 50 in-flight) limits when this obligation applies.
- **Implementability contract**
  - System boundary: Health API process boundary — both the start event (receipt of complete request headers) and the end event (emission of final response byte) are measured at this same named boundary.
  - Responsible system: Health API.
  - Actor: Any client (unauthenticated).
  - Affected interfaces: HTTP `GET /health` response latency at the process boundary.
  - Inputs and data definitions: Same as Candidate 1.
  - Preconditions, state, feature scope, and trigger: Precondition — the count of health requests concurrently in flight at the process boundary is ≤ 50. Trigger — receipt of complete request headers for a GET /health request with no request body.
  - Observable normal, alternate, and failure response: The final byte of the HTTP response is emitted at the process boundary no later than 100 ms after the complete request headers were received at that boundary.
  - Quantities, units, limits, timing, ordering, and tolerances: Maximum latency: 100 ms (milliseconds). Coverage: P100 — every request without exception. Measurement start: complete-header receipt at the process boundary. Measurement end: final-byte emission at the process boundary. Concurrency scope: ≤ 50 health requests in flight simultaneously (inclusive).
  - Domain terms, dependencies, and unresolved decisions: "Complete request headers" — the HTTP request headers have been fully received (e.g., the blank line terminating the header section in HTTP/1.1, or the HEADERS frame end in HTTP/2). "Final response byte" — the last byte of the HTTP response, including any chunked-transfer terminal chunk or stream-end frame. "Process boundary" — the named measurement boundary; network round-trip time, TCP stack, and OS socket buffering outside this boundary are excluded. "In flight" — see Candidate 1.
- **Verification contract**
  - Setup, state, input data, and stimulus: Instrument the Health API process boundary to record (a) the timestamp of complete-header receipt and (b) the timestamp of final-byte emission, for each request. Issue GET /health requests with no body under a concurrency load of 1 to 50 simultaneous requests.
  - Observable evidence or oracle: The per-request latency value computed as (final-byte emission timestamp) − (complete-header receipt timestamp), in milliseconds, at the process boundary.
  - Expected result: Every observed latency value is ≤ 100 ms.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass — all requests in the test run have latency ≤ 100 ms. Fail — any single request has latency > 100 ms. The criterion is P100: one failure is a failure of the requirement.
  - Boundary, negative, and failure cases: Solo request (1 in-flight): must meet 100 ms. Maximum in-scope concurrency (50 in-flight): every one of the 50 must meet 100 ms. Requests above 50 concurrent are out of scope and need not be measured here.
  - Required data, instrumentation, assessment, or external evidence: Instrumentation at the Health API process boundary capable of timestamping complete-header receipt and final-byte emission with sub-millisecond resolution. Test tooling capable of driving 1–50 concurrent GET /health requests. A pass/fail report covering every observed latency value.
  - Observable at the named boundary alone: yes (the measurement boundary is explicitly named and both events are at that boundary; no internal probes are required by this requirement, though process-boundary instrumentation is needed)

---

## Suggested Supporting Requirements

**S1 — Failure-path companion (Candidate 1, 2, 3):**
The source explicitly defers error responses to separate requirements. A companion requirement defining the HTTP status code and body for requests that arrive when > 50 health requests are in flight is implied but out of scope here. Confirm that the separate requirements covering this case exist and reference the same process-boundary definition.

**S2 — Boundary companion (Candidate 4):**
The source defines concurrency in terms of a specific measurement interval (complete-header receipt to final-byte emission). A companion requirement defining how the Health API counts or observes concurrency at the process boundary (e.g., via an atomic counter incremented at complete-header receipt and decremented at final-byte emission) may be needed so that two implementers arrive at the same in-flight count. This is suggested, not required, if the mechanism is left to implementation choice.

**S3 — Interface companion (Candidates 1–4):**
The source does not specify an HTTP version. If the Health API must support both HTTP/1.1 and HTTP/2 (or HTTP/3), a companion requirement scoping the applicable protocol versions would prevent ambiguity. This is an optional consideration — confirm whether protocol version is in scope.

**S4 — Operational companion (Candidate 4):**
The source does not specify an availability or uptime obligation for the endpoint. A companion requirement defining availability (e.g., the endpoint is available 99.9% of calendar time) is a common accompaniment to a health endpoint but is not stated here. Raise only if operational SLAs are in scope.

---

## Consistency Findings

**Finding 1 — Intentional exception, out-of-scope behaviors:**
Candidates 1–4 all carry the precondition "no more than 50 health requests are concurrently in flight." The source explicitly states that behavior above 50 concurrent requests is covered by separate requirements. There is no conflict within this package; the scope exclusion is intentional. The separate requirements must share the same "in-flight" and "process boundary" definitions to avoid a gap at exactly 50 concurrent requests.

**Finding 2 — Gap, exactly-50 boundary:**
The source says the requirement applies "while no more than 50 such requests are in flight." This includes exactly 50. The separate requirements covering > 50 concurrent requests must define whether the transition is based on a strict inequality (> 50) or ≥ 51. If the separate requirements use "> 50" they will be consistent with this package's "≤ 50"; if they use "≥ 50" there is a conflict at exactly 50. Flag this for the authors of the companion requirements.

**Finding 3 — No conflict among Candidates 1–4:**
The four candidates are parallel obligations triggered by the same event under the same precondition. Each covers a distinct observable: status code, Content-Type header, response body, and latency. They do not contradict each other.

**Finding 4 — Overlap, byte sequence and Content-Type:**
Candidate 2 (Content-Type) and Candidate 3 (body) together define the full HTTP response payload contract. There is no overlap between them, but they should be verified together in a single test scenario to confirm coherence (a response can satisfy one and fail the other).

---

## Assumptions

1. **"Process boundary" is the HTTP layer at the application process, not the network interface.** The measurement excludes OS TCP/IP stack buffering, load-balancer latency, and client-side network delay. This is the most natural reading of "Health API process boundary" and is consistent with the measurement being symmetric (same boundary for start and end).

2. **"No request body" means the absence of a body payload, not a body of zero declared length.** A request with `Content-Length: 0` and no body bytes is treated as having no request body. This is the common HTTP convention.

3. **The exact byte sequence `{"status":"ok"}` is UTF-8 encoded with no BOM.** All 15 characters are in the ASCII range, so the encoding is unambiguous, but the assumption is stated explicitly.

4. **"Every request must meet it" (P100) means no statistical exception.** A single request exceeding 100 ms is a failure. This is the plain reading of the source.

5. **The concurrency count of 50 is inclusive.** "No more than 50" means the count may be 0, 1, 2, … 50. The obligation does not apply when the count is 51 or higher.

---

## Rationale and Context

The source describes a standard liveness/readiness probe endpoint common in cloud-native deployments. Such endpoints are called frequently by orchestration platforms (e.g., Kubernetes) and load balancers to determine whether a service instance is healthy. The tight latency bound (100 ms, P100) and the explicit concurrency scope (≤ 50 in-flight) reflect operational requirements where slow or incorrect health responses cause cascading restarts or traffic misrouting. The exact body requirement (`{"status":"ok"}`) enables downstream tooling to perform byte-level assertions rather than JSON parsing, reducing probe complexity.

The source is unusually well-specified for a prose requirement. The measurement boundary, timing model, concurrency definition, and scope exclusions are all explicit. The main elicitation work is decomposition (splitting three joined obligations into independently verifiable candidates) and surfacing the boundary-condition gap at exactly 50 concurrent requests.

---

## Review Status

| Candidate | Status | Blocking finding |
|---|---|---|
| Candidate 1 — HTTP 200 status | ready for review | none |
| Candidate 2 — Content-Type header | ready for review | none |
| Candidate 3 — Exact response body | ready for review | none |
| Candidate 4 — 100 ms P100 latency | ready for review | none |

All four candidates have passed both the implementability contract and the verification contract. No material question remains unresolved. The suggested supporting requirements (S1–S4) and the boundary-condition gap (Finding 2) are flagged for the authors of the companion requirements that cover out-of-scope behaviors. No candidate is blocked.
