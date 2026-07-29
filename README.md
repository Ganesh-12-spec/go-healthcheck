# go-healthcheck

A concurrent URL health checker written in Go. Give it a list of URLs, and it checks them all at the same time, reports status codes and response times, and can re-check on a fixed interval.

## What it does

- Reads a list of URLs from `config.json`
- Checks all of them **concurrently**, with a configurable max concurrency limit
- Applies a per-request timeout so a slow/hanging URL can't block the whole run
- Prints results with colored output (green = healthy, red = error/down)
- Optional `--interval` flag to keep re-checking on a loop

## Why

Real backend systems need to know when a service goes down — quickly, and automatically. This is a small-scale version of that idea, built to learn Go's concurrency primitives (goroutines, sync.WaitGroup, channels, semaphores, and context) by actually using them.

## Usage

1. Edit `config.json`:

​```json
{
  "urls": ["https://google.com", "https://github.com"],
  "timeout_seconds": 5,
  "max_concurrent": 3
}
​```

2. Run once:
​```bash
go run main.go
​```

3. Run on a loop (re-checks every 30 seconds):
​```bash
go run main.go --interval 30s
​```
Press Ctrl+C to stop.

## Running tests

​```bash
go test ./... -v
go test ./... -race
​```

## Architecture

​```
main.go            — entry point, loads config, prints results, handles the interval loop
checker/config.go  — Config struct + LoadConfig(): reads and parses config.json
checker/checker.go — Result struct, CheckURL() (single check with timeout),
                      CheckAll() (concurrent checks with a worker-pool limit)
​```

`main.go` stays thin — only orchestration. All logic lives in the `checker` package, testable independently.

### Concurrency design

- Each URL is checked in its own goroutine.
- A semaphore (buffered channel) caps how many checks run at once (`max_concurrent`).
- A `context.WithTimeout` wraps each request so a hanging URL can't stall the whole batch.
- Results are collected through a channel and returned once every goroutine finishes (`sync.WaitGroup`).

## Sample output

​```
https://github.com -> 200 (142ms)
https://google.com -> 200 (466ms)
​```
