package report

import "strings"

// RequestIDReport aggregates events by requestid.
type RequestIDReport struct {
	buckets map[string]*Bucket
	total   int
}

func NewRequestIDReport() *RequestIDReport {
	return &RequestIDReport{buckets: make(map[string]*Bucket)}
}

func (report *RequestIDReport) Add(event Event) {
	key := strings.TrimSpace(event.RequestID)
	if key == "" {
		return
	}
	bucket, exists := report.buckets[key]
	if !exists {
		bucket = &Bucket{Key: key, DurationMin: event.DurationMS, DurationMax: event.DurationMS}
		report.buckets[key] = bucket
	}
	bucket.Count++
	bucket.DurationTotal += event.DurationMS
	bucket.BytesTotal += event.Bytes
	if event.DurationMS < bucket.DurationMin {
		bucket.DurationMin = event.DurationMS
	}
	if event.DurationMS > bucket.DurationMax {
		bucket.DurationMax = event.DurationMS
	}
	if event.Level == "error" || strings.HasPrefix(event.Status, "5") {
		bucket.Errors++
	}
	report.total++
}

func (report *RequestIDReport) Rows(limit int) []Bucket {
	return Rank(report.buckets, limit)
}

func (report *RequestIDReport) Total() int {
	return report.total
}

func (report *RequestIDReport) Unique() int {
	return len(report.buckets)
}

func (report *RequestIDReport) Find(key string) (Bucket, bool) {
	bucket, ok := report.buckets[key]
	if !ok {
		return Bucket{}, false
	}
	return *bucket, true
}

func (report *RequestIDReport) Top() (Bucket, bool) {
	rows := report.Rows(1)
	if len(rows) == 0 {
		return Bucket{}, false
	}
	return rows[0], true
}

func (report *RequestIDReport) Merge(other *RequestIDReport) {
	for _, row := range other.buckets {
		target, ok := report.buckets[row.Key]
		if !ok {
			copy := *row
			report.buckets[row.Key] = &copy
			continue
		}
		target.Count += row.Count
		target.Errors += row.Errors
		target.DurationTotal += row.DurationTotal
		target.BytesTotal += row.BytesTotal
		if row.DurationMin < target.DurationMin {
			target.DurationMin = row.DurationMin
		}
		if row.DurationMax > target.DurationMax {
			target.DurationMax = row.DurationMax
		}
	}
	report.total += other.total
}

func (report *RequestIDReport) Reset() {
	report.buckets = make(map[string]*Bucket)
	report.total = 0
}

func (report *RequestIDReport) JSON(limit int) ([]byte, error) {
	return Encode(report.Rows(limit))
}
