package dbml

import (
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/diag"
)

func recordsSrc(rows int) string {
	var b strings.Builder
	b.WriteString("Table t {\n  a int\n}\nRecords t(a) {\n")
	for i := 0; i < rows; i++ {
		b.WriteString("  1\n")
	}
	b.WriteString("}\n")
	return b.String()
}

func hasCode(l diag.List, c diag.Code) bool {
	for _, d := range l {
		if d.Code == c {
			return true
		}
	}
	return false
}

func TestRecordLimits(t *testing.T) {
	r := ParseFileWith(DefaultFilepath, recordsSrc(50), Options{MaxRecordRows: 10})
	if r.Database != nil || !hasCode(r.Diagnostics, diag.RECORDS_TOO_LARGE) {
		t.Fatalf("want RECORDS_TOO_LARGE, got %v", r.Diagnostics)
	}
	if n := len(r.Partial.Records[0].Values); n != 10 {
		t.Errorf("partial keeps %d rows, want 10", n)
	}
	r = ParseFileWith(DefaultFilepath, recordsSrc(50), Options{MaxRecordBytes: 64})
	if r.Database != nil || !hasCode(r.Diagnostics, diag.RECORDS_TOO_LARGE) {
		t.Fatalf("want RECORDS_TOO_LARGE for bytes, got %v", r.Diagnostics)
	}
	if r := ParseFile(DefaultFilepath, recordsSrc(50)); r.Database == nil {
		t.Fatalf("defaults must allow 50 rows: %v", r.Diagnostics)
	}
}
