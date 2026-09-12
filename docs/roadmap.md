# Roadmap

Status: `[ ]` pending · `[~]` in progress · `[x]` done

## Phase 1 — Requirements & Data Model `[~]`

**Goal:** Pin down scope and the data model before writing any code.

- Message format — DECIDED: slash-command protocol, silent-ignore unknowns,
  config-held currency (D6)
- Income tracking via `/earned` alongside expenses (D7)
- Row schema — ledger form: date, description, debit, credit, person,
  method, raw_message, update_id (D7)
- Grammar refinements: multi-word merchants, person optionality, dates (Q5)
- Balance mechanism (Q6), field symmetry (Q7)
- User scope: single user / family / open (Q2)

**Exit criteria:** every item in `docs/open-questions.md` is answered; row
schema sketched.

## Phase 2 — Architecture & Design `[~]`

**Goal:** Component boundaries, data flow, package layout, error strategy.

- Data flow diagram (see `architecture.md`)
- Go interfaces for Telegram provider + Sheets store (mockable)
- Idempotency strategy for webhook retries

**Exit criteria:** architecture doc reviewed and agreed.

## Phase 3 — Scaffolding `[ ]`

**Goal:** Repo skeleton ready for feature work.

- Go module layout, config management (env vars / secrets)
- Logging setup, Makefile, CI skeleton

## Phase 4 — Telegram Integration `[ ]`

**Goal:** Receive real messages end-to-end.

- HTTP server with `/webhook` endpoint
- `setWebhook` registration (url + secret_token)
- `X-Telegram-Bot-Api-Secret-Token` header verification
- `sendMessage` client for confirmation replies
- Local dev via ngrok / Cloudflare Tunnel

## Phase 5 — Parsing Engine `[ ]`

**Goal:** Turn message text into a structured `Transaction`.

- Rule-based parser first (regex → struct), LLM later if needed
- Pure Go logic, no external deps — table-driven tests

## Phase 6 — Google Sheets Integration `[ ]`

**Goal:** Persist a transaction as a sheet row.

- Service account auth, minimal OAuth scope
- Ledger columns: exactly one of debit/credit populated per row (D7)
- Append-row operation behind a `Store` interface
- Read operation: query transactions for month/person aggregation (D10)

## Phase 7 — Core Business Logic `[ ]`

**Goal:** Correctness under real-world messiness.

- Validation (amounts, known methods/dates)
- Deduplication on Telegram `update_id` (Telegram retries webhooks)

## Phase 8 — Queries & Reports `[ ]`

**Goal:** Read path from the sheet.

- `/balance [YYYY-MM]` — spent, earned, net for the month (D10)
- `/spendtotal [YYYY-MM] [person]` — month spend total, optional person tag (D10)
- `/report [YYYY-MM]` — full breakdown incl. per-person totals (D10)
- Aggregation in Go from a `Store` read; inline `sendMessage` replies

## Phase 9 — Testing `[ ]`

**Goal:** Confidence to refactor.

- Unit tests with mocked provider/store
- Integration tests against a Telegram test bot
- Ongoing practice from Phase 5 onward, not a one-off

## Phase 10 — Deployment & Hardening `[ ]`

**Goal:** Runs unattended.

- Hosting with public HTTPS, secrets management
- Retries, rate limits, monitoring/alerting
