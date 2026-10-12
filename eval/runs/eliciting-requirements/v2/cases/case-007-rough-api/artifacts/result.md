# Requirements Elicitation Package

## Understanding

**Requested outcome:** Users can submit invoice documents through a browser interface. The system accepts each submission, processes it asynchronously, and persists the processed result. Invoices that are malformed, invalid, or otherwise unprocessable are handled without silent failure.

**Proposed system boundary:** The system under specification is a web-based invoice upload and processing service. It exposes a browser interface for submissions, an asynchronous processing pipeline, and a persistence layer for results.

**Source material classification:**

| Statement | Kind |
|---|---|
| "Users can upload invoices" | Goal / behavioral obligation |
| "Use a React form" | Design idea (framework choice — kept out of normative sentences) |
| "put the work on a Redis queue" | Design idea (queue technology — kept out of normative sentences) |
| "store the result in a PostgreSQL table" | Design idea (storage technology — kept out of normative sentences) |
| "handle bad invoices gracefully" | Goal / behavioral obligation (vague — needs clarification on "bad" and "gracefully") |

The technology choices (React, Redis, PostgreSQL) are implementation constraints the requester has imposed. They are preserved here as design context and in the Assumptions section, but they are intentionally excluded from normative sentences because requirements describe observable behavior at system boundaries, not implementation structures.

---

## Clarifying Questions

The following questions represent the highest-value gaps. Each, if left unanswered, would allow two reasonable implementers or verifiers to produce visibly different behavior.

**Q1 — Accepted file formats and size limits (affects C1, C2)**
What file types may a user upload (e.g., PDF, PNG, CSV, XML)? Is there a maximum file size? What happens when an unsupported type or oversized file is submitted?

**Q2 — Definition of "bad invoice" (affects C4, C5)**
"Bad invoices" is ambiguous. Does it include: malformed files (corrupt, unreadable), files of a wrong type, files with missing required fields (e.g., no invoice number, no total), invoices that fail a business rule (e.g., amount exceeds a threshold), or all of the above? The answer determines how many distinct failure candidates are needed.

**Q3 — Definition of "gracefully" (affects C4, C5)**
What observable behavior is expected when a bad invoice is detected?
- Is the user notified in the browser at the time of submission, asynchronously after processing, or both?
- Is a failed invoice persisted (e.g., with an error status), discarded, or returned to the user?
- Is the user allowed to correct and resubmit?

**Q4 — User identity and authentication (affects C1, C4)**
Must the submitting user be authenticated? Is invoice ownership tracked per user, or are all submissions treated as anonymous? This affects the observable submission response (e.g., whether a user-specific identifier is returned or stored).

**Q5 — Submission response and acknowledgement (affects C1)**
When a user submits a valid invoice, what does the system return immediately? Options include: a submission identifier the user can reference, a status indicator, or no immediate acknowledgement. The answer determines whether a synchronous acknowledgement requirement is needed.

**Q6 — Result observability and retrieval (affects C3)**
After a submitted invoice is processed, how does a user or downstream consumer access the stored result? Is there a retrieval endpoint? Can the user check status? If no retrieval behavior is specified, the stored result is unobservable at the named boundary.

**Q7 — Duplicate submissions (affects C1)**
What happens when a user uploads the same invoice file more than once? Is each submission treated independently, or is deduplication required?

**Q8 — Processing timeout (affects C4)**
If asynchronous processing does not complete within some time window, is that treated as a failure? What is the time boundary, and is the user notified?

---

## Candidate Requirements

### Candidate 1

- **Status:** needs clarification
- **Template:** Event-driven
- **Requirement:** When a user submits an invoice file, the Invoice Upload Service shall accept the submission and return an acknowledgement to the user.
- **Source intent and rationale:** The primary stated goal is that users can upload invoices. A submission event at the browser interface triggers the system. The immediate response (acknowledgement) is implied but not specified — Q1 (file types, size) and Q5 (what acknowledgement looks like) are open.
- **Implementability contract**
  - System boundary: Browser interface exposed by the Invoice Upload Service
  - Responsible system: Invoice Upload Service
  - Actor: Authenticated or anonymous user (open — see Q4)
  - Affected interfaces: Browser form submission endpoint (HTTP request/response)
  - Inputs and data definitions: Invoice file — format(s) and size limits are unresolved (Q1)
  - Preconditions, state, feature scope, and trigger: User submits an invoice file via the browser interface
  - Observable normal, alternate, and failure response: System accepts the submission; returns an acknowledgement — content and form of acknowledgement are unresolved (Q5); behavior for unsupported file types is unresolved (Q1)
  - Quantities, units, limits, timing, ordering, and tolerances: File size limit unresolved (Q1); response timing unresolved (Q5)
  - Domain terms, dependencies, and unresolved decisions: "submit" = user action to send a file; "acknowledgement" = undefined observable response (Q5); accepted file formats undefined (Q1); authentication requirement undefined (Q4)
- **Verification contract**
  - Setup, state, input data, and stimulus: A user opens the invoice submission interface and submits a file of an accepted type within any size limit
  - Observable evidence or oracle: HTTP response received by the browser after submission
  - Expected result: The system returns an acknowledgement (content undefined — Q5)
  - Pass/fail criteria, including scope, thresholds, units, and timing: Pass when an acknowledgement is received; fail when no response or an error is returned — thresholds unresolved
  - Boundary, negative, and failure cases: Unsupported file type (Q1), oversized file (Q1), unauthenticated user (Q4), duplicate submission (Q7) — all unresolved
  - Required data, instrumentation, assessment, or external evidence: Access to the submission endpoint; a valid invoice file in an accepted format
  - Observable at the named boundary alone: yes
  - If no, internal evidence needed and rationale: N/A

---

### Candidate 2

- **Status:** needs clarification
- **Template:** Event-driven
- **Requirement:** When a user submits an invoice file, the Invoice Upload Service shall validate the submission before accepting it.
- **Source intent and rationale:** "Handle bad invoices gracefully" implies that the system distinguishes valid from invalid invoices. Validation at the submission boundary is a necessary companion to C1. The specific validation rules are unresolved (Q1, Q2).
- **Implementability contract**
  - System boundary: Browser interface / submission endpoint
  - Responsible system: Invoice Upload Service
  - Actor: User
  - Affected interfaces: Browser form submission endpoint
  - Inputs and data definitions: Invoice file; validation rules are unresolved — file type, size, required fields, business rules (Q1, Q2)
  - Preconditions, state, feature scope, and trigger: User submits a file
  - Observable normal, alternate, and failure response: If validation passes, submission is accepted; if validation fails, submission is rejected — rejection response content is unresolved (Q2, Q3)
  - Quantities, units, limits, timing, ordering, and tolerances: Validation rule thresholds unresolved (Q1, Q2)
  - Domain terms, dependencies, and unresolved decisions: "validate" — rules are unresolved; "bad invoice" is undefined (Q2)
- **Verification contract**
  - Setup, state, input data, and stimulus: Submit a file known to be invalid (e.g., wrong type, missing field) and a file known to be valid
  - Observable evidence or oracle: System response distinguishes accepted from rejected submissions
  - Expected result: Valid file accepted; invalid file rejected with an observable signal
  - Pass/fail criteria, including scope, thresholds, units, and timing: Unresolved — validation rules not yet defined
  - Boundary, negative, and failure cases: Each category of "bad invoice" (Q2); boundary file sizes (Q1)
  - Required data, instrumentation, assessment, or external evidence: Test files for each validation case; knowledge of defined validation rules
  - Observable at the named boundary alone: yes
  - If no, internal evidence needed and rationale: N/A

---

### Candidate 3

- **Status:** needs clarification
- **Template:** Event-driven
- **Requirement:** When an accepted invoice is ready for processing, the Invoice Processing Service shall process the invoice and persist the result.
- **Source intent and rationale:** The source requests asynchronous processing ("put the work on a queue") and storage of the result ("store the result"). This requirement describes the observable outcome of that pipeline: a processed result is persisted. The content of "the result" and how it is observable are not specified (Q6).
- **Implementability contract**
  - System boundary: Processing pipeline output / persistence boundary
  - Responsible system: Invoice Processing Service
  - Actor: System (triggered by a queued work item, not directly by a user action at this stage)
  - Affected interfaces: Persistence store (write interface); optional retrieval interface (Q6)
  - Inputs and data definitions: Accepted invoice (format as resolved in C1/C2); processed result — fields and schema are unresolved (Q6)
  - Preconditions, state, feature scope, and trigger: An accepted invoice is placed in the processing queue
  - Observable normal, alternate, and failure response: Result is persisted — observable schema and retrieval method are unresolved (Q6); failure path during processing is addressed in C4/C5
  - Quantities, units, limits, timing, ordering, and tolerances: Processing time limit unresolved (Q8); result retention period unresolved
  - Domain terms, dependencies, and unresolved decisions: "result" — content and structure undefined (Q6); retrieval behavior undefined (Q6)
- **Verification contract**
  - Setup, state, input data, and stimulus: An accepted invoice is submitted and queued for processing
  - Observable evidence or oracle: After processing completes, the result is retrievable (Q6 — retrieval method undefined)
  - Expected result: A result record exists and is accessible
  - Pass/fail criteria, including scope, thresholds, units, and timing: Unresolved — result content, schema, and retrieval path undefined (Q6)
  - Boundary, negative, and failure cases: Processing failure (C4); timeout (Q8); duplicate processing (Q7)
  - Required data, instrumentation, assessment, or external evidence: Access to the persistence store or retrieval endpoint; a known accepted invoice as input
  - Observable at the named boundary alone: no
  - If no, internal evidence needed and rationale: The persistence of the result is only verifiable by reading the persistence store or a retrieval endpoint, neither of which is yet defined (Q6). Without a defined retrieval path, the verification contract cannot be completed.

---

### Candidate 4

- **Status:** needs clarification
- **Template:** Unwanted behavior
- **Requirement:** If a submitted invoice fails validation, then the Invoice Upload Service shall reject the submission and notify the user of the reason for rejection.
- **Source intent and rationale:** "Handle bad invoices gracefully" implies the user is informed when their invoice is not accepted. This candidate covers the synchronous validation-failure path. "Gracefully" and "reason for rejection" are not yet measurable (Q2, Q3).
- **Implementability contract**
  - System boundary: Browser interface / submission endpoint
  - Responsible system: Invoice Upload Service
  - Actor: User
  - Affected interfaces: Browser form submission response
  - Inputs and data definitions: Invalid invoice file; rejection reason — defined by validation rule set (Q1, Q2)
  - Preconditions, state, feature scope, and trigger: User submits a file that fails validation
  - Observable normal, alternate, and failure response: Submission is rejected; user receives a notification with the rejection reason — form and content of notification are unresolved (Q3)
  - Quantities, units, limits, timing, ordering, and tolerances: Notification timing (synchronous vs. deferred) unresolved (Q3)
  - Domain terms, dependencies, and unresolved decisions: "gracefully" — no agreed measure (Q3); "reason for rejection" — message format and specificity undefined
- **Verification contract**
  - Setup, state, input data, and stimulus: Submit an invoice file that fails a known validation rule
  - Observable evidence or oracle: The submission endpoint response contains a rejection indicator and a reason
  - Expected result: Submission rejected; reason message visible to user
  - Pass/fail criteria, including scope, thresholds, units, and timing: Unresolved — "gracefully" and reason content are not measurable without Q2/Q3 answers
  - Boundary, negative, and failure cases: Each category of invalid invoice (Q2); submission of multiple invalid invoices in sequence; submission after correction
  - Required data, instrumentation, assessment, or external evidence: Invalid test files for each validation category; user-visible response content
  - Observable at the named boundary alone: yes
  - If no, internal evidence needed and rationale: N/A

---

### Candidate 5

- **Status:** needs clarification
- **Template:** Unwanted behavior
- **Requirement:** If processing of an accepted invoice fails, then the Invoice Processing Service shall record the failure and notify the submitting user.
- **Source intent and rationale:** "Handle bad invoices gracefully" also covers failures detected during asynchronous processing (after submission was accepted). This is a distinct path from synchronous validation failure (C4). Whether the user is notified asynchronously and how failure is recorded are unresolved (Q3, Q6).
- **Implementability contract**
  - System boundary: Processing pipeline / persistence and notification boundary
  - Responsible system: Invoice Processing Service
  - Actor: System (failure condition); submitting user (notification recipient)
  - Affected interfaces: Persistence store (failure record write); user notification channel (undefined — Q3)
  - Inputs and data definitions: Accepted invoice that fails during processing; failure record — fields undefined (Q6)
  - Preconditions, state, feature scope, and trigger: An accepted and queued invoice encounters a processing failure
  - Observable normal, alternate, and failure response: Failure is recorded; user is notified — notification channel, timing, and content are unresolved (Q3); whether the invoice is retried or discarded is unresolved
  - Quantities, units, limits, timing, ordering, and tolerances: Retry count and interval unresolved; notification timing unresolved (Q3)
  - Domain terms, dependencies, and unresolved decisions: "processing failure" — categories undefined (Q2); notification channel undefined (Q3); "record the failure" — schema and retrieval undefined (Q6)
- **Verification contract**
  - Setup, state, input data, and stimulus: Submit an accepted invoice that is engineered or known to fail during processing
  - Observable evidence or oracle: Failure record exists in the persistence store; user receives notification
  - Expected result: Failure is recorded and user is informed
  - Pass/fail criteria, including scope, thresholds, units, and timing: Unresolved — notification channel and failure record schema undefined (Q3, Q6)
  - Boundary, negative, and failure cases: Repeated processing failures; failure after partial processing; timeout expiry (Q8)
  - Required data, instrumentation, assessment, or external evidence: Access to the persistence store; evidence of user notification delivery; engineered failure scenario
  - Observable at the named boundary alone: no
  - If no, internal evidence needed and rationale: The failure record is only verifiable via the persistence store or a retrieval endpoint (both undefined — Q6). User notification requires evidence of delivery through the notification channel (undefined — Q3).

---

## Suggested Supporting Requirements

**SR-1 — File size and type constraint**
*Relationship:* Boundary companion to C1 and C2
*Missing behavior:* The system has no stated constraint on the file types or sizes it accepts. Without a defined boundary, the submission endpoint behavior for unsupported or oversized inputs is undefined.
*Suggestion:* Consider adding a requirement that the system rejects files that do not conform to the accepted type list or exceed a stated maximum size, and that it notifies the user of the specific constraint violated.
*Precondition:* Q1 must be answered first.

**SR-2 — Submission acknowledgement with tracking identifier**
*Relationship:* Required companion to C1
*Missing behavior:* Once a valid invoice is accepted, the user has no stated mechanism to know whether processing succeeded or failed later.
*Suggestion:* Consider requiring the system to return a unique submission identifier at acceptance time, enabling the user or a downstream system to correlate later status or failure notifications.
*Precondition:* Q5 and Q6 must be answered to determine whether a retrieval path exists.

**SR-3 — Processing status retrieval**
*Relationship:* Interface companion to C3 and C5
*Missing behavior:* The result is persisted (C3) but no retrieval behavior is specified. The verification contract for C3 and C5 cannot be completed without a defined retrieval path.
*Suggestion:* Consider specifying that the system exposes an interface by which an authorized caller can retrieve the processing status and result for a given submission.
*Precondition:* Q6 must be answered.

**SR-4 — Duplicate submission handling**
*Relationship:* Boundary companion to C1
*Missing behavior:* If a user uploads the same invoice twice, the expected behavior is undefined.
*Suggestion:* Determine whether the system should detect duplicates (and on what key — file hash, invoice number, etc.), reject them, or process each independently. Add a requirement once the policy is decided.
*Precondition:* Q7 must be answered.

**SR-5 — Processing timeout and retry**
*Relationship:* Failure-path companion to C3 and C5
*Missing behavior:* If processing does not complete within a time boundary, the outcome is undefined. The user may wait indefinitely.
*Suggestion:* Consider specifying a maximum processing duration after which the system treats the invoice as a processing failure and applies the C5 behavior (record failure, notify user).
*Precondition:* Q8 must be answered.

**SR-6 — User authentication and authorization**
*Relationship:* Required companion to C1 (if ownership is required)
*Missing behavior:* Whether users must authenticate before submitting is unstated. If invoice ownership matters (e.g., a user should only see their own processing results), authentication is a prerequisite.
*Suggestion:* Determine whether the service requires authentication and, if so, specify the authentication boundary and the behavior when an unauthenticated user attempts to submit.
*Precondition:* Q4 must be answered.

---

## Consistency Findings

**Finding 1 — Gap: C3 result is unobservable without a retrieval path**
C3 requires that a result be persisted, but no candidate or suggestion yet defines how that result is retrieved. Without a retrieval interface, the persistence obligation in C3 cannot be independently verified, and C5's failure record is similarly unobservable. SR-3 addresses this gap. Resolution: specify a retrieval interface requirement (accept SR-3, or define the observable retrieval behavior in C3).

**Finding 2 — Gap: Synchronous vs. asynchronous failure notification are not yet separated**
C4 covers synchronous validation failure (at submission time). C5 covers asynchronous processing failure (after acceptance). These are intentionally split. However, if Q3 reveals that the user is notified of both failure types through the same channel, the notification behavior should be checked for consistency between C4 and C5 to avoid contradictory messages for the same invoice.

**Finding 3 — Precedence question: Validation in C2 vs. processing failure in C5**
C2 implies validation occurs before the invoice enters the queue. C5 implies some bad invoices pass validation and fail during processing. The boundary between these two failure modes is not yet defined (Q2). If validation is complete at the boundary (C2), the processing failure in C5 should be limited to processing-specific failures (e.g., a service dependency failure, parsing error not caught at submission). Resolution: Q2 must define which failure types are caught at which stage.

**Finding 4 — Overlap risk: "Bad invoice" spans C2, C4, and C5**
The undefined scope of "bad invoice" risks overlap between C2 (validation), C4 (synchronous rejection), and C5 (asynchronous processing failure). Once Q2 is answered, re-check that each candidate covers a distinct failure category with no unintended overlap.

**Finding 5 — Gap: No candidate covers the normal user journey end-to-end**
No candidate currently specifies what the user observes after a successful end-to-end flow (submission → processing → result available). SR-2 and SR-3 address parts of this gap, but a complete happy-path observable sequence is not yet captured. Resolution: once Q5 and Q6 are answered, add a candidate or refine C1 and C3 to cover the full observable lifecycle.

---

## Assumptions

**A1 — Technology choices are design constraints, not requirements.**
The source names React, Redis, and PostgreSQL. These are treated as implementation design choices imposed by the requester, not as system behaviors. They are excluded from normative sentences. If any of these technologies impose externally visible constraints (e.g., Redis queue item size limits, PostgreSQL row-size limits), those constraints should be expressed as observable behavioral requirements when their values are known.

**A2 — Asynchronous processing is implied by "queue."**
The source's mention of a queue implies that invoice processing does not occur synchronously with the user's submission. This assumption separates the submission response (C1) from the processing and persistence behavior (C3).

**A3 — "Gracefully" means the user is informed, not that failures are silently discarded.**
In the absence of a definition, "handle bad invoices gracefully" is interpreted to mean the user receives an observable, actionable signal when their invoice is rejected or fails processing. This interpretation is provisional and subject to Q3.

**A4 — "Result" refers to a structured output of invoice processing, not the original invoice file.**
C3 uses "result" to mean whatever the processing pipeline produces (e.g., extracted data, a parsed record, a status). The schema and content of this result are undefined and must be resolved via Q6.

**A5 — The system boundary is a single service.**
The source implies a single system. If there are separate services (upload service, processing service, notification service), each requires its own boundary and responsible-system declaration. The candidates name notional services for clarity but do not prescribe a deployment architecture.

---

## Rationale and Context

The source request combines a user-facing goal ("users can upload invoices") with implementation directives (React, Redis, PostgreSQL) and a quality aspiration ("handle bad invoices gracefully"). The goal is to make this request implementable and verifiable by separating what the system must do (observable behavior) from how it is built (technology choices).

The split between C4 (synchronous rejection) and C5 (asynchronous processing failure) reflects a structural reality: once the source says "put the work on a queue," there are now two distinct failure moments with potentially different notification paths. Both must be specified, and Q2–Q3 are the minimum required to make them verifiable.

The five candidates cover the primary behavioral obligations implied by the request. All are marked `needs clarification` because one or more material questions remain open in each implementability or verification contract. No candidate has been promoted to `candidate` or `ready for review` while those questions are unresolved.

---

## Review Status

| Candidate | Status | Blocking questions |
|---|---|---|
| C1 — Invoice submission and acceptance | needs clarification | Q1, Q4, Q5, Q7 |
| C2 — Submission validation | needs clarification | Q1, Q2 |
| C3 — Processing and result persistence | needs clarification | Q6, Q8 |
| C4 — Synchronous rejection and user notification | needs clarification | Q2, Q3 |
| C5 — Asynchronous processing failure and user notification | needs clarification | Q2, Q3, Q6, Q8 |

**All five candidates are `needs clarification`.** The minimum set of questions to advance any candidate to `candidate` status is **Q2** (definition of "bad invoice") and **Q3** (definition of "gracefully"), which unblock C2, C4, and C5. Q1 (file types and size) unblocks C1 and C2. Q5 and Q6 are needed before C1 and C3 can pass their verification contracts.

No requirements have been approved. This package is the output of one elicitation pass. The candidates are provisional until the clarifying questions are answered and both quality contracts are re-evaluated.
