package main

import (
	"fmt"
	"log"

	"github.com/Ganesh-12-spec/go-healthcheck/checker"
)

func main() {
	cfg, err := checker.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded config: %+v\n", cfg)
	results := checker.CheckALL(cfg.URLs, cfg.MaxConcurrent)
	fmt.Printf("Check results: %+v\n", results)
}