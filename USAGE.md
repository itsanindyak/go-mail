How to use / run the Email Campaign project

This file contains step-by-step commands to run the project locally, generate test data, and start the observability stack. Commands are shown for PowerShell and for WSL/Linux where differences matter.

Prerequisites

- Install Go (the repo uses Go 1.24.x).
- Install Docker & Docker Compose for observability.

Quick setup

1. Clone and change into the project directory:

```powershell
git clone <repo-url> email-campaign
cd email-campaign
```

2. Set required environment variables for real sending (do not store secrets in git). Example (PowerShell):

```powershell
$env:RESEND_API_KEY = "your_api_key_here"
```

Generate test CSV (PowerShell)

```powershell
# Create header
"name,email" | Out-File -FilePath email.csv -Encoding utf8
# Append 5000 test rows
1..5000 | ForEach-Object { "user$_,user$_@example.com" } | Out-File -FilePath email.csv -Append -Encoding utf8
```

Run the application

- Important: run the whole `cmd` package so all files in the package are compiled (do not `go run cmd/main.go` alone).

PowerShell / WSL:

```powershell
# Run without building a binary
go run ./cmd

# Or build then run the binary
go build -o email-campaign ./cmd
# PowerShell (Windows)
.\email-campaign.exe
# WSL / Linux
./email-campaign
```

Observability stack

1. Start Prometheus + Grafana with Docker Compose:

PowerShell (use resolved path for mounts if needed):

```powershell
$src = (Resolve-Path .\prometheus.yml).Path
docker compose -f observability.yaml up -d --build
```

WSL / Linux:

```bash
docker compose -f observability.yaml up -d --build
```

2. Check Prometheus targets: open `http://localhost:9090` → Status → Targets. Confirm the `email-campaign` target (scraping `:2112`) is `UP`.

Metrics endpoint

- The app exposes Prometheus metrics at `http://localhost:2112/metrics`. Example (PowerShell):

```powershell
Invoke-WebRequest http://localhost:2112/metrics -UseBasicParsing | Select-Object -ExpandProperty Content
```

Troubleshooting

- If Docker mount errors occur, use `Resolve-Path .\prometheus.yml` and pass the absolute path when mounting.
- If you see `undefined: StartMetrics` or similar, ensure you run `go run ./cmd` (whole package) or build the package — not `go run cmd/main.go` alone.

Next helper additions

- I can add a `start-dev.ps1` helper, a dry-run flag, or a small `make`/script to generate test data. Tell me which helper you prefer.
