// Package diag defines go-dbml's diagnostics: compile errors, warnings and
// infos carrying precise source ranges.
//
// Positions follow the conventions of the Monaco editor so a Diagnostic maps
// 1:1 onto a Monaco marker:
//   - Line and Column are 1-based.
//   - Column and Offset are counted in UTF-16 code units (JavaScript string
//     indices), not bytes or runes.
//   - End is exclusive.
package diag

import "fmt"

// Severity of a diagnostic.
type Severity int

const (
	// SeverityError blocks model generation: a file with any error yields no model.
	SeverityError Severity = iota
	// SeverityWarning is advisory.
	SeverityWarning
	// SeverityInfo is advisory and lower priority than a warning (e.g. a ref
	// operator that disagrees with column nullability).
	SeverityInfo
)

// String returns "error", "warning" or "info".
func (s Severity) String() string {
	switch s {
	case SeverityWarning:
		return "warning"
	case SeverityInfo:
		return "info"
	default:
		return "error"
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s Severity) MarshalText() ([]byte, error) { return []byte(s.String()), nil }

// Monaco marker severities (monaco.MarkerSeverity).
const (
	MonacoHint    = 1
	MonacoInfo    = 2
	MonacoWarning = 4
	MonacoError   = 8
)

// Monaco returns the monaco.MarkerSeverity value for s.
func (s Severity) Monaco() int {
	switch s {
	case SeverityWarning:
		return MonacoWarning
	case SeverityInfo:
		return MonacoInfo
	default:
		return MonacoError
	}
}

// Pos is a 1-based source position.
type Pos struct {
	// Offset is the 0-based offset from the start of the file in UTF-16 code units.
	Offset int `json:"offset"`
	// Line is 1-based.
	Line int `json:"line"`
	// Column is 1-based, in UTF-16 code units.
	Column int `json:"column"`
}

// Range is a half-open source range [Start, End).
type Range struct {
	Start Pos `json:"start"`
	End   Pos `json:"end"`
}

// Diagnostic is a single compile error, warning or info.
type Diagnostic struct {
	Code     Code     `json:"code"`
	Name     string   `json:"name"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Filepath string   `json:"filepath"`
	Range    Range    `json:"range"`
}

// New builds an error-severity diagnostic.
func New(code Code, message, filepath string, r Range) Diagnostic {
	return Diagnostic{Code: code, Name: code.String(), Severity: SeverityError, Message: message, Filepath: filepath, Range: r}
}

// NewWithSeverity builds a diagnostic with an explicit severity.
func NewWithSeverity(sev Severity, code Code, message, filepath string, r Range) Diagnostic {
	d := New(code, message, filepath, r)
	d.Severity = sev
	return d
}

// Error implements the error interface.
func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s: %s (%s)", d.Filepath, d.Range.Start.Line, d.Range.Start.Column, d.Severity, d.Message, d.Code)
}

// MonacoMarker mirrors monaco.editor.IMarkerData.
type MonacoMarker struct {
	Severity        int    `json:"severity"`
	Message         string `json:"message"`
	Code            string `json:"code"`
	Source          string `json:"source"`
	StartLineNumber int    `json:"startLineNumber"`
	StartColumn     int    `json:"startColumn"`
	EndLineNumber   int    `json:"endLineNumber"`
	EndColumn       int    `json:"endColumn"`
}

// Marker converts the diagnostic into a Monaco marker.
func (d Diagnostic) Marker() MonacoMarker {
	return MonacoMarker{
		Severity:        d.Severity.Monaco(),
		Message:         d.Message,
		Code:            d.Code.String(),
		Source:          "dbml",
		StartLineNumber: d.Range.Start.Line,
		StartColumn:     d.Range.Start.Column,
		EndLineNumber:   d.Range.End.Line,
		EndColumn:       d.Range.End.Column,
	}
}

// List is a slice of diagnostics with filtering helpers.
type List []Diagnostic

// Errors returns only error-severity diagnostics.
func (l List) Errors() List { return l.filter(SeverityError) }

// Warnings returns only warning-severity diagnostics.
func (l List) Warnings() List { return l.filter(SeverityWarning) }

// Infos returns only info-severity diagnostics.
func (l List) Infos() List { return l.filter(SeverityInfo) }

// HasErrors reports whether any diagnostic is an error.
func (l List) HasErrors() bool {
	for _, d := range l {
		if d.Severity == SeverityError {
			return true
		}
	}
	return false
}

// Markers converts every diagnostic into a Monaco marker.
func (l List) Markers() []MonacoMarker {
	out := make([]MonacoMarker, 0, len(l))
	for _, d := range l {
		out = append(out, d.Marker())
	}
	return out
}

func (l List) filter(s Severity) List {
	var out List
	for _, d := range l {
		if d.Severity == s {
			out = append(out, d)
		}
	}
	return out
}
