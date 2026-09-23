package report

import "testing"

func TestSourceReport(t *testing.T) {
	report := NewSourceReport()
	report.Add(Event{Source: "alpha", Level: "error", DurationMS: 10, Bytes: 100})
	report.Add(Event{Source: "alpha", DurationMS: 20, Bytes: 50})
	report.Add(Event{Source: "beta", DurationMS: 5})
	if report.Total() != 3 || report.Unique() != 2 {
		t.Fatalf("unexpected totals")
	}
	row, ok := report.Find("alpha")
	if !ok || row.Count != 2 || row.DurationTotal != 30 || row.BytesTotal != 150 {
		t.Fatalf("unexpected bucket: %+v", row)
	}
	if report.Rows(1)[0].Key != "alpha" {
		t.Fatal("ranking failed")
	}
	merged := NewSourceReport()
	merged.Merge(report)
	if merged.Total() != 3 {
		t.Fatal("merge failed")
	}
	merged.Reset()
	if merged.Total() != 0 {
		t.Fatal("reset failed")
	}
}
