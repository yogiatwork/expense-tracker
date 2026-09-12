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

> **⚠ Superseded by D8** — input source switched to Telegram Bot API.
> This decision is preserved for historical context only.

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

**Consequences:** Service must be publicly reachable over HTTPS; webhook URL
registered with Telegram via `setWebhook`; must verify
`X-Telegram-Bot-Api-Secret-Token`; must tolerate provider retries
idempotently (see `architecture.md`).

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

**Context:** Telegram delivers every message sent to the bot to our webhook
— filtering cannot happen at the bot layer. We need deterministic parsing
and must avoid accidentally logging casual chatter aimed at the bot.

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
a public bot.

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

| date | description | debit | credit | person | method | raw_message | update_id |
|---|---|---|---|---|---|---|---|

**Rationale:** Matches spreadsheet bookkeeping conventions; `SUM(credit) −
SUM(debit)` works natively in Sheets; human-scannable without extra tooling.

**Consequences:** The parser emits one unified `Transaction{Direction
debit|credit, ...}` rather than separate expense/income types; a `method`
column exists for earned entries (whether spends also record method: Q7);
balance mechanism (sheet formulas vs `/balance` command vs scheduled rows)
is deferred to Q6.

## D8 — Telegram Bot API replaces WhatsApp Business API (2026-08-22)

**Context:** WhatsApp Business API setup proved painful (OTP rate limits,
business verification, test number restrictions). Telegram Bot API offers
identical functionality for our use case with instant bot creation, no
business verification, and simpler webhook mechanics.

**Decision:** Input source changed from WhatsApp Business API to Telegram Bot
API. This supersedes D2.

**Rationale:** Instant bot creation via BotFather (no review process);
webhook registration via a single API call (`setWebhook`); simpler signature
verification (plain header compare vs HMAC-SHA256); long-polling available
for dev without a tunnel; generous free tier; identical UX from the user's
perspective (just send commands to a bot).

**Consequences:** Architecture docs updated; D2 marked superseded; row schema
`wamid` column renamed to `update_id`; WhatsApp-specific components (Verifier,
`hub.challenge` handshake) removed from the design.

## D9 — Gin web framework (2026-08-22)

**Context:** D1 assumed stdlib `net/http` would suffice. During scaffolding,
the developer chose Gin based on prior familiarity.

**Decision:** Use Gin as the HTTP framework for the webhook server.

**Rationale:** Already familiar from prior projects — lowest friction for
moving fast. Built-in JSON binding, structured logging middleware, and
graceful shutdown are useful ergonomics for a webhook handler.

**Consequences:** Supersedes D1's "stdlib alone" consequence; raw body must
be read via `c.GetRawData()` before Gin's binding touches it (important for
signature verification); module dependency added.

## D10 — Read/reporting command protocol (2026-08-22)

**Context:** Users must be able to query what they've spent and their month-end
balance, not just log transactions. This resolves Q6 (balance mechanism) in
favour of on-demand commands over scheduled/sheet-formula approaches.

**Decision:** Three read commands, month-scoped by default:

```
/balance    [YYYY-MM]            net = earned − spent for the month
/spendtotal [YYYY-MM] [person]   total spend for the month, optionally per person tag
/report     [YYYY-MM]            full breakdown: spent, earned, net + per-person totals
```

- No arg = current month; `YYYY-MM` selects a specific month (validated `01`–`12`)
- `/spendtotal` disambiguation: first token matching `^\d{4}-\d{2}$` is a
  month, otherwise it is a person tag — deterministic, no ambiguity
- Replies are plain inline text using the configured fixed currency
- Behavior split: **unrecognized** command → silently ignored (D6);
  **recognized but malformed** command → short usage-hint reply

**Rationale:** Inline queries with zero sheet-formula maintenance; month scope
matches how personal budgeting is thought about; per-person breakouts reuse
the `person` tag already on every spend row instead of adding new fields.

**Consequences:** The `Store` interface gains a read operation — v1 reads all
rows and aggregates in Go (single sheet, small volume, so full-read-per-query
is acceptable). Query flow: webhook → router → `Store` read → aggregate →
`sendMessage` reply. Sheet remains append-only on the write path.
