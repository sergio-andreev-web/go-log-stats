package report

import (
	"encoding/json"
	"sort"
)

type Event struct {
	Time        string `json:"time"`
	Level       string `json:"level"`
	Message     string `json:"message"`
	Service     string `json:"service,omitempty"`
	Status      string `json:"status,omitempty"`
	Method      string `json:"method,omitempty"`
	Path        string `json:"path,omitempty"`
	Host        string `json:"host,omitempty"`
	User        string `json:"user,omitempty"`
	Region      string `json:"region,omitempty"`
	Environment string `json:"environment,omitempty"`
	Trace       string `json:"trace,omitempty"`
	RequestID   string `json:"request_id,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
	Source      string `json:"source,omitempty"`
	Component   string `json:"component,omitempty"`
	Version     string `json:"version,omitempty"`
	Operation   string `json:"operation,omitempty"`
	Tenant      string `json:"tenant,omitempty"`
	Device      string `json:"device,omitempty"`
	Browser     string `json:"browser,omitempty"`
	Platform    string `json:"platform,omitempty"`
	Outcome     string `json:"outcome,omitempty"`
	Category    string `json:"category,omitempty"`
	Queue       string `json:"queue,omitempty"`
	Worker      string `json:"worker,omitempty"`
	DurationMS  int64  `json:"duration_ms,omitempty"`
	Bytes       int64  `json:"bytes,omitempty"`
}

type Bucket struct {
	Key           string `json:"key"`
	Count         int    `json:"count"`
	Errors        int    `json:"errors"`
	DurationTotal int64  `json:"duration_total_ms"`
	DurationMin   int64  `json:"duration_min_ms"`
	DurationMax   int64  `json:"duration_max_ms"`
	BytesTotal    int64  `json:"bytes_total"`
}

func (bucket Bucket) MeanDuration() float64 {
	if bucket.Count == 0 {
		return 0
	}
	return float64(bucket.DurationTotal) / float64(bucket.Count)
}

func Rank(buckets map[string]*Bucket, limit int) []Bucket {
	rows := make([]Bucket, 0, len(buckets))
	for _, item := range buckets {
		rows = append(rows, *item)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Count == rows[j].Count {
			return rows[i].Key < rows[j].Key
		}
		return rows[i].Count > rows[j].Count
	})
	if limit > 0 && limit < len(rows) {
		return rows[:limit]
	}
	return rows
}

func Encode(rows []Bucket) ([]byte, error) { return json.MarshalIndent(rows, "", "  ") }

type Analyzer interface {
	Add(Event)
	Rows(int) []Bucket
	Total() int
	Reset()
}
