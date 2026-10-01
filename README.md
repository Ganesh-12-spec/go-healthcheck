# Go Healthcheck

A concurrent URL health checker built in Go to demonstrate **goroutines, channels, worker pools, semaphores, `sync.WaitGroup`, and context-based timeouts**.

Give it a list of URLs, and it checks them concurrently, reports their HTTP status codes and response times, and can continuously re-check them at a fixed interval.

## What It Does

* Reads target URLs from `config.json`
* Checks multiple URLs **concurrently**
* Limits concurrent requests with a configurable `max_concurrent` value
* Applies a **per-request timeout** using Go's `context`
* Reports HTTP status codes and response latency
* Displays healthy and failed checks with colored terminal output
* Supports continuous monitoring with the `--interval` flag
* Includes tests and **race-detector support**

## Example

```text
https://github.com  → 200 OK   (142ms)
https://google.com  → 200 OK   (466ms)
https://example.com → 200 OK   (183ms)
```

If a target fails or exceeds its timeout:

```text
https://unavailable.example → DOWN (timeout)
```

---

## Why I Built This

Backend systems often need to detect service failures quickly and reliably.

This project is a small-scale implementation of that idea, built primarily to understand **Go concurrency through real engineering problems rather than isolated examples**.

The project focuses on:

* Goroutines
* Channels
* Worker-pool patterns
* Buffered channels as semaphores
* `sync.WaitGroup`
* Context cancellation and timeouts
* Concurrent result collection
* Race-condition detection
* Clean package separation

---

## Configuration

Create or edit `config.json`:

```json
{
  "urls": [
    "https://google.com",
    "https://github.com"
  ],
  "timeout_seconds": 5,
  "max_concurrent": 3
}
```

### Configuration

| Field             | Description                                     |
| ----------------- | ----------------------------------------------- |
| `urls`            | URLs that should be monitored                   |
| `timeout_seconds` | Maximum time allowed for each request           |
| `max_concurrent`  | Maximum number of checks running simultaneously |

---

## Usage

### Run a single health-check cycle

```bash
go run main.go
```

### Run continuously

Re-check all URLs every 30 seconds:

```bash
go run main.go --interval 30s
```

Stop the monitor with:

```text
Ctrl+C
```

---

## Testing

Run the complete test suite:

```bash
go test ./... -v
```

Run the race detector:

```bash
go test ./... -race
```

The race detector helps identify unsafe concurrent access to shared memory.

---

## Architecture

```text
                    ┌──────────────┐
                    │   main.go    │
                    │  Orchestrator│
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   config.go  │
                    │ Load Config  │
                    └──────┬───────┘
                           │
                           ▼
                  ┌──────────────────┐
                  │  checker package │
                  └────────┬─────────┘
                           │
                 ┌─────────┴─────────┐
                 │                   │
                 ▼                   ▼
             CheckURL()          CheckAll()
           Single check        Concurrent checks
                 │                   │
                 │          ┌────────┴────────┐
                 │          │                 │
                 │      Goroutines        WaitGroup
                 │          │                 │
                 │      Semaphore             │
                 │          │                 │
                 └──────────┴─────────────────┘
                           │
                           ▼
                     Result Channel
                           │
                           ▼
                     main.go prints
                       the results
```

### Project Structure

```text
go-healthcheck/
│
├── main.go
│   └── Entry point, orchestration, output, interval loop
│
├── config.go
│   └── Config structure and config.json loading
│
├── checker/
│   └── checker.go
│       ├── Result
│       ├── CheckURL()
│       └── CheckAll()
│
├── config.json
│   └── Target URLs and runtime configuration
│
└── *_test.go
    └── Automated tests
```

The `main.go` file intentionally stays **thin**.

The actual health-checking and concurrency logic lives inside the `checker` package so that it can be tested independently.

---

# Concurrency Design

The core of the project is `CheckAll()`.

Instead of checking URLs sequentially:

```text
URL 1 → wait → URL 2 → wait → URL 3
```

the checker runs them concurrently:

```text
             CheckAll()
                 │
       ┌─────────┼─────────┐
       ▼         ▼         ▼
    URL 1      URL 2      URL 3
   goroutine  goroutine  goroutine
       │         │         │
       └─────────┼─────────┘
                 ▼
            Result channel
```

### Concurrency Limit

The number of simultaneously running checks is controlled by a **buffered channel used as a semaphore**.

For example:

```json
"max_concurrent": 3
```

means at most three URL checks can be active at the same time.

This prevents an unnecessarily large number of concurrent network requests.

Conceptually:

```text
max_concurrent = 3

┌───────┐
│ Slot 1│ → URL check
├───────┤
│ Slot 2│ → URL check
├───────┤
│ Slot 3│ → URL check
└───────┘

New check waits until a slot becomes available.
```

---

## Per-Request Timeout

Every HTTP request is wrapped in a `context.WithTimeout`.

Conceptually:

```text
URL Check
   │
   ▼
Context with timeout
   │
   ├── Request completes
   │       ↓
   │     Result
   │
   └── Timeout occurs
           ↓
         Error
```

This prevents a hanging or extremely slow endpoint from blocking the entire health-check cycle indefinitely.

---

## Synchronization

`sync.WaitGroup` is used to track all running goroutines.

The basic lifecycle is:

```text
Start goroutine
      ↓
WaitGroup.Add(1)
      ↓
Run health check
      ↓
Send result
      ↓
WaitGroup.Done()
      ↓
All goroutines complete
      ↓
Return results
```

The result channel allows the goroutines to send their results back to the coordinator safely.

---

# Key Go Concepts Demonstrated

### Goroutines

Run individual URL checks concurrently.

### Channels

Used to communicate results between goroutines.

### Buffered Channels

Used as a semaphore to limit concurrency.

### `sync.WaitGroup`

Waits until every health-check goroutine has finished.

### `context`

Provides cancellation and per-request timeouts.

### HTTP Client

Performs the actual URL health checks.

### Worker-Pool Thinking

The concurrency limit introduces a basic worker-pool pattern where available execution capacity controls how many jobs run simultaneously.

### Race Detection

The project can be tested with:

```bash
go test ./... -race
```

to detect certain classes of data races during concurrent execution.

---

# Design Principles

The project follows a simple separation of responsibilities:

```text
main.go
   │
   └── Orchestration
        │
        ├── Load configuration
        ├── Start checks
        ├── Handle interval
        └── Print results

checker/
   │
   └── Health-check logic
        │
        ├── HTTP request
        ├── Timeout
        ├── Concurrency
        ├── Result collection
        └── Synchronization
```

This keeps the entry point readable while making the core logic independently testable.

---

# What I Learned

Building this project helped me understand how Go handles concurrent I/O workloads in practice.

The main lessons were:

* How goroutines allow multiple network operations to run concurrently
* How channels provide communication between goroutines
* How buffered channels can act as concurrency limits
* How `sync.WaitGroup` coordinates multiple goroutines
* Why uncontrolled concurrency can become a problem
* How context-based timeouts prevent hanging operations
* How to structure a Go project into testable packages
* How to test concurrent code with Go's race detector

---

# Example Output

```text
Go Healthcheck
────────────────────────────────────────

✓ https://github.com   →
```
