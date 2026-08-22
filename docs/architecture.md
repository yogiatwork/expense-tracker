# Architecture

## High-Level Data Flow

```
┌──────┐  1. "/spend 200 tesco sapna"   ┌─────────┐  2. HTTP POST (JSON)   ┌──────────────┐
│ User │ ─────────────────────────────► │  Meta   │ ─────────────────────► │   Go Web     │
│      │                                │ Servers │     to /webhook        │   Service    │
└──────┘                                └────▲────┘                        └──────┬───────┘
                                             │                                    │ 3. parse + append row
                                             │ 4. send confirmation               ▼
                                             └──────────────────────────── ┌──────────────┐
                                                 (outbound Graph API call) │ Google Sheet │
                                                                           └──────────────┘
```

## Components (planned)

| Component | Responsibility |
|---|---|
| **Webhook handler** | Accept POST from Meta; verify signature; ack 200 fast |
| **Verifier** | One-time GET handshake (`hub.challenge` echo) |
| **Command router + parser** | Dispatch slash commands; text → `Transaction` (debit \| credit), regex-first |
| **Store interface** | Append/query rows; Google Sheets implementation |
| **Provider interface** | Send replies; Meta/Twilio implementations |
| **Deduper** | Drop retried messages by WhatsApp message ID |

Design rule: external systems sit behind Go interfaces so implementations can
be swapped or mocked in tests.

## Webhook Mechanics

### Registration (one-time)

1. Expose the service at a public HTTPS URL.
2. Register the URL in the Meta Developer Console.

### Verification handshake (one-time GET)

Meta proves you own the endpoint before activating it:

```
GET /webhook?hub.mode=subscribe&hub.verify_token=SECRET&hub.challenge=1039847
```

Handler must: if `hub.verify_token` matches the configured secret, respond
`200` with `hub.challenge` as the plain-text body; otherwise `403`.

### Message delivery (recurring POST)

```json
{
  "object": "whatsapp_business_account",
  "entry": [{
    "changes": [{
      "value": {
        "messages": [{
          "from": "919876543210",
          "text": { "body": "/spend 200 tesco sapna" },
          "type": "text",
          "timestamp": "1755772800"
        }]
      },
      "field": "messages"
    }]
  }]
}
```

The reply direction (step 4 above) is a normal outbound API call — not a webhook.

### Signature verification (every POST)

Each POST carries `X-Hub-Signature-256`: HMAC-SHA256 of the raw body keyed
with the app secret. Verify before parsing so attackers can't forge expenses
by hitting the public URL directly. Requires reading the raw body *before*
JSON decoding.

## Command Protocol

Only messages starting with a known slash command are processed; everything
else is silently ignored — Meta delivers *all* messages, so filtering happens
here, not at WhatsApp (D6).

```
/spend [today|yesterday] <amount> <merchant> <person-tag>
/earned <amount> <method>
```

| Message | Result |
|---|---|
| `/spend 200 tesco sapna` | debit 200 · merchant tesco · person sapna · today |
| `/spend 45 chai sapna` | debit 45 · merchant chai · person sapna |
| `/earned 149 online` | credit 149 · method online |
| `/earned 23 cash` | credit 23 · method cash |
| anything else | ignored |

Rules:

- Currency never appears in messages; one fixed currency lives in service
  config.
- Both commands map onto a single model the rest of the pipeline consumes:

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

- Grammar refinements pending: Q5 (word boundaries, optionality, dates),
  Q7 (field symmetry).

## Operational Constraints

- **Public HTTPS required** — local dev uses ngrok / Cloudflare Tunnel.
- **Meta retries on non-200 / slow responses** → same message may arrive more
  than once → dedupe on the message ID (`wamid...`) before writing rows.
- **Ack fast:** return `200` immediately; process async if handling grows slow.
- **Idempotency is a correctness requirement**, not an optimization — a retried
  webhook must not create a duplicate transaction row.
