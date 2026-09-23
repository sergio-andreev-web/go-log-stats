package report

import "strings"

// ComponentReport aggregates events by component.
type ComponentReport struct {
	buckets map[string]*Bucket
	total   int
}

func NewComponentReport() *ComponentReport {
	return &ComponentReport{buckets: make(map[string]*Bucket)}
}

func (report *ComponentReport) Add(event Event) {
	key := strings.TrimSpace(event.Component)
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

func (report *ComponentReport) Rows(limit int) []Bucket {
	return Rank(report.buckets, limit)
}

func (report *ComponentReport) Total() int {
	return report.total
}

func (report *ComponentReport) Unique() int {
	return len(report.buckets)
}

func (report *ComponentReport) Find(key string) (Bucket, bool) {
	bucket, ok := report.buckets[key]
	if !ok {
		return Bucket{}, false
	}
	return *bucket, true
}

func (report *ComponentReport) Top() (Bucket, bool) {
	rows := report.Rows(1)
	if len(rows) == 0 {
		return Bucket{}, false
	}
	return rows[0], true
}

func (report *ComponentReport) Merge(other *ComponentReport) {
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

func (report *ComponentReport) Reset() {
	report.buckets = make(map[string]*Bucket)
	report.total = 0
}

func (report *ComponentReport) JSON(limit int) ([]byte, error) {
	return Encode(report.Rows(limit))
}
