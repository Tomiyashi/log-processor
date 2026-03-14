package worker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LogEntry struct {
	Timestamp    string
	IP           string
	Method       string
	URL          string
	StatusCode   int
	ResponseTime int
}


type Statistics struct {
	TotalRequests   int
	ErrorCount      int
	RequestsByIP    map[string]int
	AverageRespTime float64
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
