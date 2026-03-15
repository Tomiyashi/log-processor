package main

import (
	"Netology/internal/worker"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go <путь_к_файлу>")
		fmt.Println("Пример: go run main.go testdata/logs.csv")
		os.Exit(1)
	}

	filename := os.Args[1]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logsCh, err := worker.ReadLogs(ctx, filename)
	if err != nil {
		log.Fatal(err)
	}

	processedCh := worker.ProcessLogs(ctx, logsCh, 5)

	stats := worker.CalculateStats(processedCh)

	fmt.Printf("Статистика логов:\n")
	fmt.Printf("====================\n")
	fmt.Printf("Всего запросов: %d\n", stats.TotalRequests)
	fmt.Printf("Ошибок (4xx/5xx): %d\n", stats.ErrorCount)
	fmt.Printf("Среднее время ответа: %.2f ms\n", stats.AverageRespTime)

	worker.PrintTopIPs(stats.RequestsByIP, 3)
}
