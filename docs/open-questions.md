# Open Questions

Unresolved design questions. Each gets answered, then moves into
`decisions.md` as a dated decision. Newest at bottom.

---

## Q1 — Input provider choice *(answered 2026-08-22 → D8)*

**Answered:** Telegram Bot API chosen — instant bot creation via BotFather,
no business verification, simpler webhook setup, long-polling available for
dev. Supersedes the earlier Meta/Twilio comparison.

## Q2 — User scope: who can log expenses? *(open)*

**Options & tradeoffs:**

- **Just me** — simplest; no auth model beyond ignoring unknown numbers.
- **Me + family/friends** — phone-number allowlist; still single sheet.
- **Open to anyone** — needs per-user sheets/onboarding; heavy for v1.

**Leaning:** start single-user with an allowlist check anyway (cheap to add,
needed for security regardless of scope).

## Q3 — Sheet layout: single tab vs monthly tabs? *(open)*

**Options & tradeoffs:**

- **Single append-only tab** — simplest writes; reports filter by date column;
  cross-month totals trivial.
- **Monthly tabs** — human-friendly browsing per month; harder cross-month
  reports; more write logic (pick/create tab).

**Leaning:** single tab + date column; pivot tables in Sheets handle monthly
views natively.

## Q4 — Message format: structured shorthand vs natural language? *(answered 2026-08-22 → D6)*

**Options & tradeoffs:**

- **Structured shorthand** ("spent 250 on groceries", "250 food") — regex-
  friendly, deterministic, matches D5.
- **Natural language** ("had biryani with friends, paid 800 by card") — needs
  LLM/NLP; flexible but costly, slower, occasionally wrong.
- **Hybrid** — regex MVP now, LLM parser behind the same interface later.

**Answered:** slash-command protocol adopted (`/spend`, `/earned`); unknown
messages silently ignored; currency held in config. Refinements moved to
Q5/Q7.

## Q5 — /spend grammar refinements *(open)*

Raised while defining D6; deferred:

- **Multi-word merchants:** positional parsing breaks on
  `/spend 200 big bazaar sapna`. Options: underscores (`big_bazaar`),
  quoted merchant, or last-token-is-always-person.
- **Person field:** required (removes all ambiguity) vs optional (needs a
  disambiguation rule).
- **Date vocabulary:** keywords only (`today` default, `yesterday`) vs also
  ISO dates for backfill vs day names within the week.

## Q6 — Balance/query mechanism *(answered 2026-08-22 → D10)*

**Decision:** on-demand commands (`/balance`, `/spendtotal`, `/report`),
month-scoped, computed by the service from a `Store` read. Scheduled summary
rows rejected (adds infra + dedupe concerns). Sheet formulas remain an
optional read-only convenience in the spreadsheet, not part of the service.

- **Sheet formulas** (SUMIFS summary block/tab) — zero service code, always
  current.
- **`/balance` command** — service computes totals and replies in chat.
- **Scheduled summary rows** — cron appends weekly rows; adds infra plus its
  own dedupe concerns.

**Answered:** command-based approach (D10); full grammar in `decisions.md`.

## Q7 — Field symmetry between /spend and /earned *(open)*

- Should `/earned` accept an optional source word ("freelance")?
- Should `/spend` also record payment method (cash/card/online)?
- Do we need a category concept beyond free-form merchant + person tag?
