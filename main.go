package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type ReportData struct {
	TotalTime          time.Duration
	TotalRequests      int
	SuccessfulRequests int
	StatusCounts       map[int]int
}

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

	startTime := time.Now()

	var wg sync.WaitGroup
	resultsChan := make(chan int, *requests)
	statusCounts := make(map[int]int)
	var successfulRequests int
	var mu sync.Mutex

	semaphore := make(chan struct{}, *concurrency)

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(reqNum int) {
			defer wg.Done()
			defer func() { <-semaphore }()
			resp, err := http.Get(*url)
			if err != nil {
				resultsChan <- 0
				return
			}
			defer resp.Body.Close()
			resultsChan <- resp.StatusCode
		}(i)
	}

	wg.Wait()
	close(resultsChan)

	for statusCode := range resultsChan {
		mu.Lock()
		if statusCode == http.StatusOK {
			successfulRequests++
		}
		statusCounts[statusCode]++
		mu.Unlock()
	}

	totalTime := time.Since(startTime)

	report := ReportData{
		TotalTime:          totalTime,
		TotalRequests:      *requests,
		SuccessfulRequests: successfulRequests,
		StatusCounts:       statusCounts,
	}

	printReport(report)
}

func printReport(data ReportData) {
	fmt.Println("\n--- Relatório do Teste de Carga ---")
	fmt.Printf("Tempo total gasto na execução: %s\n", data.TotalTime)
	fmt.Printf("Quantidade total de requests realizados: %d\n", data.TotalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", data.SuccessfulRequests)
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for status, count := range data.StatusCounts {
		if status != http.StatusOK && status != 0 { // Exclui 200 (já reportado) e 0 (erros de request)
			fmt.Printf("  Status %d: %d\n", status, count)
		}
	}
	if countZero, ok := data.StatusCounts[0]; ok {
		fmt.Printf("  Requests com erro (não foi possível obter status): %d\n", countZero)
	}
	fmt.Println("-----------------------------------")
}
