# Engage — adaptive student engagement

Observe → Understand → Adapt → Measure → Learn. Behavioural telemetry → deterministic engagement
assessment → deterministic intervention → personalised activity (rules, optionally LLM-rewritten) →
new attempt → measured effectiveness.

## Run
```bash
# terminal 1 — backend (http://localhost:8080). Migrates + seeds a starter library on first run.
cd backend && go mod tidy && go run ./cmd/api

# terminal 2 — frontend (http://localhost:5173)
cd frontend && npm install && npm run dev
```
Open the Dashboard and press **Run live demo**, then see **History**.

## How it answers the challenge
- **Detects the full spectrum**: struggling, bored, and drifting-away students are told apart from aggregate behaviour (no content capture); the demo plays all three stories.
- **Listens to the student**: a one-tap check-in (stuck / bored / curious / want company / want to play) shapes the next activity with a connection, play or curiosity twist.
- **Community**: an anonymous "community pulse" shows what learners feel right now.
- **Closed loop**: every adaptation is measured (EFFECTIVE / NEUTRAL / INEFFECTIVE) and shown in History.
- **Low-connectivity friendly**: lightweight API, aggregate telemetry, and a cached activity library with an offline notice. Starting attempts and syncing progress still need the server.

## Optional: LLM content adaptation
The LLM only rewrites activity content; it never decides engagement state or intervention, its
output is validated, and any failure falls back to the rule-based adapter.

Either paste a key in the Dashboard's "AI content adaptation" card (kept in server memory only), or:
```bash
export ANTHROPIC_API_KEY=...        # Claude (default model claude-sonnet-5-5)
# or
export OPENAI_API_KEY=...           # any OpenAI-compatible API
export LLM_BASE_URL=https://...     # optional, e.g. OpenRouter/Groq/local
export LLM_MODEL=...                # optional
export LLM_PROVIDER=anthropic|openai
```

## Privacy
Telemetry is aggregate only (durations, idle time, tab switches, interaction counts). No keystrokes,
mouse coordinates or message contents are captured.
