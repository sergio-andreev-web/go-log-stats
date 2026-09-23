package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/sergio-andreev-web/go-log-stats/internal/report"
)

type Event = report.Event

type Filters struct {
	Level   string
	Service string
	Since   time.Time
	Until   time.Time
	Group   string
}

type Count struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

type Summary struct {
	Total       int             `json:"total"`
	ByLevel     map[string]int  `json:"by_level"`
	ByService   map[string]int  `json:"by_service"`
	TopMessages []Count         `json:"top_messages"`
	Grouped     []report.Bucket `json:"grouped,omitempty"`
}

func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time %q", value)
}

func ParseEvent(line []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(line, &event); err != nil {
		return event, err
	}
	event.Level = strings.ToLower(strings.TrimSpace(event.Level))
	event.Service = strings.TrimSpace(event.Service)
	if event.Level == "" || strings.TrimSpace(event.Message) == "" {
		return event, fmt.Errorf("level and message are required")
	}
	if _, err := parseTime(event.Time); err != nil {
		return event, err
	}
	return event, nil
}

func (filters Filters) matches(event Event) bool {
	if filters.Level != "" && filters.Level != event.Level {
		return false
	}
	if filters.Service != "" && filters.Service != event.Service {
		return false
	}
	eventTime, _ := parseTime(event.Time)
	if !filters.Since.IsZero() && (eventTime.IsZero() || eventTime.Before(filters.Since)) {
		return false
	}
	if !filters.Until.IsZero() && (eventTime.IsZero() || eventTime.After(filters.Until)) {
		return false
	}
	return true
}

func topCounts(counts map[string]int, limit int) []Count {
	ranked := make([]Count, 0, len(counts))
	for value, count := range counts {
		ranked = append(ranked, Count{value, count})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Count == ranked[j].Count {
			return ranked[i].Value < ranked[j].Value
		}
		return ranked[i].Count > ranked[j].Count
	})
	if limit < len(ranked) {
		ranked = ranked[:limit]
	}
	return ranked
}

func Analyze(reader io.Reader, filters Filters, top int) (Summary, error) {
	summary := Summary{ByLevel: map[string]int{}, ByService: map[string]int{}, TopMessages: []Count{}}
	if top < 1 || top > 100 {
		return summary, fmt.Errorf("top must be between 1 and 100")
	}
	var analyzer report.Analyzer
	if filters.Group != "" {
		var err error
		analyzer, err = report.NewByName(filters.Group)
		if err != nil {
			return summary, err
		}
	}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	messages := map[string]int{}
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Bytes()
		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}
		event, err := ParseEvent(line)
		if err != nil {
			return summary, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if !filters.matches(event) {
			continue
		}
		summary.Total++
		summary.ByLevel[event.Level]++
		if event.Service != "" {
			summary.ByService[event.Service]++
		}
		messages[event.Message]++
		if analyzer != nil {
			analyzer.Add(event)
		}
	}
	if err := scanner.Err(); err != nil {
		return summary, err
	}
	summary.TopMessages = topCounts(messages, top)
	if analyzer != nil {
		summary.Grouped = analyzer.Rows(top)
	}
	return summary, nil
}
