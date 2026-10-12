# Requirements Elicitation Package

## Understanding

**Requested outcome:** When the upstream identity provider fails to respond within an acceptable time, the gateway must protect callers from an indefinite wait by returning a standard HTTP 503 response and a `Retry-After` header.

**Proposed system boundary:** The **gateway** — the component that sits between callers and the upstream identity provider. The gateway is responsible for detecting the identity-provider timeout and producing the observable HTTP response to the caller.

**Source classification:**

| Source fragment | Kind |
|---|---|
| "If the upstream identity provider times out" | Unwanted condition |
| "the gateway shall return HTTP 503" | Behavioral obligation |
| "a Retry-After header" | Behavioral obligation (companion response) |
| "to the caller" | Interface/actor identification |

The source uses `shall`, matching normative obligation language. The condition is an unwanted failure event (timeout of an external dependency), which selects the **Unwanted behavior** EARS pattern. Two independently testable responses are implied — the status code and the header — but because both are part of a single observable HTTP response to one event, they are captured in one candidate with a compound response clause. If they must be negotiated or verified independently, they should be split; a clarifying question is raised below.

---

## Clarifying Questions

1. **Retry-After value (material):** What value shall the gateway place in the `Retry-After` header? Options include a fixed configured duration, a value derived from the identity provider's own retry guidance, a calculated back-off, or a product of policy. Two implementers will produce different behavior without an answer.

2. **Timeout definition (material):** How is "times out" defined at the gateway boundary? Is it: (a) no response received within a configured interval, (b) a TCP/TLS connection failure, (c) an application-layer keepalive or health check expiry, or (d) something else? The verification contract cannot be completed without this.

3. **Scope of the condition (material):** Does this requirement apply only while the gateway is actively proxying an inbound caller request, or does it also apply to background health checks and proactive token refresh flows? The answer determines whether a `While` precondition is needed.

4. **Existing in-flight request body:** If the gateway has already begun streaming a response body to the caller when the identity-provider timeout occurs, shall it still return HTTP 503, or is this path out of scope?

5. **Single vs. split candidates:** Shall the HTTP 503 status and the `Retry-After` header be treated as independently negotiable obligations (requiring two separate candidates), or as a single atomic response? The current draft treats them as one atomic response.

6. **`Retry-After` format:** RFC 9110 permits either an HTTP-date or a delay-seconds integer. Is one form required, or is the format left to the implementation?

---

## Candidate Requirements

### Candidate 1

- **Status:** needs clarification
- **Template:** Unwanted behavior
- **Requirement:** If the upstream identity provider times out, then the gateway shall return HTTP status code 503 and a `Retry-After` header to the caller.

**Source intent and rationale:**  
The source requests that callers receive a standards-conformant signal when the identity provider is unreachable due to a timeout, allowing them to retry after an appropriate delay rather than receiving an ambiguous error or waiting indefinitely.

**Implementability contract**

- **System boundary:** Gateway — the intermediary that forwards authentication or authorization requests to the upstream identity provider and returns responses to callers.
- **Responsible system:** Gateway.
- **Actor:** Caller (the client that sent the original request to the gateway).
- **Affected interfaces:** Outbound interface from the gateway to the caller (HTTP response); inbound interface from the gateway to the upstream identity provider (where the timeout is detected).
- **Inputs and data definitions:** An inbound caller request that causes the gateway to contact the upstream identity provider. A timeout event on the gateway-to-identity-provider connection or response.
- **Preconditions, state, feature scope, and trigger:** The gateway has initiated contact with the upstream identity provider on behalf of a caller. The identity provider does not respond within the timeout threshold. (The exact threshold definition is **open** — see Question 2.)
- **Observable normal, alternate, and failure response:** The gateway returns an HTTP response with status code `503 Service Unavailable` and a `Retry-After` header to the caller. The value of `Retry-After` is **unresolved** — see Question 1. The format of `Retry-After` (HTTP-date vs. delay-seconds) is **unresolved** — see Question 6.
- **Quantities, units, limits, timing, ordering, and tolerances:** Timeout threshold: **unresolved** (see Question 2). `Retry-After` value: **unresolved** (see Question 1).
- **Domain terms, dependencies, and unresolved decisions:**
  - *Upstream identity provider*: an external system responsible for authenticating or authorizing the caller; assumed to be contacted synchronously per caller request. Confirm if background flows are in scope (Question 3).
  - *Times out*: definition is **open** (Question 2).
  - *Retry-After*: defined by RFC 9110 §10.2.4; format is **open** (Question 6).
  - *Caller*: the client sending the original HTTP request to the gateway.

**Verification contract**

- **Setup, state, input data, and stimulus:** Configure the gateway with a known timeout threshold for the identity provider. Direct the identity provider (or a test double) to withhold its response beyond that threshold. Send a caller request to the gateway that requires identity-provider contact.
- **Observable evidence or oracle:** The HTTP response received by the caller.
- **Expected result:** HTTP status code is exactly `503`; response includes a `Retry-After` header with a value conforming to RFC 9110.
- **Pass/fail criteria:** Pass if and only if status code equals `503` AND `Retry-After` header is present AND its value is syntactically valid per RFC 9110. Fail if status code differs, header is absent, or header value is malformed. The specific `Retry-After` value assertion is **pending** resolution of Question 1.
- **Boundary, negative, and failure cases:** (a) Identity provider responds just within the threshold — gateway must not return 503. (b) Identity provider times out on one of multiple identity providers (if more than one is configured) — behavior is **unresolved**. (c) Caller disconnects before the timeout is detected — out of scope unless Question 3 or 4 brings it in.
- **Required data, instrumentation, assessment, or external evidence:** Ability to control or simulate identity-provider response delay; HTTP response capture at the caller boundary.
- **Observable at the named boundary alone:** Yes — the HTTP status code and headers are observable at the gateway–caller interface without internal instrumentation, provided the timeout threshold definition is known.
- **If no, internal evidence needed and rationale:** N/A.

---

## Suggested Supporting Requirements

**SR-1 — Boundary companion:** *Timeout threshold configuration*  
The gateway should expose a configurable timeout threshold for the upstream identity-provider connection and response. Without this, the trigger condition in Candidate 1 cannot be set or verified independently of a deployment-specific default. Relationship: boundary companion to Candidate 1.

**SR-2 — Failure-path companion:** *Logging the timeout event*  
If the upstream identity provider times out, then the gateway should record a structured log entry that identifies the identity provider, the elapsed time, and the caller request reference. This supports operational diagnosis without depending on the caller's 503 response as the sole evidence. Relationship: failure-path companion to Candidate 1.

**SR-3 — Boundary companion:** *Behavior when the identity provider is unavailable from the start (connection refused or DNS failure)*  
The source addresses a timeout; connection-refused and DNS-failure conditions may require the same or different responses. If they require the same 503 + Retry-After response, a separate or broadened requirement is needed. If they are out of scope, confirm explicitly. Relationship: boundary companion to Candidate 1.

**SR-4 — Optional consideration:** *Rate-limiting repeated retries from the same caller*  
After returning 503 + Retry-After, the gateway may receive immediate retry requests that further load the failing identity provider. A rate-limiting or circuit-breaker requirement may be warranted. This is not justified by the source alone; ask whether it is in scope. Relationship: optional consideration.

**SR-5 — Interface companion:** *Identity provider timeout propagation when multiple providers are configured*  
If the gateway may contact more than one upstream identity provider, the requirement should state whether a timeout of one provider triggers 503 immediately or only after all providers are tried. Relationship: interface companion to Candidate 1.

---

## Consistency Findings

**Finding 1 — Gap:** The source does not define behavior when the identity provider returns an error response (e.g., HTTP 500 or 401) rather than timing out. The Candidate 1 trigger is specifically a timeout; error responses are a separate gap not covered by this package.

**Finding 2 — Gap:** The `Retry-After` value is mentioned in the source but its determination is undefined. Without a companion requirement or constraint on its value, two conformant implementations could return wildly different retry intervals, making interoperability testing impossible.

**Finding 3 — Precedence question:** If another requirement governs the gateway's general error-response format (e.g., a JSON error body or a specific `Content-Type`), it is unclear whether that requirement takes precedence, is additive, or conflicts with the 503 response defined here. No such requirement was supplied; flag for review if one exists.

No conflicts or duplicates were found within this package (only one candidate is present).

---

## Assumptions

1. **A1 (provisional):** "Times out" means the upstream identity provider does not return any response within a configured time interval measured from when the gateway initiates the request. This is the most common interpretation; it must be confirmed (Question 2).

2. **A2 (provisional):** The `Retry-After` header value is a fixed or configured duration expressed as a delay-seconds integer, unless the identity provider supplies its own retry guidance. This must be confirmed (Question 1 and 6).

3. **A3 (provisional):** The requirement applies only while the gateway is servicing an active inbound caller request — not to background health-check or token-refresh flows. Confirm via Question 3.

4. **A4 (provisional):** The gateway contacts exactly one upstream identity provider per caller request. If multiple providers may be tried in sequence or parallel, SR-5 becomes a required companion.

---

## Rationale and Context

The source models a standard resilience pattern: when a critical upstream dependency fails to respond, a service gateway should return a well-known HTTP failure code (`503 Service Unavailable`) and signal to the caller when it may retry (`Retry-After`), rather than returning an ambiguous error or leaving the caller blocked. This allows callers to implement compliant back-off strategies and avoids thundering-herd re-connections during an outage.

The `Retry-After` header is defined in RFC 9110 §10.2.4. Its presence is a normative part of the requested behavior, but its value and format are not specified in the source, making those the highest-priority open questions for this package.

---

## Review Status

| Item | Status |
|---|---|
| Candidate 1 | **needs clarification** — `Retry-After` value (Q1), timeout definition (Q2), and scope (Q3) are unresolved; the verification contract cannot be completed without them. |
| SR-1 through SR-5 | Suggestions only; not yet approved candidates. |
| Consistency findings | 3 findings (2 gaps, 1 precedence question); no conflicts or duplicates detected. |
| Overall package | Not ready for review. Resolve Q1 and Q2 at minimum before returning for another elicitation pass. |
