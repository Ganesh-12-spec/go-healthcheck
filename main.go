package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Ganesh-12-spec/go-healthcheck/checker"
)

const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

func main() {
	interval := flag.Duration("interval", 0, "re-check interval, e.g. 30s (0 = run once)")
	flag.Parse()

	cfg, err := checker.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	fmt.Println("DEBUG timeout value:", timeout)

	for {
		results := checker.CheckALL(cfg.URLs, cfg.MaxConcurrent, timeout)

		for _, r := range results {
			if r.Err != nil {
				fmt.Printf("%s%s -> ERROR: %v%s\n", colorRed, r.URL, r.Err, colorReset)
			} else if r.StatusCode >= 200 && r.StatusCode < 300 {
				fmt.Printf("%s%s -> %d (%v)%s\n", colorGreen, r.URL, r.StatusCode, r.Duration, colorReset)
			} else {
				fmt.Printf("%s%s -> %d (%v)%s\n", colorRed, r.URL, r.StatusCode, r.Duration, colorReset)
			}
		}

		if *interval == 0 {
			break
		}
		fmt.Println("---waiting for next check---")
		time.Sleep(*interval)
	}
}