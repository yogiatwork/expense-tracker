# Expense Tracker

Turns WhatsApp messages into Google Sheets rows. A user sends something like
"/spend 200 tesco sapna" (or "/earned 149 online") to a WhatsApp Business
number; Meta's servers push the message to our webhook endpoint; a Go service
parses it, appends a debit/credit ledger row to a Google Sheet, and replies
with a confirmation.

## Tech Stack

- **Language:** Go
- **Input:** WhatsApp Business API (provider TBD — see `docs/open-questions.md`)
- **Storage:** Google Sheets (service account auth)
- **Integration pattern:** webhooks (push), not polling

## Working Agreement (important)

- **The developer writes all production code.** The assistant acts strictly as
  mentor and code reviewer: explain concepts, discuss design, review diffs,
  point out idiomatic patterns — never implement features on the user's behalf.
- Requirements and design are discussed before code is written.
- External systems (WhatsApp provider, Google Sheets) must sit behind Go
  interfaces so they can be mocked in tests and swapped later.
- Prefer small, table-driven-tested packages over framework magic.

## Session Context

Read these before doing any work — they carry decisions from prior sessions:

| File | Contents |
|---|---|
| `docs/roadmap.md` | Project phases and current status |
| `docs/architecture.md` | Components, data flow, webhook mechanics |
| `docs/decisions.md` | Append-only decision log with rationale |
| `docs/open-questions.md` | Unresolved design questions |

## Current Status

Phase 1–2 in progress (requirements gathering + architecture discussion).
No code exists yet.
