# Decision Log

Append-only. Format: context → decision → rationale → consequences.
Newest entries go at the bottom.

---

## D1 — Go as the implementation language (2026-08-22)

**Context:** Single-developer personal project; needs a small deployable HTTP
service talking to external APIs.

**Decision:** Write the service in Go.

**Rationale:** Strong stdlib HTTP support (`net/http`) — no web framework
needed; static binaries simplify deployment; first-class concurrency if async
processing is needed later; good fit for learning idiomatic Go (interfaces,
table-driven tests).

**Consequences:** No framework magic; standard project layout under `cmd/` +
`internal/`; tests use stdlib `testing` + table-driven style.

## D2 — WhatsApp Business API as the sole input source (2026-08-22)

**Context:** Expenses should be loggable instantly from a phone with zero UI
to build.

**Decision:** Users log expenses by messaging a WhatsApp Business number;
no separate app/web UI for entry.

**Rationale:** Zero-friction capture where the user already is; messaging UX
is familiar; avoids building/maintaining any frontend in v1.

**Consequences:** Must handle WhatsApp-specific concerns: webhook delivery,
verification handshake, signature checks, retry-based duplicates, and
eventually non-text messages (images/receipts) are out of scope for v1.

## D3 — Google Sheets as the persistence layer (2026-08-22)

**Context:** Storage must be visible/editable by a non-technical user without
extra tooling.

**Decision:** Store expenses as rows in a Google Sheet, accessed via service
account auth.

**Rationale:** The user gets a live spreadsheet view "for free"; no database
to host; Sheets API supports append and basic reads.

**Consequences:** Sheets API has quotas and higher latency than a DB — keep
writes single-row appends; complex reporting may be limited; wrap access in a
`Store` interface so a DB can replace it later without touching business
logic.

## D4 — Webhook (push) integration pattern, not polling (2026-08-22)

**Context:** Provider delivers new messages either by pushing to our endpoint
or by us repeatedly asking for updates.

**Decision:** Use webhooks — provider pushes each message to our HTTPS
endpoint as it arrives.

**Rationale:** Near real-time logging; no wasted polling requests; simpler
state model (no "last seen timestamp" bookkeeping).

**Consequences:** Service must be publicly reachable over HTTPS; must answer
the verification handshake; must verify `X-Hub-Signature-256`; must tolerate
provider retries idempotently (see `architecture.md`).

## D5 — Regex-first parsing philosophy (2026-08-22)

**Context:** Message text must become structured data (amount, category,
description). Options range from rigid patterns to LLM extraction.

**Decision:** Start with deterministic rule-based parsing (patterns → struct);
keep the parser behind an interface so an LLM/NLP implementation can slot in
later if needed.

**Rationale:** Deterministic, fast, free, fully unit-testable; failures are
predictable and debuggable; avoids external API dependency and cost in v1;
interface boundary preserves the upgrade path.

**Consequences:** v1 accepts only well-formed commands; handling of
unrecognized input is defined by D6.

## D6 — Slash-command message protocol (2026-08-22)

**Context:** Meta delivers every message sent to the business number to our
webhook — filtering cannot happen on WhatsApp's side. We need deterministic
parsing and must avoid accidentally logging casual chatter aimed at the
number.

**Decision:** Process only messages starting with a recognized slash command;
silently ignore everything else (no error replies).

v1 grammar:

```
/spend [today|yesterday] <amount> <merchant> <person-tag>
/earned <amount> <method>       e.g. /earned 149 online · /earned 23 cash
```

- `person` is a tag recording who/what an expense is for (e.g. "sapna")
- Amounts carry no currency token; one fixed currency lives in service config
- Date defaults to today; `yesterday` keyword overrides

**Rationale:** The command prefix provides an unambiguous processing gate and
a clean dispatch point (switch on command); a deterministic grammar fits D5's
regex-first approach; silent-ignore avoids replying to strangers/spam hitting
a public number.

**Consequences:** Typos get no feedback (accepted tradeoff); future commands
(`/report`, `/undo`, …) slot into the same router; grammar refinements
(multi-word merchants, person optionality, date vocabulary) tracked in Q5.

## D7 — Credit/debit ledger sheet layout (2026-08-22)

**Context:** Scope expanded beyond expenses — income is logged via `/earned`,
and balance visibility is wanted at week's end.

**Decision:** Store transactions in ledger form with distinct **debit** and
**credit** columns; exactly one is populated per row. Spends are debits;
earnings are credits. Weekly balance must be available end-of-week.

Working row schema:

| date | description | debit | credit | person | method | raw_message | wamid |
|---|---|---|---|---|---|---|---|

**Rationale:** Matches spreadsheet bookkeeping conventions; `SUM(credit) −
SUM(debit)` works natively in Sheets; human-scannable without extra tooling.

**Consequences:** The parser emits one unified `Transaction{Direction
debit|credit, ...}` rather than separate expense/income types; a `method`
column exists for earned entries (whether spends also record method: Q7);
balance mechanism (sheet formulas vs `/balance` command vs scheduled rows)
is deferred to Q6.
