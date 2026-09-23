package report

import "strings"

// OutcomeReport aggregates events by outcome.
type OutcomeReport struct {
	buckets map[string]*Bucket
	total   int
}

func NewOutcomeReport() *OutcomeReport {
	return &OutcomeReport{buckets: make(map[string]*Bucket)}
}

func (report *OutcomeReport) Add(event Event) {
	key := strings.TrimSpace(event.Outcome)
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

func (report *OutcomeReport) Rows(limit int) []Bucket {
	return Rank(report.buckets, limit)
}

func (report *OutcomeReport) Total() int {
	return report.total
}

func (report *OutcomeReport) Unique() int {
	return len(report.buckets)
}

func (report *OutcomeReport) Find(key string) (Bucket, bool) {
	bucket, ok := report.buckets[key]
	if !ok {
		return Bucket{}, false
	}
	return *bucket, true
}

func (report *OutcomeReport) Top() (Bucket, bool) {
	rows := report.Rows(1)
	if len(rows) == 0 {
		return Bucket{}, false
	}
	return rows[0], true
}

func (report *OutcomeReport) Merge(other *OutcomeReport) {
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

func (report *OutcomeReport) Reset() {
	report.buckets = make(map[string]*Bucket)
	report.total = 0
}

func (report *OutcomeReport) JSON(limit int) ([]byte, error) {
	return Encode(report.Rows(limit))
}
