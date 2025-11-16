
# Email Campaign — Planning, Concurrency Theory, and Observability

This README focuses on the conceptual approach to building a high-throughput email campaign service in Go. It intentionally emphasizes planning and theory before any implementation details. The goal is to outline why specific design choices (producer/consumer, channels, DLQ, metrics) work for sending thousands of emails reliably and observably.

---

## Plan: objectives, constraints, and design goals

Start by writing down explicit goals and limits — these drive all technical choices.

- Goals

  - Throughput: deliver 5k–10k emails in a controlled window.
  - Reliability: retries for transient failures, persistent DLQ for permanent failures.
  - Observability: real-time metrics and logs to understand system health.
  - Resource control: bounded memory and predictable CPU/network usage.

- Constraints

  - External provider limits (SMTP or API rate/connection limits).
  - Network bandwidth and host CPU.
  - Message size and templates (affects latency and bandwidth).

- Design implications
  - Concurrency must be tunable to match provider limits and host capacity.
  - The data pipeline should be simple and observable: ingest → queue → worker pool → DLQ.
  - Make shutdown and recovery explicit: preserve in-flight messages, or document acceptable loss behavior.

Why this upfront plan matters

When you plan first, you avoid reactive changes later. For example, deciding to treat the DLQ as durable storage (file/DB) early prevents accidental message loss during scale tests. Explicit rate-limiting strategies reduce the chance of being blocked by email providers.

---

## Send Mail via Goroutines — conceptual model

This section explains the mental model for concurrency and data transfer without getting into specific code. Think of goroutines as independent workers and channels as the pipes that connect them.

Core components (conceptually)

- Producer (single logical actor)

  - Reads recipients from a source (CSV, DB, or API) and pushes work items into a queue.
  - Operates as the only writer to the queue so the source of truth is simple and sequenced.

- Queue (channel)

  - A bounded or unbounded buffer that decouples producer speed from worker speed.
  - Acts as backpressure: when the queue fills, the producer slows (or blocks) preventing runaway memory use.

- Worker pool (many concurrent consumers)

  - A configurable number of concurrent workers process items from the queue.
  - Each worker performs the send operation and reports success/failure to observability.

- Dead-letter queue (DLQ)
  - Permanently failed items are sent to DLQ for later inspection and reprocessing.
  - DLQ is typically persisted by one consumer to avoid concurrent writes.

Coordination and lifecycle

- Graceful shutdown: carry a cancellation token or signal that tells producer to stop producing and workers to finish in-flight messages or stop immediately depending on policy.
- Completion detection: a producer closing the queue signals no more messages; consumers exit when queue is drained.

Backpressure, buffering, and trade-offs

- Small buffers: keep memory low and force producer to match consumer speed — useful when you want tight control.
- Large buffers: absorb spikes but increase memory usage and reduce immediate feedback on downstream issues.
- Unbuffered channels: force handoff — producer blocks until a worker accepts the message (strong backpressure).

Retry and failure strategy (mental model)

- Retry transient errors with exponential backoff and a capped retry budget per message.
- On exceeding retries, move to DLQ with context (error, attempt count, timestamp).

Scaling horizontally

- If a single instance cannot reach desired throughput or you want redundancy, run multiple instances that share a centralized work source (queue) or partition work by ranges of the input (e.g., CSV slices or DB offsets).

---

## Observability — what to capture and why

Observability is essential: without it, tuning concurrency is guesswork.

Key metric families (conceptual)

- Counters

  - emails_sent_total: absolute successes.
  - emails_failed_total: permanent failures.
  - email_retries_total: how many retry attempts occur (shows instability).

- Histograms / Summaries

  - email_send_duration_seconds: distribution of send latency (P50, P95, P99).

- Gauges
  - workers_active: how many workers are currently processing items.
  - queue_length: current number of items waiting in queue (if trackable).

Logs and traces

- Structured logs: include message id or hashed recipient, attempt number, error codes, and timestamps.
- Tracing: add spans for send attempts to see cross-system latency (if using external APIs).

Dashboards and alerts (conceptual)

- Dashboard panels

  - Send rate (per second), success ratio, retry rate, queue depth, worker utilization, recent DLQ samples.

- Alerts to configure
  - Elevated failure rate (over a sliding window).
  - Rising retry rate without corresponding recoveries.
  - Queue depth steadily growing (consumers can't keep up).

How to wire observability (theory)

- Expose metrics from the process over an HTTP endpoint in Prometheus format.
- Use a metrics backend (Prometheus) to scrape and store metrics; use Grafana for visualization and alerting.

---

## Capacity planning and tuning (guidelines)

Start with conservative settings and iterate with measurements.

- Worker count: begin with a low number (e.g., 5–10), measure send latency and success rate, then increase slowly while watching failure and retry rates.
- Channel buffer: small buffer sizes (e.g., 50–200) are safe for memory; larger sizes require testing.
- Rate limiting: implement a token-bucket if provider quotas are strict; tune tokens/second to the provider's allowed throughput.

Example trade-offs

- More workers increase throughput but also increase concurrency at the provider (more connections, more rate-limited responses).
- Larger buffers reduce producer backpressure but hide immediate downstream problems and increase memory footprint.

---

## Next steps (recommended)

- Add a dry-run mode for safe load testing without sending real emails.
- Replace hard-coded credentials with environment variables and document required secrets.
- Add a short load-testing plan: CSV generator, monitoring checklist, and target metrics to observe.
- Optionally, provide a short tuning guide with sample results from small experiments on different machine sizes.

If you want, I can now add a small conceptual tuning guide or create a load-test plan and a dry-run generator. Which would you like next?
