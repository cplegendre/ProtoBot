# Requirements Elicitation Package

## Understanding

**Source request:** "The UI must be responsive."

**Requested outcome:** The user interface must perform or present itself in a way that qualifies as "responsive." This is ambiguous: "responsive" in UI contexts has two common, distinct meanings:

1. **Layout responsiveness** — the UI adapts its layout and presentation to different screen sizes, viewport widths, or device classes (e.g., mobile, tablet, desktop).
2. **Interaction responsiveness** — the UI reacts to user input within a perceptible time threshold (e.g., a button press produces visible feedback within N milliseconds).

Both interpretations are legitimate. Neither can be confirmed from the source alone.

**Proposed system boundary:** The user interface layer of an unspecified system. The system name, the technology stack, the interaction modalities (touch, pointer, keyboard), and the deployment targets are all unspecified.

**Status of the source statement:** The term "responsive" is ambiguous, no system is named, no measurable criterion is given, and no actor or device scope is established. The source cannot be implemented or verified by two independent teams in the same way. Both candidate requirements below are therefore `needs clarification`.

---

## Clarifying Questions

The following questions are listed in priority order. Questions 1 and 2 are blocking — without answers, no normative requirement can be drafted.

1. **Which meaning of "responsive" is intended?**
   - (a) Layout adaptation to different screen sizes or device classes, or
   - (b) Interaction latency — the UI reacts to user input within a time limit, or
   - (c) Both?

2. **What is the name of the system or component whose UI is being constrained?**
   (e.g., "the Patient Portal," "the Admin Dashboard," "the mobile app")

3. **If layout responsiveness (a) is in scope:** What are the target viewport widths, device classes, or breakpoints? (e.g., 320 px–480 px mobile, 768 px–1024 px tablet, 1280 px+ desktop.) Are there specific orientations (portrait/landscape) to support?

4. **If interaction responsiveness (b) is in scope:** What is the maximum acceptable delay, in milliseconds, between a user action (e.g., tap, click, keystroke) and a visible acknowledgement? Is the limit the same for all interaction types, or does it vary by action (e.g., button feedback vs. data load)?

5. **What is the verifiable oracle?** How will a reviewer determine whether the UI is responsive — manual inspection, automated test, a profiling tool, a specific measurement device, or a combination?

6. **Are there conditions under which the requirement does not apply?** For example, during a loading state, a background data fetch, or an error state?

---

## Candidate Requirements

### Candidate 1 — Layout responsiveness

- **Status:** needs clarification
- **Template:** unresolved

**Requirement:** Drafting is blocked. The system name, the set of target viewport widths or device classes, the observable layout behavior at each breakpoint, and the pass/fail criterion are all unknown. A provisional sentence form would be:

> While the \<system name\> UI is displayed at \<unresolved viewport range\>, the \<system name\> shall present its content without \<unresolved failure criterion, e.g., horizontal overflow, content clipping, or loss of functional control\>.

This sentence cannot be finalized until questions 1, 2, 3, and 5 are answered.

- **Source intent and rationale:** The requester wants the UI to look and function correctly across the range of devices or screen sizes used by its target audience. This is a presentation-layer quality obligation.

- **Implementability contract**
  - System boundary: unresolved — no system name supplied
  - Responsible system: unresolved
  - Actor: end user viewing or interacting with the UI on an unspecified device
  - Affected interfaces: the visual rendering surface (browser viewport, native screen, etc.) — unresolved
  - Inputs and data definitions: viewport width, height, orientation — all unresolved
  - Preconditions, state, feature scope, and trigger: the UI is displayed; specific viewport or device class — unresolved
  - Observable normal, alternate, and failure response: unresolved — "responsive" has no agreed measure
  - Quantities, units, limits, timing, ordering, and tolerances: breakpoints, pixel ranges — unresolved
  - Domain terms, dependencies, and unresolved decisions: "responsive," "breakpoint," target device set — all unresolved; see questions 1–5

- **Verification contract**
  - Setup, state, input data, and stimulus: unresolved — device or viewport to test, content state to load, and user action to perform are all unspecified
  - Observable evidence or oracle: unresolved — no measurement method specified
  - Expected result: unresolved
  - Pass/fail criteria, including scope, thresholds, units, and timing: unresolved
  - Boundary, negative, and failure cases: unresolved — minimum and maximum viewport, orientation change, font-scale change not defined
  - Required data, instrumentation, assessment, or external evidence: unresolved
  - Observable at the named boundary alone: unresolved
  - If no, internal evidence needed and rationale: unresolved

---

### Candidate 2 — Interaction responsiveness

- **Status:** needs clarification
- **Template:** unresolved

**Requirement:** Drafting is blocked. The system name, the triggering user action, the maximum permissible delay, and the observable acknowledgement are all unknown. A provisional sentence form would be:

> When a user performs \<unresolved action, e.g., activates a control\>, the \<system name\> shall display a visible acknowledgement within \<unresolved N\> milliseconds.

This sentence cannot be finalized until questions 1, 2, 4, and 5 are answered.

- **Source intent and rationale:** The requester wants users to perceive that the system reacts promptly to their actions, reducing the perception of lag or unresponsiveness.

- **Implementability contract**
  - System boundary: unresolved — no system name supplied
  - Responsible system: unresolved
  - Actor: end user interacting with the UI
  - Affected interfaces: input device (pointer, touch, keyboard) and display surface — unresolved
  - Inputs and data definitions: user action type — unresolved
  - Preconditions, state, feature scope, and trigger: user performs a discrete action — trigger type unresolved
  - Observable normal, alternate, and failure response: visible acknowledgement — definition and form unresolved
  - Quantities, units, limits, timing, ordering, and tolerances: maximum delay in milliseconds — unresolved
  - Domain terms, dependencies, and unresolved decisions: "responsive," "acknowledgement," "action" — unresolved; see questions 1, 2, 4, 5

- **Verification contract**
  - Setup, state, input data, and stimulus: unresolved — action type and device not specified
  - Observable evidence or oracle: unresolved — no timing method or tool specified
  - Expected result: unresolved — no time threshold given
  - Pass/fail criteria, including scope, thresholds, units, and timing: unresolved
  - Boundary, negative, and failure cases: behavior when the system is under load or in a degraded state — unresolved
  - Required data, instrumentation, assessment, or external evidence: unresolved — profiling tool, automated test, or manual observation not chosen
  - Observable at the named boundary alone: unresolved
  - If no, internal evidence needed and rationale: unresolved

---

## Suggested Supporting Requirements

The following are suggestions only. None is an approved requirement. Each is listed with its relationship to the source, the missing behavior it addresses, and a short reason.

1. **Relationship:** Required companion (layout interpretation)
   **Missing behavior:** Behavior at orientation change — when a user rotates the device, the UI shall reflow its content without loss of data or functional controls.
   **Reason:** Layout responsiveness is often tested at a fixed viewport; rotation introduces a dynamic viewport change that must be handled separately to avoid a verification gap.

2. **Relationship:** Required companion (interaction interpretation)
   **Missing behavior:** Loading feedback — when a user action triggers a background operation that cannot complete within the interaction-response threshold, the UI shall display a progress indicator until the operation completes or fails.
   **Reason:** Without this, the interaction-responsiveness requirement has no defined behavior for slow operations; a verifier cannot distinguish an unresponsive UI from a correctly loading one.

3. **Relationship:** Failure-path companion (interaction interpretation)
   **Missing behavior:** If the UI cannot render a response within the agreed maximum delay, the UI shall display an error or degraded-state indicator rather than a blank or frozen display.
   **Reason:** The source gives no failure-mode behavior; two implementers would choose different behaviors, making verification impossible.

4. **Relationship:** Boundary companion (layout interpretation)
   **Missing behavior:** Minimum content legibility — when the UI is displayed at the minimum supported viewport width, all text shall remain legible and all interactive controls shall remain operable.
   **Reason:** "Responsive layout" implies minimum usability guarantees at the smallest supported size, which the source does not specify.

5. **Relationship:** Optional consideration
   **Missing behavior:** Accessibility at variable text scale — when the user sets a system font scale above the default, the UI shall not clip or overlap text.
   **Reason:** Responsive layout and accessibility overlap when users increase text size; this may be in scope but the source does not confirm it. Ask whether it belongs.

---

## Consistency Findings

**Finding 1 — Overlap risk between Candidate 1 and Candidate 2**
- **Kind:** Overlap
- **Affected candidates:** Candidate 1 (layout), Candidate 2 (interaction)
- **Condition:** If the requester intends both meanings of "responsive," the two candidates address independent obligations and must be maintained as separate requirements. If only one meaning is intended, the unused candidate should be withdrawn to avoid a spurious requirement.
- **Resolution options:** Answer question 1 to confirm the intended scope; then either retain both or withdraw the inapplicable one.

**Finding 2 — Gap: no failure or degradation behavior**
- **Kind:** Gap
- **Condition:** Neither candidate defines what the UI must do when it cannot meet the responsiveness criterion (e.g., under extreme load, on an unsupported device, or during a network failure).
- **Resolution options:** Add a failure-path companion (see Suggested Supporting Requirements 2 and 3) once the primary candidates are confirmed.

**Finding 3 — Vacuous risk**
- **Kind:** Infeasible or vacuous
- **Condition:** Without a measurable threshold or defined device scope, either candidate could be satisfied by any UI implementation or could be interpreted as requiring infinite performance. The requirement is currently unverifiable.
- **Resolution options:** Resolve questions 3–5 to introduce a measurable criterion before proceeding.

---

## Assumptions

All items below are provisional. They are not requirements.

1. The word "must" in the source is treated as expressing a normative obligation (`shall` in EARS), not a preference or aspiration.
2. "UI" refers to a user-facing interface rendered on a screen. It does not refer to a command-line interface or a voice interface unless the user confirms otherwise.
3. The requirement applies to the normal operating state of the UI; it is not assumed to apply during installation, initial loading before the first render, or explicit error or maintenance states unless confirmed.

---

## Rationale and Context

The source statement "The UI must be responsive" is a widely used but highly ambiguous quality aspiration. It cannot be implemented or verified in its current form because:

- "Responsive" conflates two distinct and independently verifiable obligations (layout adaptation and interaction latency).
- No measurable criterion is given for either meaning.
- No system name, device scope, or actor is identified.
- No failure behavior is described.

The reference for this class of source statement notes explicitly: *"The UI must be responsive. Ask for interaction, workload, device scope, response-time measure, and pass/fail observation."*

Neither candidate can advance past `needs clarification` until at least questions 1, 2, and either 3 or 4 (depending on the meaning) are answered.

---

## Review Status

| Candidate | Status | Blocking items |
|---|---|---|
| Candidate 1 — Layout responsiveness | needs clarification | Questions 1, 2, 3, 5 unanswered; measurable criterion and system name absent; implementability and verification contracts incomplete |
| Candidate 2 — Interaction responsiveness | needs clarification | Questions 1, 2, 4, 5 unanswered; time threshold and system name absent; implementability and verification contracts incomplete |

No candidate is `ready for review`. The package should be returned to the requester with the clarifying questions above. Once answers are received, the elicitation loop should resume and both contracts re-evaluated before advancing either candidate.
