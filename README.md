# Email Campaign — Lightweight, Production‑Ready Sender (Go)

A fast, observable, production‑minded email campaign backend built in Go.
Small codebase, big ideas: goroutine worker pools, DLQ, retries, metrics, and tracing — everything you need to send thousands of emails reliably *and* understand exactly what happened.

---

## 🚀 Why this repo is interesting

### **⚡ Concurrency done right**
Producer → bounded channel → worker pool.
This gives predictable memory usage, natural backpressure, and safe high throughput.

### **🛡️ Reliability built‑in**
Retries with exponential backoff + a durable DLQ. Nothing is silently lost.

### **📊 Observability‑first**
Prometheus metrics ("what"), Grafana dashboards, and OpenTelemetry traces ("why").

### **📈 Scales with you**
Starts as a simple monolith. Later plug in Redis/Kafka and run multiple worker services without rewriting logic.

### **⚙️ Transparent engineering trade‑offs**
Channel buffer, worker count, retry budget, rate limits — all configurable knobs.
No hidden magic.

---

## ✨ Features

- CSV → producer → channel ingestion
- Configurable worker pool
- Retries with exponential backoff
- Dead‑Letter Queue (file/DB)
- Prometheus metrics at `/metrics`
- Ready‑to‑import Grafana dashboard
- OpenTelemetry tracing (optional)
- Supports Resend / SES / SendGrid mail providers

---

## 🏁 Getting started (2 minutes)

Works on Windows (PowerShell), WSL, and Linux.


### **1) Start observability stack (Prometheus + Grafana)**

**PowerShell:**
```powershell
docker compose -f observability.yaml up -d --build
```

**WSL/Linux:**
```bash
docker compose -f observability.yaml up -d --build
```

**Grafana:** http://localhost:3000 → login: `admin/admin`
- If Grafana runs in Docker: Prometheus URL = `http://host.docker.internal:9090`
- If Grafana runs locally: Prometheus URL = `http://localhost:9090`

**Prometheus:** http://localhost:9090 → Status → Targets → must show **UP**

---

### **2) Run the app**
Use the whole `cmd` package — not `main.go` directly.

```bash
go run ./cmd
```

Or build:
```bash
go build -o email-campaign ./cmd
# Windows
./email-campaign.exe
# Linux/WSL
./email-campaign
```

Metrics exposed at: **http://localhost:2112/metrics**

---

## 🔧 Quick config knobs

- **worker_count** — how many goroutines (start with 5–20)
- **channel_buffer** — queue depth (50–500)
- **retry_attempts** — max retry count
- **backoff** — time between retries
- **mail_provider** — switch between SMTP/Resend/SES

---

## 📐 Architecture (1‑line diagram)
```
CSV Producer → bounded channel → N workers (goroutines) → mail provider
                                      ↓
                                 DLQ (file/DB)
```

---

## 🔍 Observability
### Metrics included
- `email_sent_total`
- `email_failed_total`
- `email_send_duration_seconds`
- `email_worker_active_count`
- `email_dlq_total`

### Dashboards
Grafana dashboard JSON is inside:
```
grafana/provisioning/dashboards/email-dashboard.json
```
Import manually or auto‑load via provisioning.

### Tracing (optional)
OpenTelemetry spans show:
- worker → send attempt → retry → provider → DLQ

View traces in Grafana (Tempo/Jaeger).

---

## 🏭 Production checklist

- Use SES / SendGrid / Resend (avoid raw SMTP in production)
- Configure SPF / DKIM / DMARC
- Consider Redis/Kafka for distributed queue
- Persist DLQ to Postgres
- Add Alertmanager alerts:
  - failure spikes
  - retry storms
  - queue depth growth
  - latency SLO violations

---

## 🛠️ Useful commands

Run tests:
```bash
go test ./...
```


Build Docker image:
```bash
docker build -t email-app:latest .
```

Start observability:
```bash
docker compose -f observability.yaml up -d --build
```

