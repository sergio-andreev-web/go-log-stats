package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
)

type Event struct { Level string `json:"level"`; Message string `json:"message"` }

func main() {
    if len(os.Args) != 2 { fmt.Fprintln(os.Stderr, "usage: go run . <events.jsonl>"); os.Exit(2) }
    file, err := os.Open(os.Args[1]); if err != nil { panic(err) }; defer file.Close()
    counts := map[string]int{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var event Event
        if err := json.Unmarshal(scanner.Bytes(), &event); err != nil { panic(err) }
        counts[event.Level]++
    }
    if err := scanner.Err(); err != nil { panic(err) }
    output, _ := json.MarshalIndent(counts, "", "  ")
    fmt.Println(string(output))
}
