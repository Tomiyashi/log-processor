package worker

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

func ReadLogs(filename string) (<-chan LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не получилось открыть файл %v", filename)
	}

	ch := make(chan LogEntry, 100)

	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := ParseLogLine(line)
			if err != nil {
				log.Print("Ошибка чтения", err)
				continue
			}
			ch <- entry
		}
		close(ch)
	}()

	return ch, nil
}

func ProcessLogs(ctx context.Context, input <-chan LogEntry, numWorkers int) <-chan LogEntry {
	var wg sync.WaitGroup
	output := make(chan LogEntry, 50)
	go func() {
		wg.Wait()
		close(output)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case entry, ok := <-input:
					if !ok {
						return
					}
					if entry.StatusCode >= 400 {
						output <- entry
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return output
}

func CalculateStats(input <-chan LogEntry) Statistics {
	var totalRespTime int

	stats := Statistics{
		RequestsByIP: make(map[string]int),
	}

	for entry := range input {
		stats.TotalRequests++
		if entry.StatusCode >= 400 {
			stats.ErrorCount++
		}
		stats.RequestsByIP[entry.IP]++
		totalRespTime += entry.ResponseTime
	}
	if stats.TotalRequests == 0 {
		stats.AverageRespTime = 0
	} else {
		stats.AverageRespTime = float64(totalRespTime) / float64(stats.TotalRequests)
	}

	return stats
}
