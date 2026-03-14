package worker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LogEntry struct {
	Timestamp    string // или time.Time, если будем парсить время
	IP           string // IP адрес клиента
	Method       string // HTTP метод (GET, POST и т.д.)
	URL          string // путь запроса
	StatusCode   int    // HTTP статус код (200, 404, 500)
	ResponseTime int    // время ответа в миллисекундах
}

// Statistics хранит агрегированные данные
type Statistics struct {
	TotalRequests   int            // общее количество запросов
	ErrorCount      int            // количество ошибок (статус >= 400)
	RequestsByIP    map[string]int // количество запросов с каждого IP
	AverageRespTime float64        // среднее время ответа
}

func ParseLogLine(line string) (LogEntry, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 6 {
		return LogEntry{}, errors.New("неверное количество полей")
	}

	intStatus, err := strconv.Atoi(fields[4])
	if err != nil {
		return LogEntry{}, fmt.Errorf("не удалось конвертировать StatusCode: %s", fields[4])
	}
	intResponse, err := strconv.Atoi(fields[5])
	if err != nil {
		return LogEntry{}, fmt.Errorf("не удалось конвертировать ResponseTime: %s", fields[5])
	}

	entry := LogEntry{
		Timestamp:    fields[0],
		IP:           fields[1],
		Method:       fields[2],
		URL:          fields[3],
		StatusCode:   intStatus,
		ResponseTime: intResponse,
	}

	return entry, nil
}
