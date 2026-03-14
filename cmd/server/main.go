package main

import (
	"Netology/internal/worker"
	"context"
	"fmt"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logsCh, err := worker.ReadLogs("testdata/logs.csv")
	if err != nil {
		log.Fatal(err)
	}
	filteredCh := worker.ProcessLogs(ctx, logsCh, 5)
	stats := worker.CalculateStats(filteredCh)

	fmt.Printf("Всего запросов: %d\n", stats.TotalRequests)
	fmt.Printf("Ошибок: %d\n", stats.ErrorCount)
	fmt.Printf("Среднее время ответа: %.2f ms\n", stats.AverageRespTime)
}
