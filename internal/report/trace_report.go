package report

import "strings"

// TraceReport aggregates events by trace.
type TraceReport struct {
	buckets map[string]*Bucket
	total   int
}

func NewTraceReport() *TraceReport {
	return &TraceReport{buckets: make(map[string]*Bucket)}
}

func (report *TraceReport) Add(event Event) {
	key := strings.TrimSpace(event.Trace)
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

func (report *TraceReport) Rows(limit int) []Bucket {
	return Rank(report.buckets, limit)
}

func (report *TraceReport) Total() int {
	return report.total
}

func (report *TraceReport) Unique() int {
	return len(report.buckets)
}

func (report *TraceReport) Find(key string) (Bucket, bool) {
	bucket, ok := report.buckets[key]
	if !ok {
		return Bucket{}, false
	}
	return *bucket, true
}

func (report *TraceReport) Top() (Bucket, bool) {
	rows := report.Rows(1)
	if len(rows) == 0 {
		return Bucket{}, false
	}
	return rows[0], true
}

func (report *TraceReport) Merge(other *TraceReport) {
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

func (report *TraceReport) Reset() {
	report.buckets = make(map[string]*Bucket)
	report.total = 0
}

func (report *TraceReport) JSON(limit int) ([]byte, error) {
	return Encode(report.Rows(limit))
}
