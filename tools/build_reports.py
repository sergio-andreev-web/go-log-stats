from pathlib import Path

root = Path(__file__).resolve().parents[1]
fields = {
    'Level':'Level', 'Service':'Service', 'Message':'Message', 'Status':'Status',
    'Method':'Method', 'Path':'Path', 'Host':'Host', 'User':'User',
    'Region':'Region', 'Environment':'Environment', 'Trace':'Trace',
    'RequestID':'RequestID', 'ErrorCode':'ErrorCode', 'Source':'Source',
    'Component':'Component', 'Version':'Version', 'Operation':'Operation',
    'Tenant':'Tenant', 'Device':'Device', 'Browser':'Browser',
    'Platform':'Platform', 'Outcome':'Outcome', 'Category':'Category',
    'Queue':'Queue', 'Worker':'Worker',
}
for name, field in fields.items():
    path = root/'internal/report'/f'{name.lower()}_report.go'
    path.write_text(f'''package report

import "strings"

// {name}Report aggregates events by {name.lower()}.
type {name}Report struct {{
    buckets map[string]*Bucket
    total int
}}

func New{name}Report() *{name}Report {{
    return &{name}Report{{buckets: make(map[string]*Bucket)}}
}}

func (report *{name}Report) Add(event Event) {{
    key := strings.TrimSpace(event.{field})
    if key == "" {{ return }}
    bucket, exists := report.buckets[key]
    if !exists {{
        bucket = &Bucket{{Key: key, DurationMin: event.DurationMS, DurationMax: event.DurationMS}}
        report.buckets[key] = bucket
    }}
    bucket.Count++
    bucket.DurationTotal += event.DurationMS
    bucket.BytesTotal += event.Bytes
    if event.DurationMS < bucket.DurationMin {{ bucket.DurationMin = event.DurationMS }}
    if event.DurationMS > bucket.DurationMax {{ bucket.DurationMax = event.DurationMS }}
    if event.Level == "error" || strings.HasPrefix(event.Status, "5") {{ bucket.Errors++ }}
    report.total++
}}

func (report *{name}Report) Rows(limit int) []Bucket {{
    return Rank(report.buckets, limit)
}}

func (report *{name}Report) Total() int {{
    return report.total
}}

func (report *{name}Report) Unique() int {{
    return len(report.buckets)
}}

func (report *{name}Report) Find(key string) (Bucket, bool) {{
    bucket, ok := report.buckets[key]
    if !ok {{ return Bucket{{}}, false }}
    return *bucket, true
}}

func (report *{name}Report) Top() (Bucket, bool) {{
    rows := report.Rows(1)
    if len(rows) == 0 {{ return Bucket{{}}, false }}
    return rows[0], true
}}

func (report *{name}Report) Merge(other *{name}Report) {{
    for _, row := range other.buckets {{
        target, ok := report.buckets[row.Key]
        if !ok {{
            copy := *row
            report.buckets[row.Key] = &copy
            continue
        }}
        target.Count += row.Count
        target.Errors += row.Errors
        target.DurationTotal += row.DurationTotal
        target.BytesTotal += row.BytesTotal
        if row.DurationMin < target.DurationMin {{ target.DurationMin = row.DurationMin }}
        if row.DurationMax > target.DurationMax {{ target.DurationMax = row.DurationMax }}
    }}
    report.total += other.total
}}

func (report *{name}Report) Reset() {{
    report.buckets = make(map[string]*Bucket)
    report.total = 0
}}

func (report *{name}Report) JSON(limit int) ([]byte, error) {{
    return Encode(report.Rows(limit))
}}
''')
    (root/'internal/report'/f'{name.lower()}_report_test.go').write_text(f'''package report

import "testing"

func Test{name}Report(t *testing.T) {{
    report := New{name}Report()
    report.Add(Event{{{field}: "alpha", {'' if name == 'Level' else 'Level: "error", '}DurationMS: 10, Bytes: 100}})
    report.Add(Event{{{field}: "alpha", DurationMS: 20, Bytes: 50}})
    report.Add(Event{{{field}: "beta", DurationMS: 5}})
    if report.Total() != 3 || report.Unique() != 2 {{ t.Fatalf("unexpected totals") }}
    row, ok := report.Find("alpha")
    if !ok || row.Count != 2 || row.DurationTotal != 30 || row.BytesTotal != 150 {{ t.Fatalf("unexpected bucket: %+v", row) }}
    if report.Rows(1)[0].Key != "alpha" {{ t.Fatal("ranking failed") }}
    merged := New{name}Report()
    merged.Merge(report)
    if merged.Total() != 3 {{ t.Fatal("merge failed") }}
    merged.Reset()
    if merged.Total() != 0 {{ t.Fatal("reset failed") }}
}}
''')
imports = '\n'.join([])
cases = '\n'.join(f'case "{name.lower()}": return New{name}Report(), nil' for name in fields)
(root/'internal/report/registry.go').write_text('''package report

import "fmt"

func NewByName(name string) (Analyzer, error) {
    switch name {
'''+cases+'''
    default: return nil, fmt.Errorf("unknown group: %s", name)
    }
}
''')
print('Generated',len(fields),'report analyzers')
