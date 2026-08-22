# Roadmap

Status: `[ ]` pending · `[~]` in progress · `[x]` done

## Phase 1 — Requirements & Data Model `[~]`

**Goal:** Pin down scope and the data model before writing any code.

- Message format — DECIDED: slash-command protocol, silent-ignore unknowns,
  config-held currency (D6)
- Income tracking via `/earned` alongside expenses (D7)
- Row schema — ledger form: date, description, debit, credit, person,
  method, raw_message, wamid (D7)
- Grammar refinements: multi-word merchants, person optionality, dates (Q5)
- Balance mechanism (Q6), field symmetry (Q7)
- User scope: single user / family / open (Q2)

**Exit criteria:** every item in `docs/open-questions.md` is answered; row
schema sketched.

## Phase 2 — Architecture & Design `[~]`

**Goal:** Component boundaries, data flow, package layout, error strategy.

- Data flow diagram (see `architecture.md`)
- Go interfaces for WhatsApp provider + Sheets store (mockable)
- Idempotency strategy for webhook retries

**Exit criteria:** architecture doc reviewed and agreed.

## Phase 3 — Scaffolding `[ ]`

**Goal:** Repo skeleton ready for feature work.

- Go module layout, config management (env vars / secrets)
- Logging setup, Makefile, CI skeleton

## Phase 4 — WhatsApp Integration `[ ]`

**Goal:** Receive real messages end-to-end.

- HTTP server with `/webhook` endpoint
- Verification handshake (GET `hub.challenge` echo)
- `X-Hub-Signature-256` HMAC verification
- Send-message client for confirmation replies
- Local dev via ngrok / Cloudflare Tunnel

## Phase 5 — Parsing Engine `[ ]`

**Goal:** Turn message text into a structured `Expense`.

- Rule-based parser first (regex → struct), LLM later if needed
- Pure Go logic, no external deps — table-driven tests

## Phase 6 — Google Sheets Integration `[ ]`

**Goal:** Persist a transaction as a sheet row.

- Service account auth, minimal OAuth scope
- Ledger columns: exactly one of debit/credit populated per row (D7)
- Append-row operation behind a `Store` interface

## Phase 7 — Core Business Logic `[ ]`

**Goal:** Correctness under real-world messiness.

- Validation (amounts, known methods/dates)
- Deduplication on WhatsApp message ID (Meta retries webhooks)

## Phase 8 — Queries & Reports `[ ]`

**Goal:** Read path from the sheet.

- "How much this month?", "report food" style queries
- Weekly/monthly balance views — mechanism per Q6

## Phase 9 — Testing `[ ]`

**Goal:** Confidence to refactor.

- Unit tests with mocked provider/store
- Integration tests against sandbox accounts
- Ongoing practice from Phase 5 onward, not a one-off

## Phase 10 — Deployment & Hardening `[ ]`

**Goal:** Runs unattended.

- Hosting with public HTTPS, secrets management
- Retries, rate limits, monitoring/alerting
