package main

import (
	"strings"
	"testing"
	"time"
)

const sample = `{"time":"2026-01-01","level":"info","service":"api","message":"started"}
{"time":"2026-01-02","level":"error","service":"api","message":"timeout"}
{"time":"2026-01-03","level":"error","service":"worker","message":"timeout"}
`

func TestAnalyze(t *testing.T) {
	summary, err := Analyze(strings.NewReader(sample), Filters{}, 3)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Total != 3 || summary.ByLevel["error"] != 2 || summary.TopMessages[0].Count != 2 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}

func TestFilters(t *testing.T) {
	since, _ := time.Parse("2006-01-02", "2026-01-02")
	summary, err := Analyze(strings.NewReader(sample), Filters{Level: "error", Service: "api", Since: since}, 2)
	if err != nil || summary.Total != 1 || summary.TopMessages[0].Value != "timeout" {
		t.Fatalf("unexpected filtered result: %+v, %v", summary, err)
	}
}

func TestInvalidEventReportsLine(t *testing.T) {
	_, err := Analyze(strings.NewReader(sample+"{bad json}"), Filters{}, 5)
	if err == nil || !strings.Contains(err.Error(), "line 4") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRejectsInvalidTop(t *testing.T) {
	_, err := Analyze(strings.NewReader(sample), Filters{}, 0)
	if err == nil {
		t.Fatal("expected an error")
	}
}
