package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	level := flag.String("level", "", "filter by level")
	service := flag.String("service", "", "filter by service")
	sinceText := flag.String("since", "", "earliest date (YYYY-MM-DD or RFC3339)")
	untilText := flag.String("until", "", "latest date (YYYY-MM-DD or RFC3339)")
	top := flag.Int("top", 5, "number of top messages")
	jsonOutput := flag.Bool("json", false, "print JSON")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run . [options] <events.jsonl>")
		os.Exit(2)
	}
	since, err := parseTime(*sinceText)
	if err != nil {
		fail(err)
	}
	until, err := parseTime(*untilText)
	if err != nil {
		fail(err)
	}
	if !since.IsZero() && !until.IsZero() && since.After(until) {
		fail(fmt.Errorf("since must be before until"))
	}
	file, err := os.Open(flag.Arg(0))
	if err != nil {
		fail(err)
	}
	defer file.Close()
	summary, err := Analyze(file, Filters{strings.ToLower(*level), *service, since, until}, *top)
	if err != nil {
		fail(err)
	}
	if *jsonOutput {
		encoded, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			fail(err)
		}
		fmt.Println(string(encoded))
		return
	}
	fmt.Printf("Events: %d\n", summary.Total)
	for level, count := range summary.ByLevel {
		fmt.Printf("%-8s %d\n", level, count)
	}
	fmt.Println("Top messages:")
	for _, item := range summary.TopMessages {
		fmt.Printf("%4d %s\n", item.Count, item.Value)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

