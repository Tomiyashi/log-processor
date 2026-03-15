package worker

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

func ReadLogs(ctx context.Context, filename string) (<-chan LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не получилось открыть файл %v", filename)
	}

	ch := make(chan LogEntry, 100)

	go func() {
		defer file.Close()
		defer close(ch)
		scanner := bufio.NewScanner(file)
		scanner.Scan()

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}
			line := scanner.Text()
			entry, err := ParseLogLine(line)
			if err != nil {
				log.Printf("ошибка парсинга строки '%s': %v", line, err)
				continue
			}
			ch <- entry
		}
		if err := scanner.Err(); err != nil {
			log.Printf("ошибка чтения файла: %v", err)
		}
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
					output <- entry

				case <-ctx.Done():
					return
				}
			}
		}()
	}
	return output
}

func FilterLogs(input <-chan LogEntry, minStatus int) <-chan LogEntry {
	output := make(chan LogEntry, 50)

	go func() {
		defer close(output)

		for entry := range input {
			if entry.StatusCode >= minStatus {
				output <- entry
			}
		}
	}()

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

func PrintTopIPs(RequestsByIP map[string]int, n int) {

	slice := []IpCount{}

	for ip, count := range RequestsByIP {
		slice = append(slice, IpCount{
			IP:    ip,
			Count: count,
		})
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Count > slice[j].Count
	})

	fmt.Println("\nТоп IP-адресов:")

	limit := n
	if len(slice) < limit {
		limit = len(slice)
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("  %d. %s — %d запросов\n", i+1, slice[i].IP, slice[i].Count)
	}
}
