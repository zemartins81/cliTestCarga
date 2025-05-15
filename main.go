package main

import (
	"flag"     // Pacote para análise de argumentos da linha de comando
	"fmt"      // Pacote para formatação de entrada/saída
	"net/http" // Pacote para fazer requisições HTTP
	"os"       // Pacote para interação com o sistema operacional
	"sync"     // Pacote para primitivas de sincronização
	"time"     // Pacote para manipulação de tempo
)

// ReportData estrutura para armazenar os resultados do teste de carga
type ReportData struct {
	TotalTime          time.Duration // Tempo total da execução do teste
	TotalRequests      int           // Número total de requisições realizadas
	SuccessfulRequests int           // Número de requisições com status HTTP 200
	StatusCounts       map[int]int   // Contagem de respostas por código de status HTTP
}

func main() {
	// Definição e análise dos argumentos da linha de comando
	url := flag.String("url", "", "url")                      // URL alvo do teste de carga
	requests := flag.Int("requests", 100, "requests")         // Número total de requisições
	concurrency := flag.Int("concurrency", 10, "concurrency") // Número máximo de requisições simultâneas
	flag.Parse()                                              // Analisa os argumentos fornecidos

	// Validação dos argumentos fornecidos
	if *url == "" {
		fmt.Println("url is required")
		os.Exit(1) // Encerra o programa com código de erro
	}

	if *requests <= 0 {
		fmt.Println("requests must be greater than 0")
		os.Exit(1)
	}

	if *concurrency <= 0 {
		fmt.Println("concurrency must be greater than 0")
		os.Exit(1)
	}

	// Ajusta a concorrência se for maior que o número de requisições
	if *requests < *concurrency {
		fmt.Println("requests must be greater than concurrency")
		*concurrency = *requests
	}

	// Marca o tempo de início do teste
	startTime := time.Now()

	// Configuração das estruturas de sincronização e comunicação
	var wg sync.WaitGroup                    // Grupo de espera para aguardar todas as goroutines finalizarem
	resultsChan := make(chan int, *requests) // Canal para comunicar os códigos de status das requisições
	statusCounts := make(map[int]int)        // Mapa para contar ocorrências de cada código de status
	var successfulRequests int               // Contador de requisições bem-sucedidas (status 200)
	var mu sync.Mutex                        // Mutex para acesso seguro aos contadores compartilhados

	// Semáforo para limitar o número de goroutines concorrentes
	semaphore := make(chan struct{}, *concurrency)

	// Inicia as goroutines para fazer as requisições HTTP
	for i := 0; i < *requests; i++ {
		wg.Add(1)               // Incrementa o contador do WaitGroup
		semaphore <- struct{}{} // Adquire um slot no semáforo
		go func(reqNum int) {
			defer wg.Done()                // Decrementa o contador do WaitGroup ao finalizar
			defer func() { <-semaphore }() // Libera um slot no semáforo ao finalizar

			// Realiza a requisição HTTP
			resp, err := http.Get(*url)
			if err != nil {
				resultsChan <- 0 // Envia 0 para o canal em caso de erro
				return
			}
			defer resp.Body.Close()        // Garante que o corpo da resposta seja fechado
			resultsChan <- resp.StatusCode // Envia o código de status para o canal
		}(i)
	}

	wg.Wait()          // Aguarda todas as goroutines finalizarem
	close(resultsChan) // Fecha o canal de resultados após todas as goroutines finalizarem

	// Processa os resultados recebidos pelo canal
	for statusCode := range resultsChan {
		mu.Lock() // Adquire o mutex para acesso seguro aos contadores
		if statusCode == http.StatusOK {
			successfulRequests++ // Incrementa o contador de requisições bem-sucedidas
		}
		statusCounts[statusCode]++ // Incrementa a contagem para este código de status
		mu.Unlock()                // Libera o mutex
	}

	// Calcula o tempo total gasto
	totalTime := time.Since(startTime)

	// Prepara os dados do relatório
	report := ReportData{
		TotalTime:          totalTime,
		TotalRequests:      *requests,
		SuccessfulRequests: successfulRequests,
		StatusCounts:       statusCounts,
	}

	// Exibe o relatório final
	printReport(report)
}

// printReport exibe na saída padrão um relatório formatado com os resultados do teste de carga
func printReport(data ReportData) {
	fmt.Println("\n--- Relatório do Teste de Carga ---")
	fmt.Printf("Tempo total gasto na execução: %s\n", data.TotalTime)
	fmt.Printf("Quantidade total de requests realizados: %d\n", data.TotalRequests)
	fmt.Printf("Quantidade de requests com status HTTP 200: %d\n", data.SuccessfulRequests)
	fmt.Println("Distribuição de outros códigos de status HTTP:")
	for status, count := range data.StatusCounts {
		// Exibe todos os códigos exceto 200 (já reportado) e 0 (erros de requisição)
		if status != http.StatusOK && status != 0 {
			fmt.Printf("  Status %d: %d\n", status, count)
		}
	}
	// Exibe a quantidade de requisições com erro (que não retornaram código de status)
	if countZero, ok := data.StatusCounts[0]; ok {
		fmt.Printf("  Requests com erro (não foi possível obter status): %d\n", countZero)
	}
	fmt.Println("-----------------------------------")
}
