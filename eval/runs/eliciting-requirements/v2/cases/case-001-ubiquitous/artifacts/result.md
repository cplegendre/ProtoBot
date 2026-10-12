# Requirements Elicitation Package

## Understanding

**Requested outcome:** The audit service must keep every accepted audit event for 30 days and allow an authorized auditor to retrieve any such event during that period.

**Proposed system boundary:** The audit service — the named system that accepts, stores, and exposes audit events to external actors.

**Actors identified:**
- **Audit service** — the responsible system.
- **Authorized auditor** — an external actor who retrieves stored events.
- **Event submitter** — the producer that submits audit events (implied; not described in source).

**Domain terms requiring definition:**
- *Accepted audit event* — the source implies a lifecycle state ("accepted") but does not define what acceptance means, what data an event contains, or how acceptance is signaled. Assumed to mean: an event that the audit service has received and positively validated or acknowledged.
- *Authorized auditor* — the source does not define how authorization is established, granted, or verified. The mechanism is left open.
- *Retrieve* — the source does not specify how retrieval is requested (query by ID, by time range, by attribute, etc.) or what the response contains.
- *30-day retention period* — the source does not specify whether the 30-day clock starts at event acceptance time or at some other reference point, or what happens at the boundary (wall-clock days, rolling 24-hour periods, etc.).

**Obligations identified from the source:**
1. The audit service shall retain every accepted audit event for 30 days.
2. An authorized auditor shall be able to retrieve an accepted event during the retention period.

Both obligations are unconditional in the source — they apply without any qualifying state, trigger, optional feature, or failure condition stated. Each is a candidate for the **Ubiquitous** EARS pattern; however, both have unresolved domain terms in their contracts.

---

## Clarifying Questions

**Q1 — Definition of "accepted":**
What does it mean for an audit event to be "accepted"? Is acceptance a validation step, an acknowledgment returned to the submitter, a write to durable storage, or some other observable outcome? The answer determines when the 30-day clock starts and what data must be retained.

**Q2 — Retention clock reference point:**
Does the 30-day retention period begin at the moment the audit service accepts the event (i.e., per-event timestamp), or does it apply as a rolling window from the current time? And are "days" defined as 24-hour periods (UTC or local time zone) or calendar days?

**Q3 — Definition of "authorized auditor":**
How is an authorized auditor distinguished from an unauthorized caller? Is authorization enforced by the audit service itself (e.g., checking a role or credential), or is it enforced by an upstream gateway? The answer determines whether the audit service has a direct authorization obligation.

**Q4 — Retrieval interface:**
How does an authorized auditor retrieve an event — by event identifier, by time range, by attribute filter, or some other means? What is the expected response (full event payload, summary, specific fields)? Without this, two implementers would design incompatible interfaces.

**Q5 — Behavior at retention boundary:**
What shall the audit service do with an accepted event when the 30-day retention period expires — delete it, archive it elsewhere, anonymize it? Is there a requirement on what is observable at day 30 boundary (inclusive or exclusive)?

**Q6 — Behavior when retrieval is requested after expiry:**
If an authorized auditor requests an event that has passed its retention period, what shall the audit service return — an error, an empty result, or a specific status? This is a failure-path question that two implementers would likely resolve differently without a requirement.

**Q7 — Scope of "every":**
Does "every accepted audit event" include events accepted before the audit service was updated to this 30-day retention policy, or only events accepted after the policy takes effect?

---

## Candidate Requirements

### Candidate 1

- **Status:** needs clarification
- **Template:** Ubiquitous
- **Requirement:** The audit service shall retain every accepted audit event for 30 days from the time of acceptance.

- **Source intent and rationale:** The source states that the audit service "must retain every accepted audit event for 30 days." This is an unconditional obligation on the audit service. The 30-day duration is explicit. "From the time of acceptance" is an interpretation of when the clock starts (see Q1, Q2); it is the most common meaning but is not confirmed.

- **Implementability contract**
  - System boundary: The audit service — the system that receives and stores audit events.
  - Responsible system: The audit service.
  - Actor: Any event submitter whose event is accepted (the retention obligation is on the service, not the actor).
  - Affected interfaces: Durable storage interface of the audit service; any interface that could delete or expunge events before 30 days must be constrained.
  - Inputs and data definitions: *Accepted audit event* — **open; see Q1.** The event identity and its accepted-at timestamp must be defined.
  - Preconditions, state, feature scope, and trigger: None stated. The obligation is unconditional once an event is accepted.
  - Observable normal, alternate, and failure response: An accepted event remains retrievable for 30 days. Behavior at or after expiry is **open; see Q5.**
  - Quantities, units, limits, timing, ordering, and tolerances: 30 days — unit (calendar days vs. 24-hour periods) and reference timezone are **open; see Q2.**
  - Domain terms, dependencies, and unresolved decisions: *Accepted*, *30 days* — both open; see Q1, Q2.

- **Verification contract**
  - Setup, state, input data, and stimulus: Submit an event that the audit service accepts. Record the acceptance timestamp.
  - Observable evidence or oracle: The event must be retrievable at 1 day, at 15 days, and at 29 days after acceptance.
  - Expected result: The audit service returns the retained event at each retrieval attempt within the 30-day window.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Event is present and retrievable for the full 30-day retention period. The exact day boundary (inclusive/exclusive) is **open; see Q2.**
  - Boundary, negative, and failure cases: Retrieval attempt at exactly 30 days (boundary); retrieval attempt after 30 days (see Q5, Q6). Both are **open.**
  - Required data, instrumentation, assessment, or external evidence: Access to the audit service retrieval interface; a mechanism to observe the accepted-at timestamp and confirm event presence/absence.
  - Observable at the named boundary alone: yes — provided a retrieval interface exists (see Q4; that interface is currently unspecified).
  - If no, internal evidence needed and rationale: N/A pending Q4 resolution.

---

### Candidate 2

- **Status:** needs clarification
- **Template:** Ubiquitous
- **Requirement:** The audit service shall return an accepted audit event to an authorized auditor who requests it within the event's retention period.

- **Source intent and rationale:** The source states an "authorized auditor must be able to retrieve an accepted event during that retention period." This is an unconditional obligation on the audit service: when a valid retrieval request is made within the retention window, the service must fulfill it. "Authorized auditor" and the retrieval interface are both undefined (see Q3, Q4).

- **Implementability contract**
  - System boundary: The audit service — the system that exposes stored events to external actors.
  - Responsible system: The audit service.
  - Actor: Authorized auditor — **open; authorization mechanism undefined; see Q3.**
  - Affected interfaces: The retrieval interface of the audit service — **open; see Q4.**
  - Inputs and data definitions: The retrieval request (parameters, format, credentials) is **open; see Q4.** Event identity and accepted-at timestamp required (see Candidate 1).
  - Preconditions, state, feature scope, and trigger: The request must be made within the 30-day retention period. Authorization must be established. Both conditions are implied by the source but their mechanisms are open.
  - Observable normal, alternate, and failure response: Normal — the audit service returns the requested event. Unauthorized access response is **open; see Q3.** Post-expiry response is **open; see Q6.**
  - Quantities, units, limits, timing, ordering, and tolerances: Within 30-day retention period — inherits the open unit question from Candidate 1 (Q2). No response-time or latency requirement is stated.
  - Domain terms, dependencies, and unresolved decisions: *Authorized auditor*, *retrieve*, *retention period* — all open; see Q1–Q4.

- **Verification contract**
  - Setup, state, input data, and stimulus: An accepted event is present within its retention period. An authorized auditor submits a valid retrieval request for that event.
  - Observable evidence or oracle: The audit service returns the requested event's data.
  - Expected result: The authorized auditor receives the accepted event.
  - Pass/fail criteria, including scope, thresholds, units, and timing: Retrieval succeeds for any accepted event within its 30-day retention window. No latency threshold specified; add one if required.
  - Boundary, negative, and failure cases: Retrieval by an unauthorized caller (see Q3); retrieval of a non-existent event; retrieval after retention expiry (see Q6). All open.
  - Required data, instrumentation, assessment, or external evidence: A defined retrieval interface (Q4); a means to establish and verify auditor authorization (Q3); event presence/absence oracle (linked to Candidate 1).
  - Observable at the named boundary alone: yes — provided retrieval interface and authorization mechanism are defined.
  - If no, internal evidence needed and rationale: N/A pending Q3, Q4 resolution.

---

## Suggested Supporting Requirements

**S1 — Failure-path companion (to Candidate 2):**
*Missing behavior:* What the audit service shall do when an unauthorized caller requests an accepted event.
*Suggestion:* If a retrieval request is made by a caller who is not an authorized auditor, then the audit service shall deny the request without disclosing any retained event data.
*Reason:* Without this, two implementers would choose different denial behaviors (silent empty response, explicit error, disclosure of metadata). Authorization failure is a material failure path that affects security posture.

**S2 — Failure-path companion (to Candidates 1 and 2):**
*Missing behavior:* What the audit service shall do when a retrieval request is made for an event whose retention period has expired.
*Suggestion:* If an authorized auditor requests an accepted audit event whose 30-day retention period has expired, then the audit service shall return a response indicating the event is no longer available.
*Reason:* The source defines the retrieval right only within the retention period; behavior after expiry is unspecified and two implementers would produce different observable results.

**S3 — Boundary companion (to Candidate 1):**
*Missing behavior:* What the audit service shall do with a retained event when its 30-day retention period expires.
*Suggestion:* When the 30-day retention period for an accepted audit event expires, the audit service shall [delete / archive / anonymize] the event.
*Reason:* The source requires 30-day retention but does not state the disposition action. Without this, implementers may silently accumulate events indefinitely or purge them inconsistently.
*Note: Disposition action is bracketed because it is a policy decision for the user to confirm.*

**S4 — Required companion (to Candidate 2):**
*Missing behavior:* Definition of how the audit service determines whether a caller is an authorized auditor.
*Suggestion:* The audit service shall verify that a retrieval request originates from an authorized auditor before returning any retained event data.
*Reason:* "Authorized auditor" appears in the source but no authorization mechanism is defined. Without this companion, the obligation in Candidate 2 cannot be implemented or verified at the audit service boundary.

**S5 — Interface companion (to Candidate 2):**
*Missing behavior:* The retrieval interface contract — what parameters identify an event, what the response contains.
*Suggestion:* Define the retrieval request parameters (e.g., event identifier, time range) and the required response content (e.g., full event payload, specific fields, format).
*Reason:* Without a defined interface, the retrieval obligation in Candidate 2 cannot be implemented consistently or verified.

**S6 — Optional consideration:**
*Missing behavior:* Whether the audit service must guarantee that retained events are immutable (tamper-evident) during the 30-day period.
*Suggestion:* Consider whether the audit service shall detect or prevent modification of an accepted event during its retention period, and what evidence it shall provide on retrieval.
*Reason:* Audit event integrity is a common implicit expectation in audit services, but it is not stated in the source. Confirm whether it is in scope.

---

## Consistency Findings

**Finding 1 — Gap: retention expiry and retrieval right overlap undefined**
Candidates 1 and 2 together imply that retrieval is available for 30 days. However, neither candidate defines the observable boundary behavior at exactly day 30 (inclusive vs. exclusive). Two verifiers would apply different pass/fail criteria at the boundary. Resolution: confirm whether day 30 is included in the retrieval window, and express the boundary explicitly in both candidates.

**Finding 2 — Gap: no failure path for unauthorized access**
Candidate 2 defines the retrieval obligation for an authorized auditor, but no candidate defines what the audit service shall do when an unauthorized caller attempts retrieval. S1 (suggested) covers this path. Without it, the candidate set is incomplete for authorization.

**Finding 3 — Gap: post-expiry retrieval behavior**
Candidate 1 requires retention for 30 days; Candidate 2 requires retrieval within the retention period. Neither candidate states what the service returns when retrieval is requested after the window closes. S2 and S3 (suggested) address this. Until resolved, the candidate set has an unspecified failure path that two implementers would resolve differently.

**Finding 4 — Potential overlap: Candidate 1 and Candidate 2**
Candidate 1 establishes that retained events exist; Candidate 2 establishes that they are accessible. These are distinct obligations (existence vs. access) and are correctly separated. No overlap or conflict, but the linkage should be preserved: Candidate 2 depends on Candidate 1 being satisfied.

No conflicts found between the two candidates.

---

## Assumptions

**A1:** "Accepted audit event" means an event that the audit service has received and positively acknowledged or validated — i.e., the service has made the event durable. This is provisional; see Q1.

**A2:** The 30-day retention clock starts at the moment of acceptance for each individual event, measured as 30 consecutive 24-hour periods in UTC. This is provisional; see Q2.

**A3:** Authorization of an auditor is a property that the audit service can observe or verify at its boundary (e.g., via a credential, role, or token). This is provisional; see Q3.

**A4:** "Retrieve" implies that the authorized auditor submits a request to the audit service and receives the event data in response. The specific request/response protocol is not assumed. See Q4.

---

## Rationale and Context

The source expresses two related obligations from what appears to be a compliance or governance context: (1) event preservation for a fixed window, and (2) access for a privileged actor during that window. These are common in audit trail requirements regulated by security, privacy, or financial compliance frameworks.

The two obligations are separated into distinct candidates because they can be negotiated and verified independently: an audit service could satisfy the retention obligation while failing the retrieval obligation (e.g., events are retained but no retrieval API exists), or vice versa.

The Ubiquitous pattern is the correct candidate for both obligations because neither the source statement nor the elicited context introduces a limiting state, a discrete trigger, an optional feature scope, or an unwanted condition. Both obligations are always active once an event is accepted.

Both candidates remain at `needs clarification` because material domain terms (*accepted*, *authorized auditor*, *retrieve*, *30 days* boundary) have no confirmed definition, and at least one reasonable implementer or verifier would make a materially different choice without that information.

---

## Review Status

| Candidate | Status | Blocking item(s) |
|-----------|--------|-----------------|
| Candidate 1 — Retention | needs clarification | Q1 (accepted), Q2 (30-day unit/boundary), Q5 (expiry behavior) |
| Candidate 2 — Retrieval | needs clarification | Q1 (accepted), Q2 (retention period boundary), Q3 (authorization), Q4 (retrieval interface), Q6 (post-expiry response) |

**Suggested supporting requirements:** S1–S6 are outside the candidate list and await user acceptance.

**Next step:** Answer Q1–Q7 to allow both contracts to be completed and candidates to advance to `candidate` or `ready for review`.
