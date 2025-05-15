package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	url := flag.String("url", "", "url")
	requests := flag.Int("requests", 100, "requests")
	concurrency := flag.Int("concurrency", 10, "concurrency")
	flag.Parse()

	if *url == "" {
		fmt.Println("url is required")
		os.Exit(1)
	}

	if *requests <= 0 {
		fmt.Println("requests must be greater than 0")
		os.Exit(1)
	}

	if *concurrency <= 0 {
		fmt.Println("concurrency must be greater than 0")
		os.Exit(1)
	}

	if *requests < *concurrency {
		fmt.Println("requests must be greater than concurrency")
		*concurrency = *requests
	}

}
