package dbml

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/diag"
)

// referenceThreshold gates the differential regression suite: inputs under
// testdata/reference/<set>/ paired with the output of @dbml/parse 10.3.1
// (`<name>.golden.json`, generated once by scripts/ — no Node at test time).
const referenceThreshold = 0.95

const referenceDir = "testdata/reference"

// knownReferenceFailures lists cases that are allowed to differ from the
// reference (set/name -> reason). Any other failing case fails the test, so
// a single regression is caught even while the rate stays above the gate.
var knownReferenceFailures = map[string]string{}

type refDiag struct {
	Level string `json:"level"`
	Code  string `json:"code"`
	Start []int  `json:"start"`
}

type refGolden struct {
	Database any       `json:"database"`
	Diags    []refDiag `json:"diags"`
	Crash    string    `json:"crash"`
}

// stripEmpty drops null, [] and {} recursively (keeps false/0/""), the
// normalisation used by the differential harness. Array order is kept.
func stripEmpty(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, x := range t {
			x = stripEmpty(x)
			if !isEmptyJS(x) {
				out[k] = x
			}
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, x := range t {
			out[i] = stripEmpty(x)
		}
		return out
	}
	return v
}

func refDiagKeys(ds []refDiag) []string {
	var out []string
	for _, d := range ds {
		if len(d.Start) < 2 {
			out = append(out, d.Level+" "+d.Code)
			continue
		}
		out = append(out, fmt.Sprintf("%s %s %d:%d", d.Level, d.Code, d.Start[0], d.Start[1]))
	}
	sort.Strings(out)
	return out
}

func goDiagKeys(ds diag.List) []string {
	var out []string
	for _, d := range ds {
		out = append(out, fmt.Sprintf("%s %s %d:%d", d.Severity, d.Code, d.Range.Start.Line-1, d.Range.Start.Column-1))
	}
	sort.Strings(out)
	return out
}

type refCase struct {
	set, name string
	run       func() Result
}

// Result mirrors model.Result for the reference harness.
type Result = struct {
	Database any
	Diags    diag.List
}

func referenceCases(t *testing.T) []refCase {
	t.Helper()
	var cases []refCase
	sets, err := os.ReadDir(referenceDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range sets {
		dir := filepath.Join(referenceDir, s.Name())
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			set := s.Name()
			switch {
			case e.IsDir():
				root := filepath.Join(dir, e.Name())
				cases = append(cases, refCase{set, e.Name(), func() Result {
					files := map[string]string{}
					_ = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
						if err == nil && !d.IsDir() {
							b, _ := os.ReadFile(p)
							rel, _ := filepath.Rel(root, p)
							files["/"+filepath.ToSlash(rel)] = string(b)
						}
						return nil
					})
					r := ParseProject("/main.dbml", files)
					return Result{dbOrNil(r.Database), r.Diagnostics}
				}})
			case strings.HasSuffix(e.Name(), ".dbml"):
				p := filepath.Join(dir, e.Name())
				cases = append(cases, refCase{set, strings.TrimSuffix(e.Name(), ".dbml"), func() Result {
					b, _ := os.ReadFile(p)
					r := ParseFile(DefaultFilepath, string(b))
					return Result{dbOrNil(r.Database), r.Diagnostics}
				}})
			}
		}
	}
	return cases
}

func dbOrNil[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}

func TestReferenceSuite(t *testing.T) {
	cases := referenceCases(t)
	pass := 0
	perSet := map[string][2]int{}
	var fails, regressions []string
	for _, c := range cases {
		raw, err := os.ReadFile(filepath.Join(referenceDir, c.set, c.name+".golden.json"))
		if err != nil {
			t.Fatal(err)
		}
		var g refGolden
		if err := json.Unmarshal(raw, &g); err != nil {
			t.Fatal(err)
		}
		r := c.run()
		var got any
		b, _ := json.Marshal(r.Database)
		_ = json.Unmarshal(b, &got)
		want := stripEmpty(g.Database)
		gotS := stripEmpty(got)
		wk, gk := refDiagKeys(g.Diags), goDiagKeys(r.Diags)
		ok := reflect.DeepEqual(wk, gk) && reflect.DeepEqual(want, gotS)
		ps := perSet[c.set]
		ps[1]++
		if ok {
			pass++
			ps[0]++
		} else {
			why := ""
			if !reflect.DeepEqual(wk, gk) {
				why = fmt.Sprintf("diags want %v got %v", wk, gk)
			} else {
				why = firstDiff("db", want, gotS)
			}
			id := c.set + "/" + c.name
			if _, known := knownReferenceFailures[id]; !known {
				regressions = append(regressions, id)
			}
			fails = append(fails, fmt.Sprintf("%s: %s", id, trunc(why)))
		}
		perSet[c.set] = ps
	}
	rate := float64(pass) / float64(len(cases))
	var sets []string
	for s := range perSet {
		sets = append(sets, s)
	}
	sort.Strings(sets)
	for _, s := range sets {
		t.Logf("  reference/%-9s %d/%d", s, perSet[s][0], perSet[s][1])
	}
	t.Logf("@dbml/parse 10.3.1 reference suite: %d/%d exact (%.1f%%)", pass, len(cases), rate*100)
	for _, f := range fails {
		t.Logf("  FAIL %s", f)
	}
	if len(regressions) > 0 {
		t.Errorf("%d reference case(s) regressed (not in knownReferenceFailures): %v", len(regressions), regressions)
	}
	if rate < referenceThreshold {
		t.Errorf("reference suite %.1f%% is below the %.0f%% gate", rate*100, referenceThreshold*100)
	}
}
