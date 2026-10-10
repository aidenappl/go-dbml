package dbml

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aidenappl/go-dbml/diag"
)

// conformanceThreshold is the MVP gate: the share of Holistics interpreter
// fixtures that must match exactly (after the same normalisation the
// reference snapshot tests apply).
const conformanceThreshold = 0.95

const fixtureDir = "testdata/holistics/interpreter"

type expectedDiag struct {
	Level string `json:"level"`
	Code  string `json:"code"`
	Node  *struct {
		Context struct {
			ID string `json:"id"`
		} `json:"context"`
	} `json:"node"`
	Token *struct {
		Context struct {
			ID string `json:"id"`
		} `json:"context"`
	} `json:"token"`
}

type expectedFixture struct {
	Database any            `json:"database"`
	Errors   []expectedDiag `json:"errors"`
	Warnings []expectedDiag `json:"warnings"`
	Infos    []expectedDiag `json:"infos"`
}

var ctxPos = regexp.MustCompile(`\[L(-?\d+):C(-?\d+),`)

func (d expectedDiag) key() string {
	id := ""
	if d.Node != nil {
		id = d.Node.Context.ID
	} else if d.Token != nil {
		id = d.Token.Context.ID
	}
	m := ctxPos.FindStringSubmatch(id)
	if m == nil {
		return d.Level + " " + d.Code + " ?"
	}
	return fmt.Sprintf("%s %s L%s:C%s", d.Level, d.Code, m[1], m[2])
}

func diagKey(d diag.Diagnostic) string {
	return fmt.Sprintf("%s %s L%d:C%d", d.Severity, d.Code, d.Range.Start.Line-1, d.Range.Start.Column-1)
}

func sortedKeys(xs []string) []string { sort.Strings(xs); return xs }

// ---------------------------------------------------------------- snapshot normalisation
//
// Mirrors toSnapshot/sortObject/sortArray from
// packages/dbml-parse/__tests__/utils/testHelpers.ts: object keys are sorted,
// empty values (null, [], {}) are dropped, and arrays are ordered by a
// stable kind/position ranking (schema elements by token start offset).

func isEmptyJS(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case []any:
		return len(t) == 0
	case map[string]any:
		return len(t) == 0
	}
	return false
}

func normalize(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, val := range t {
			nv := normalize(val)
			if isEmptyJS(nv) {
				continue
			}
			out[k] = nv
		}
		return out
	case []any:
		cp := append([]any{}, t...)
		sort.SliceStable(cp, func(i, j int) bool { return compareJS(cp[i], cp[j]) < 0 })
		out := make([]any, len(cp))
		for i, e := range cp {
			out[i] = normalize(e)
		}
		return out
	}
	return v
}

func interRank(v any) int {
	switch t := v.(type) {
	case nil:
		return -2
	case string:
		return 0
	case float64:
		return 1
	case bool:
		return 2
	case map[string]any:
		if tok, ok := t["token"]; ok && tok != nil {
			return 9
		}
	}
	return 1000
}

func intraRank(v any) any {
	switch t := v.(type) {
	case nil:
		return 0.0
	case string, float64:
		return t
	case bool:
		if t {
			return 1.0
		}
		return 0.0
	case map[string]any:
		if tok, ok := t["token"].(map[string]any); ok {
			if st, ok := tok["start"].(map[string]any); ok {
				if off, ok := st["offset"].(float64); ok {
					return off
				}
			}
			return 0.0
		}
		keys := make([]string, 0, len(t))
		for k, val := range t {
			if isEmptyJS(val) {
				if _, isStr := val.(string); !isStr {
					continue
				}
			}
			keys = append(keys, k)
		}
		if len(keys) == 0 {
			return 0.0
		}
		sort.Strings(keys)
		return intraRank(t[keys[0]])
	case []any:
		keys := make([]string, 0, len(t))
		for i, val := range t {
			if isEmptyJS(val) {
				continue
			}
			keys = append(keys, strconv.Itoa(i))
		}
		if len(keys) == 0 {
			return 0.0
		}
		sort.Strings(keys)
		i, _ := strconv.Atoi(keys[0])
		return intraRank(t[i])
	}
	return 0.0
}

func compareRank(a, b any) int {
	fa, aNum := a.(float64)
	fb, bNum := b.(float64)
	sa, aStr := a.(string)
	sb, bStr := b.(string)
	switch {
	case aNum && bNum:
		switch {
		case math.IsNaN(fa) && math.IsNaN(fb):
			return 0
		case fa < fb:
			return -1
		case fa > fb:
			return 1
		}
		return 0
	case aStr && bStr:
		return strings.Compare(sa, sb)
	case aStr && bNum:
		n, err := strconv.ParseFloat(sa, 64)
		if err != nil {
			return 0
		}
		return compareRank(n, fb)
	case aNum && bStr:
		n, err := strconv.ParseFloat(sb, 64)
		if err != nil {
			return 0
		}
		return compareRank(fa, n)
	}
	return 0
}

func compareJS(a, b any) int {
	if d := interRank(a) - interRank(b); d != 0 {
		return d
	}
	return compareRank(intraRank(a), intraRank(b))
}

// ---------------------------------------------------------------- test

type fixtureResult struct {
	name     string
	pass     bool
	infosOK  bool
	category string
	detail   string
}

func runFixture(t *testing.T, dir, name string) fixtureResult {
	t.Helper()
	src, err := os.ReadFile(filepath.Join(dir, name+".in.dbml"))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(dir, name+".out.json"))
	if err != nil {
		t.Fatal(err)
	}
	var exp expectedFixture
	if err := json.Unmarshal(raw, &exp); err != nil {
		t.Fatal(err)
	}
	res := model_compile(string(src))
	gotErrs := res.Diagnostics.Errors()

	var infosWant, infosGot []string
	for _, d := range exp.Infos {
		infosWant = append(infosWant, d.key())
	}
	for _, d := range res.Diagnostics.Infos() {
		infosGot = append(infosGot, diagKey(d))
	}
	infosOK := reflect.DeepEqual(sortedKeys(infosWant), sortedKeys(infosGot))

	if len(exp.Errors) > 0 {
		var want, got []string
		for _, d := range exp.Errors {
			want = append(want, d.key())
		}
		for _, d := range gotErrs {
			got = append(got, diagKey(d))
		}
		sortedKeys(want)
		sortedKeys(got)
		if res.Database == nil && reflect.DeepEqual(want, got) {
			return fixtureResult{name: name, pass: true, infosOK: true}
		}
		return fixtureResult{name: name, category: "diagnostics mismatch", infosOK: true,
			detail: fmt.Sprintf("want %v\n      got  %v", want, got)}
	}
	if len(gotErrs) > 0 {
		var got []string
		for _, d := range gotErrs {
			got = append(got, diagKey(d)+" "+d.Message)
		}
		return fixtureResult{name: name, category: "unexpected errors", infosOK: infosOK, detail: strings.Join(got, "; ")}
	}
	gotJSON, err := json.Marshal(res.Database)
	if err != nil {
		t.Fatal(err)
	}
	var gotAny any
	if err := json.Unmarshal(gotJSON, &gotAny); err != nil {
		t.Fatal(err)
	}
	want := normalize(exp.Database)
	got := normalize(gotAny)
	if reflect.DeepEqual(want, got) {
		return fixtureResult{name: name, pass: true, infosOK: infosOK}
	}
	return fixtureResult{name: name, category: "model mismatch", infosOK: infosOK, detail: firstDiff("database", want, got)}
}

func model_compile(src string) struct {
	Database    any
	Diagnostics diag.List
} {
	r := ParseFile(DefaultFilepath, src)
	out := struct {
		Database    any
		Diagnostics diag.List
	}{Diagnostics: r.Diagnostics}
	if r.Database != nil {
		out.Database = r.Database
	}
	return out
}

func firstDiff(path string, want, got any) string {
	switch w := want.(type) {
	case map[string]any:
		g, ok := got.(map[string]any)
		if !ok {
			return fmt.Sprintf("%s: want object, got %T", path, got)
		}
		var keys []string
		for k := range w {
			keys = append(keys, k)
		}
		for k := range g {
			if _, ok := w[k]; !ok {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			if !reflect.DeepEqual(w[k], g[k]) {
				return firstDiff(path+"."+k, w[k], g[k])
			}
		}
	case []any:
		g, ok := got.([]any)
		if !ok {
			return fmt.Sprintf("%s: want array, got %T", path, got)
		}
		if len(w) != len(g) {
			return fmt.Sprintf("%s: want %d items, got %d", path, len(w), len(g))
		}
		for i := range w {
			if !reflect.DeepEqual(w[i], g[i]) {
				return firstDiff(fmt.Sprintf("%s[%d]", path, i), w[i], g[i])
			}
		}
	}
	wj, _ := json.Marshal(want)
	gj, _ := json.Marshal(got)
	return fmt.Sprintf("%s: want %s, got %s", path, trunc(string(wj)), trunc(string(gj)))
}

func trunc(s string) string {
	if len(s) > 300 {
		return s[:300] + "…"
	}
	return s
}

func TestHolisticsConformance(t *testing.T) {
	type fx struct{ dir, name string }
	var names []fx
	for _, dir := range []string{fixtureDir, "testdata/holistics/nan"} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".in.dbml") {
				names = append(names, fx{dir, strings.TrimSuffix(e.Name(), ".in.dbml")})
			}
		}
	}
	pass, infos := 0, 0
	var failures []fixtureResult
	for _, n := range names {
		r := runFixture(t, n.dir, n.name)
		if r.pass {
			pass++
		} else {
			failures = append(failures, r)
		}
		if r.infosOK {
			infos++
		}
	}
	rate := float64(pass) / float64(len(names))
	t.Logf("Holistics interpreter+nan conformance: %d/%d fixtures pass (%.1f%%); info diagnostics match in %d/%d",
		pass, len(names), rate*100, infos, len(names))
	for _, f := range failures {
		t.Logf("  FAIL %-45s [%s] %s", f.name, f.category, f.detail)
	}
	if rate < conformanceThreshold {
		t.Errorf("conformance %.1f%% is below the %.0f%% MVP gate", rate*100, conformanceThreshold*100)
	}
}

// TestInfiniteLoopSnapshots mirrors the reference's infinite_loops suite:
// each input must finish compiling (the reference allows 10s).
func TestInfiniteLoopSnapshots(t *testing.T) {
	dir := "testdata/holistics/infinite_loops"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		src, _ := os.ReadFile(filepath.Join(dir, e.Name()))
		done := make(chan struct{})
		go func() { ParseFile(DefaultFilepath, string(src)); close(done) }()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("%s did not finish within 10s", e.Name())
		}
	}
}
