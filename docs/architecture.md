# Architecture

## High-Level Data Flow

```
┌──────┐  1. "/spend 200 tesco sapna"   ┌──────────┐  2. HTTP POST (JSON)   ┌──────────────┐
│ User │ ─────────────────────────────► │ Telegram │ ────────────────────► │  Go Service  │
└──────┘                                │  Servers │    to /webhook       └──────┬───────┘
                                        └────▲─────┘                              │
                                             │                                    │ 3. parse
                                             │ 4. sendMessage (Bot API call)      ▼
                                             │                            ┌──────────────┐
                                             └────────────────────────────►│ Google Sheet │
                                                                            └──────────────┘
```

## Components (planned)

| Component | Responsibility |
|---|---|
| **Webhook handler** | Accept POST from Telegram; verify `X-Telegram-Bot-Api-Secret-Token`; ack 200 fast |
| **Command router + parser** | Dispatch slash commands; text → `Transaction` (debit \| credit), regex-first |
| **Store interface** | Append rows; read transactions for aggregation(query); Google Sheets implementation |
| **Provider interface** | Send replies; Telegram Bot API implementation |
| **Deduper** | Drop retried messages by `update_id` |

Design rule: external systems sit behind Go interfaces so implementations can
be swapped or mocked in tests.

## Webhook Mechanics

### Registration (one-time)

1. Expose the service at a public HTTPS URL (ngrok / Cloudflare Tunnel).
2. Register the URL with Telegram via the Bot API:

```
POST https://api.telegram.org/bot<TOKEN>/setWebhook?url=https://<your-tunnel>/webhook&secret_token=<RANDOM>
```

No GET handshake — Telegram validates ownership purely through the
`secret_token` it sends on each subsequent POST.

### Message delivery (recurring POST)

```json
{
  "update_id": 9527,
  "message": {
    "message_id": 1,
    "from": { "id": 12345, "first_name": "Yogi", "username": "yogi" },
    "chat": { "id": 12345, "type": "private" },
    "date": 1755772800,
    "text": "/spend 200 tesco sapna"
  }
}
```

Key fields mapped to our schema:

- `message.text` → command input for the parser
- `message.chat.id` → user identity (for the allowlist, Q2)
- `update_id` → dedupe key (replaces `wamid`)
- `message.date` → Unix timestamp

The reply direction (step 4 above) is a normal outbound API call — not a
webhook:

```
POST https://api.telegram.org/bot<TOKEN>/sendMessage
{ "chat_id": 12345, "text": "debited 200 (tesco, sapna)" }
```

### Signature verification (every POST)

Each POST carries `X-Telegram-Bot-Api-Secret-Token`: a pre-shared plain
string set during `setWebhook`. Verify with a constant-time string compare
(`crypto/subtle.ConstantTimeCompare`) before parsing — prevents third
parties from forging updates by hitting the public URL directly.

## Command Protocol

Only messages starting with a known slash command are processed; everything
else is silently ignored — Telegram delivers *all* messages, so filtering
happens here, not at the bot layer (D6).

```
Write commands:
/spend [today|yesterday] <amount> <merchant> <person-tag>
/earned <amount> <method>

Read commands:
/balance    [YYYY-MM]
/spendtotal [YYYY-MM] [person]
/report     [YYYY-MM]
```

| Message | Result |
|---|---|
| `/spend 200 tesco sapna` | debit 200 · merchant tesco · person sapna · today |
| `/spend 45 chai sapna` | debit 45 · merchant chai · person sapna |
| `/earned 149 online` | credit 149 · method online |
| `/earned 23 cash` | credit 23 · method cash |
| `/balance` | current month: spent · earned · net |
| `/balance 2026-08` | that month's spent · earned · net |
| `/spendtotal sapna` | current month spend, person = sapna |
| `/report 2026-08` | full month breakdown incl. per-person totals |
| unrecognized text | ignored |
| recognized but malformed | short usage-hint reply |

Rules:

- Currency never appears in messages; one fixed currency lives in service
  config.
- Read commands: no arg = current month; `YYYY-MM` selects a specific month
  (validated `01`–`12`). `/spendtotal` disambiguation: first token matching
  `^\d{4}-\d{2}$` is a month, otherwise a person tag (D10).
- All commands map onto a single model the rest of the pipeline consumes:

```go
type Transaction struct {
    Direction Direction // debit | credit
    Amount    int
    Date      time.Time
    Merchant  string    // /spend only
    Person    string    // /spend only (tag)
    Method    string    // /earned only
}
```

- **Read path** for query commands: webhook → router → `Store` read
  (v1: read all rows, aggregate in Go) → digest → `sendMessage` reply.
  The sheet stays append-only on the write path; rows are clean because
  dedupe happens before write.
- Grammar refinements pending: Q5 (word boundaries, optionality, dates),
  Q7 (field symmetry).

## Operational Constraints

- **Public HTTPS required for webhooks** — local dev uses ngrok / Cloudflare
  Tunnel; long-polling (`getUpdates`) available for early testing without a
  tunnel but not suitable for production.
- **Telegram retries on non-200 / slow responses** → same update may arrive
  more than once → dedupe on `update_id` before writing rows.
- **Rate limit:** 30 messages/sec to different chats — single-user v1 well
  within limits.
- **Ack fast:** return `200` immediately; process async if handling grows slow.
- **Idempotency is a correctness requirement**, not an optimization — a retried
  webhook must not create a duplicate transaction row.
