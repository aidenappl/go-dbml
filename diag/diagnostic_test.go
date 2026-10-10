package diag

import "testing"

func TestMarker(t *testing.T) {
	tests := []struct {
		name string
		sev  Severity
		want int
	}{
		{"error", SeverityError, MonacoError},
		{"warning", SeverityWarning, MonacoWarning},
		{"info", SeverityInfo, MonacoInfo},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewWithSeverity(tt.sev, BINDING_ERROR, "msg", "/main.dbml", Range{Start: Pos{Line: 2, Column: 3}, End: Pos{Line: 2, Column: 9}})
			m := d.Marker()
			if m.Severity != tt.want || m.StartLineNumber != 2 || m.StartColumn != 3 || m.EndColumn != 9 || m.Code != "BINDING_ERROR" {
				t.Fatalf("bad marker %+v", m)
			}
		})
	}
}

func TestListFilters(t *testing.T) {
	l := List{
		New(UNEXPECTED_TOKEN, "a", "/f", Range{}),
		NewWithSeverity(SeverityInfo, INVALID_REF_RELATIONSHIP, "b", "/f", Range{}),
	}
	if !l.HasErrors() || len(l.Errors()) != 1 || len(l.Infos()) != 1 || len(l.Warnings()) != 0 || len(l.Markers()) != 2 {
		t.Fatal("bad filters")
	}
	if INVALID_REF_RELATIONSHIP != 3038 || BINDING_ERROR != 4000 || DUPLICATE_METADATA_FIELD != 5012 {
		t.Fatalf("codes drifted from @dbml/parse: %d %d %d", INVALID_REF_RELATIONSHIP, BINDING_ERROR, DUPLICATE_METADATA_FIELD)
	}
}
