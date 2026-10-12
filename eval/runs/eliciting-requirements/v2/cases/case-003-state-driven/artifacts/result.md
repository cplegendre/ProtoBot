# Requirements Elicitation Package

## Understanding

**Source request:**
> While an import job is paused, the import service shall not process new records for that job.

**Requested outcome:** When an import job enters the paused state, the import service must suspend processing of new records belonging to that job until the paused state ends.

**Proposed system boundary:** The import service — the component responsible for consuming and processing records associated with an import job.

**Information sort:**

| Kind | Content |
|---|---|
| Behavioral obligation | The import service shall not process new records for a paused job — a state-driven inhibition. |
| Domain terms requiring definition | "import job", "paused" (a job lifecycle state), "process", "new records" (records not yet processed vs. newly arriving records — see Q1). |
| Assumption | The source uses "shall", a normative keyword, directly. Treated as an intended normative obligation, not a design idea or goal. |
| Assumptions | "New records" is interpreted provisionally as records not yet processed for the job; see Q1 for the distinction. |
| Design idea | None present; no internal mechanism is prescribed. |
| Questions | See Clarifying Questions below. |

**Pattern selection:** The source obligation holds continuously throughout the paused state. This is not a discrete event and does not depend on an optional product capability or an unwanted condition. The **State-driven** pattern (`While <state>, the <system> shall <response>`) is the correct choice.

---

## Clarifying Questions

**Q1 — Definition of "new records"**
The source says "new records for that job." Two interpretations are plausible:

- **(a) Unprocessed records** — records already queued or staged for the job that have not yet been processed.
- **(b) Newly arriving records** — records that arrive (are submitted, ingested, or enqueued) after the job enters the paused state.

These interpretations have different observable behavior: under (a) the service must also suspend in-flight or queued work; under (b) it need only reject or defer records that arrive during the pause. Which interpretation is intended, or does the obligation cover both?

**Q2 — Definition of "process"**
What observable actions constitute "processing" a record at the import service boundary? Plausible answers include: reading the record from a queue, validating it, transforming it, writing it to a destination, or issuing any downstream call on its behalf. The answer affects what a verifier would observe.

**Q3 — In-flight records at pause entry**
If the import service is actively processing a record at the moment the job transitions to paused, what is the required behavior?

- **(a)** Complete the in-progress record and then inhibit further processing.
- **(b)** Halt immediately (abort or roll back the in-progress record).
- **(c)** Undefined / outside scope of this requirement.

**Q4 — Definition and lifecycle of "paused"**
What determines whether an import job is in the paused state at the import service boundary? Specifically:

- What event or signal causes a job to enter the paused state (e.g., an operator command, an API call, an automatic policy)?
- What event or signal causes a job to leave the paused state?
- Is "paused" a distinct named state in a documented job lifecycle, or is it implied?

**Q5 — Scope: one job or all jobs**
The source says "for that job," which is read as: the inhibition applies only to the specific paused job, not to other concurrently running jobs. Please confirm this is correct.

**Q6 — Records deferred vs. discarded**
While a job is paused, should records that would otherwise be processed be:

- **(a)** Deferred and processed when the job resumes?
- **(b)** Discarded?
- **(c)** Returned to the caller with an error?
- **(d)** Undefined by this requirement?

---

## Candidate Requirements

### Candidate 1

- **Status:** needs clarification
- **Template:** State-driven
- **Requirement:** `While an import job is paused, the import service shall not process records for that job.`

  *(Note: "new records" retained below as a provisional marker; see Q1. The sentence above uses "records" as the unresolved term pending Q1's answer. A revision will restate the sentence once Q1–Q3 are resolved.)*

- **Source intent and rationale:** The source request expresses a state-driven inhibition: when an import job is in the paused state, the import service must not advance any records belonging to that job. The likely rationale is to give operators or the system the ability to suspend an import job without data being partially processed in an undefined state, preserving consistency and enabling safe resumption.

- **Implementability contract**
  - **System boundary:** The import service — the system component that reads, validates, transforms, or otherwise handles records belonging to an import job.
  - **Responsible system:** The import service.
  - **Actor:** Implicit — the import service acts autonomously on job records; no external actor triggers per-record processing explicitly in the source.
  - **Affected interfaces:** The interface at which the import service reads or accepts records for a job (e.g., a queue, a file endpoint, or an API endpoint). Internal to the import service boundary; the exact interface is **open** pending Q2.
  - **Inputs and data definitions:** A record belonging to an import job; the job identifier; the job's current lifecycle state. Definitions of "record," "job," and "paused state" are **open** pending Q1, Q2, and Q4.
  - **Preconditions, state, feature scope, and trigger:** Precondition: the import job is in the paused state. No additional trigger. How the service determines that a job is paused is **open** pending Q4.
  - **Observable normal, alternate, and failure response:** Normal: the import service does not advance (read, validate, transform, write) any record for the paused job while the paused state holds. Alternate: behavior for records already in flight at pause entry is **open** pending Q3. Behavior for records that arrive during the pause (deferred, rejected, discarded) is **open** pending Q6.
  - **Quantities, units, limits, timing, ordering, and tolerances:** No quantities stated. Timing of effect at state entry (immediate vs. after completing in-flight records) is **open** pending Q3.
  - **Domain terms, dependencies, and unresolved decisions:** "paused" (job lifecycle state — open, Q4); "process" (open, Q2); "new records" (open, Q1); in-flight behavior (open, Q3); disposition of deferred records (open, Q6).

- **Verification contract**
  - **Setup, state, input data, and stimulus:** Configure an import job with a known set of unprocessed records. Transition the job to the paused state. Attempt to trigger or observe processing of records for that job during the paused state.
  - **Observable evidence or oracle:** No record belonging to the paused job advances through any processing stage (as defined by Q2) while the job remains paused. Evidence is observable at the import service boundary — e.g., no output records are written, no downstream calls are issued, no dequeue events occur for that job's records.
  - **Expected result:** Zero records for the paused job are processed while the paused state holds.
  - **Pass/fail criteria, including scope, thresholds, units, and timing:** Pass: no processing actions (per Q2's definition) for the paused job's records occur during the paused state interval. Fail: any such action is observed. Scope is limited to the specific paused job (Q5 confirms this).
  - **Boundary, negative, and failure cases:** (a) Job is paused with zero pending records — inhibition is vacuously satisfied; no error expected. (b) Job is paused while a record is in flight — behavior **open** (Q3). (c) Other jobs' records continue to be processed normally — verifiable by checking that concurrent jobs are unaffected. (d) The paused state is immediately followed by resumption — records must be processed after resumption (covered by supporting requirement SR-2).
  - **Required data, instrumentation, assessment, or external evidence:** Ability to observe or query the processing state of individual records per job; access to job lifecycle state transitions; processing audit log or observable output at the service boundary.
  - **Observable at the named boundary alone:** Partially. Whether a record was "processed" may require internal instrumentation if the processing boundary is not exposed externally. Pending Q2.
  - **If no, internal evidence needed and rationale:** If "process" includes internal-only operations (e.g., validation before writing output), a processing audit log or internal instrumentation at the service's record-handling interface is needed to confirm the inhibition.

---

## Suggested Supporting Requirements

**SR-1 — State transition definition (Required companion)**
Relates to: Candidate 1 (precondition "import job is paused").
Missing behavior: What causes an import job to enter and exit the paused state is undefined in this requirement. Without it, the precondition of Candidate 1 is not implementable or verifiable.
Suggested form: `When [an authorized actor issues a pause command for an import job / a defined condition is met], the import service shall transition that job to the paused state.` (Pattern and actor pending Q4.)

**SR-2 — Resumption behavior (Required companion)**
Relates to: Candidate 1 (end of paused state).
Missing behavior: What happens to deferred records when the paused state ends is not addressed by Candidate 1. A complete picture of the pause feature requires that records are processed (or explicitly not processed) upon resumption.
Suggested form: `When a paused import job transitions to [running / active], the import service shall [resume processing records for that job / discard deferred records].` (Response pending Q6.)

**SR-3 — In-flight record behavior at pause entry (Failure-path companion)**
Relates to: Candidate 1 (timing of inhibition).
Missing behavior: The behavior of a record that is actively being processed at the moment the job transitions to paused is not defined. Two implementations (complete-then-stop vs. halt-immediately) are visibly different.
Suggested form: `While an import job is paused, if a record for that job is in progress at the time the job enters the paused state, the import service shall [complete processing that record / halt processing and [roll back / retain] partial state].` (Response pending Q3.)

**SR-4 — Concurrent-job isolation (Boundary companion)**
Relates to: Candidate 1 (scope "for that job").
Missing behavior: The requirement scopes inhibition to one job, implying other jobs are unaffected, but no requirement explicitly states that other jobs shall continue processing normally. An explicit statement avoids ambiguity.
Suggested form: `While an import job is paused, the import service shall continue to process records for other import jobs that are not paused.`

**SR-5 — Auditability / observability of paused state (Operational companion)**
Relates to: Candidate 1 (verifiability).
Missing behavior: The requirement does not state whether the import service must expose or record that processing was inhibited. An observable audit trail or status indicator supports both operation and verification.
Suggested form: *(optional consideration)* `While an import job is paused, the import service shall [record / make observable] the paused status of that job [at its status interface / in its audit log].` Ask whether auditability is in scope.

---

## Consistency Findings

**Finding 1 — Gap: pause lifecycle undefined**
Candidate 1's precondition depends on a "paused" state, but no requirement in the current set defines when a job enters or exits that state. Without SR-1, the set is incomplete and Candidate 1's precondition cannot be evaluated. *Resolution: accept SR-1 or provide an existing lifecycle requirement as a reference.*

**Finding 2 — Gap: resumption behavior undefined**
Candidate 1 prohibits processing during the paused state but does not define what happens to affected records when the state ends. SR-2 addresses this gap. Without it, two implementations — one that discards deferred records and one that processes them — both satisfy Candidate 1. *Resolution: accept SR-2 or explicitly mark out-of-scope.*

**Finding 3 — Gap: in-flight behavior at state entry**
Candidate 1 does not address behavior at the boundary between the running and paused states. SR-3 addresses this gap. The gap is material because two implementations (complete-and-stop vs. immediate-halt) are visibly different at the service boundary. *Resolution: accept SR-3 or explicitly mark out-of-scope.*

**Finding 4 — Potential overlap with a broader inhibition requirement**
If the system contains a general requirement such as "a paused job shall not perform any operation," Candidate 1 may be a special case of or duplicate to that broader requirement. No such requirement is present in the current set; flag if one exists in the host system.

---

## Assumptions

**A1** — "Shall" in the source is treated as a normative obligation keyword, not as source prose using "shall" informally. If "shall" was used colloquially (meaning "will" or "is expected to"), the normative force of Candidate 1 should be confirmed.

**A2** — "New records" is interpreted provisionally as records that have not yet been processed at the time of observation (i.e., pending records), not strictly as records that arrive after the pause begins. This is unresolved pending Q1.

**A3** — "For that job" is interpreted as scoping the inhibition to only the specific paused import job; other concurrently running jobs are unaffected.

**A4** — "The import service" is a single named system responsible for processing records. If "import service" is a label for multiple components or microservices, the system boundary and responsible component must be identified before Candidate 1 is implementable.

---

## Rationale and Context

The source request describes a classical state-driven inhibition pattern: a named system must suppress a continuous behavior (processing records) while an enduring condition (the job is paused) holds. This is a well-understood safety-through-inhibition pattern used in batch and streaming data pipelines.

The primary elicitation risk is the undefined lifecycle of the paused state: if the entry and exit conditions for "paused" are not specified, the requirement cannot be verified in isolation. The secondary risk is the ambiguity in "new records," which affects the scope of the inhibition and the disposition of in-flight work.

The requirement as stated is structurally clean and maps directly to the State-driven EARS pattern. Its status is `needs clarification` because the implementability and verification contracts both contain open items (Q1–Q6) that would cause two reasonable implementers or verifiers to choose different behavior.

---

## Review Status

| Item | Status |
|---|---|
| Candidate 1 | **needs clarification** — open items: Q1 (definition of "new records"), Q2 (definition of "process"), Q3 (in-flight behavior at pause entry), Q4 (definition and lifecycle of "paused"), Q6 (disposition of deferred records). |
| SR-1 (pause lifecycle) | Suggested — not yet a candidate. Required for Candidate 1's precondition to be implementable. |
| SR-2 (resumption behavior) | Suggested — not yet a candidate. Required for completeness. |
| SR-3 (in-flight at pause entry) | Suggested — not yet a candidate. Required for completeness. |
| SR-4 (concurrent-job isolation) | Suggested — not yet a candidate. Boundary clarification. |
| SR-5 (auditability) | Suggested as optional consideration. Pending scope decision. |
| Consistency findings | Three gaps identified; no conflicts or duplicates within the current candidate set. |
| Unresolved questions | Q1, Q2, Q3, Q4, Q5 (low-risk, provisionally answered by A3), Q6. |
