package dbml

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/diag"
	"github.com/aidenappl/go-dbml/lexer"
	"github.com/aidenappl/go-dbml/model"
)

// diagnosticsThreshold gates the lexer/parser/validator/binder snapshot
// suites, compared on diagnostics (level, code, start position) and, for the
// lexer, on the token stream.
const diagnosticsThreshold = 0.95

type diagSuite struct {
	name string
	run  func(src string) diag.List
}

var diagSuites = []diagSuite{
	{"parser", func(src string) diag.List { return diag.List(ParseSyntax(src).Errors) }},
	{"validator", func(src string) diag.List { return model.Validate(src, DefaultFilepath) }},
	{"binder", func(src string) diag.List { return ParseFile(DefaultFilepath, src).Diagnostics }},
}

func fixtureNames(t *testing.T, dir string) []string {
	t.Helper()
	es, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range es {
		if strings.HasSuffix(e.Name(), ".in.dbml") {
			out = append(out, strings.TrimSuffix(e.Name(), ".in.dbml"))
		}
	}
	sort.Strings(out)
	return out
}

func readFixture(t *testing.T, dir, name string, v any) string {
	t.Helper()
	src, err := os.ReadFile(filepath.Join(dir, name+".in.dbml"))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, name+".out.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		t.Fatal(err)
	}
	return string(src)
}

func TestHolisticsDiagnosticsConformance(t *testing.T) {
	total, passed := 0, 0
	for _, s := range diagSuites {
		dir := filepath.Join("testdata/holistics", s.name)
		names := fixtureNames(t, dir)
		pass := 0
		for _, n := range names {
			var exp expectedFixture
			src := readFixture(t, dir, n, &exp)
			var want, got []string
			for _, d := range exp.Errors {
				want = append(want, d.key())
			}
			for _, d := range s.run(src).Errors() {
				got = append(got, diagKey(d))
			}
			sort.Strings(want)
			sort.Strings(got)
			if reflect.DeepEqual(want, got) {
				pass++
				continue
			}
			t.Logf("  FAIL %s/%s\n      want %v\n      got  %v", s.name, n, want, got)
		}
		t.Logf("Holistics %s diagnostics conformance: %d/%d", s.name, pass, len(names))
		total += len(names)
		passed += pass
	}
	lp, ln := lexerConformance(t)
	total += ln
	passed += lp
	rate := float64(passed) / float64(total)
	t.Logf("Holistics lexer/parser/validator/binder conformance: %d/%d (%.1f%%)", passed, total, rate*100)
	if rate < diagnosticsThreshold {
		t.Errorf("diagnostics conformance %.1f%% is below the %.0f%% gate", rate*100, diagnosticsThreshold*100)
	}
}

type expectedToken struct {
	Context struct {
		ID string `json:"id"`
	} `json:"context"`
	Value          string `json:"value"`
	LeadingTrivia  string `json:"leadingTrivia"`
	TrailingTrivia string `json:"trailingTrivia"`
}

var tokID = regexp.MustCompile(`(?s)^token@(<[^>]+>)@.*@\[L(\d+):C(\d+), L(\d+):C(\d+)\]$`)

func lexerConformance(t *testing.T) (int, int) {
	dir := "testdata/holistics/lexer"
	names := fixtureNames(t, dir)
	pass := 0
	for _, n := range names {
		var exp struct {
			Tokens []expectedToken `json:"tokens"`
			Errors []expectedDiag  `json:"errors"`
		}
		src := readFixture(t, dir, n, &exp)
		toks, errs := lexer.Lex(src, DefaultFilepath)
		var want, got []string
		for _, et := range exp.Tokens {
			m := tokID.FindStringSubmatch(et.Context.ID)
			if m == nil {
				want = append(want, "?"+et.Context.ID)
				continue
			}
			want = append(want, fmt.Sprintf("%s %q L%s:C%s-L%s:C%s lead=%q trail=%q", m[1], et.Value, m[2], m[3], m[4], m[5], et.LeadingTrivia, et.TrailingTrivia))
		}
		for _, tk := range toks {
			got = append(got, fmt.Sprintf("%s %q L%d:C%d-L%d:C%d lead=%q trail=%q", tk.Kind, tk.Value, tk.Start.Line, tk.Start.Column, tk.End.Line, tk.End.Column, trivia(tk.LeadingTrivia()), trivia(tk.TrailingTrivia())))
		}
		var wantErr, gotErr []string
		for _, d := range exp.Errors {
			wantErr = append(wantErr, d.key())
		}
		for _, d := range errs {
			gotErr = append(gotErr, diagKey(d))
		}
		sort.Strings(wantErr)
		sort.Strings(gotErr)
		if reflect.DeepEqual(want, got) && reflect.DeepEqual(wantErr, gotErr) {
			pass++
			continue
		}
		t.Logf("  FAIL lexer/%s", n)
		for i := range want {
			if i >= len(got) || want[i] != got[i] {
				g := "<none>"
				if i < len(got) {
					g = got[i]
				}
				t.Logf("      token %d want %s\n               got  %s", i, want[i], g)
				break
			}
		}
		if !reflect.DeepEqual(wantErr, gotErr) {
			t.Logf("      errors want %v got %v", wantErr, gotErr)
		}
	}
	t.Logf("Holistics lexer conformance: %d/%d", pass, len(names))
	return pass, len(names)
}

func trivia(ts []*lexer.Token) string {
	var sb strings.Builder
	for _, t := range ts {
		sb.WriteString(t.Value)
	}
	return sb.String()
}
