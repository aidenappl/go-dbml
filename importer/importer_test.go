package importer_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/importer"
)

// importThresholds gate the exact-match rate (trailing whitespace
// normalised) against @dbml/core 10.3.1 importer.import goldens in
// testdata/import/<format>/ (scripts/reference/import.mjs).
var importThresholds = map[importer.Format]float64{
	importer.MySQL: 0.90, importer.Postgres: 0.90, importer.MSSQL: 0.90, importer.SchemaRB: 0.90,
}

// knownImportFailures lists goldens allowed to differ (format/set/name ->
// reason). Anything else failing is a regression.
// exprSyntax: the reference's ANTLR grammar rejects (and fails the whole
// import on) expression syntax foreign to the dialect inside CHECK /
// DEFAULT / index expressions or reserved words used as names; go-dbml
// does not validate expression syntax and imports these.
const exprSyntax = "reference rejects foreign expression syntax or reserved words; go-dbml imports them"

var knownImportFailures = map[string]string{
	"mysql/hand/h07_mariadb": "lenient: MariaDB CREATE OR REPLACE TABLE is imported; the reference grammar rejects it",
	"postgres/exported/core__importer__mssql_importer__output__check_constraint.out":    exprSyntax,
	"postgres/exported/core__importer__mssql_importer__output__checks.out":              exprSyntax,
	"postgres/exported/core__importer__mssql_importer__output__general_schema.out":      exprSyntax,
	"postgres/exported/core__importer__mssql_importer__output__multiple_schema.out":     exprSyntax,
	"postgres/exported/core__importer__mysql_importer__output__checks.out":              exprSyntax,
	"postgres/exported/core__importer__snowflake_importer__output__create_table.out":    exprSyntax,
	"postgres/exported/core__model_exporter__dbml_exporter__output__default_tables.out": exprSyntax,
	"postgres/exported/core__model_exporter__dbml_exporter__output__index_tables.out":   exprSyntax,
	"postgres/exported/fixtures__default_tables.in":                                     exprSyntax,
	"mssql/exported/core__importer__mysql_importer__output__checks.out":                 exprSyntax,
	"mssql/exported/core__importer__oracle_importer__output__alter_table_add_check.out": exprSyntax,
	"mssql/exported/core__importer__oracle_importer__output__column_def_checks.out":     exprSyntax,
	"mssql/exported/core__importer__oracle_importer__output__in_table_checks.out":       exprSyntax,
	"mssql/exported/core__importer__snowflake_importer__output__create_table.out":       exprSyntax,
	"mssql/exported/core__model_exporter__dbml_exporter__output__checks.out":            exprSyntax,
	"mssql/exported/core__model_exporter__dbml_exporter__output__default_tables.out":    exprSyntax,
	"mssql/exported/fixtures__column_caller_type.in":                                    exprSyntax,
	"mssql/exported/fixtures__default_tables.in":                                        exprSyntax,
	"mssql/exported/hand4__x09_quotes_names":                                            exprSyntax,
}

func normalize(s string) string {
	lines := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.TrimRight(strings.Join(lines, "\n"), "\n")
}

func importInput(f importer.Format, set, name string) string {
	hc := "../testdata/holistics/dbml-core"
	var cands []string
	switch set {
	case "importer":
		cands = []string{filepath.Join(hc, "importer", string(f)+"_importer", "input", name+".in.sql")}
	case "parser":
		cands = []string{filepath.Join(hc, "parser", string(f)+"-parse", "input", name+".in.sql"),
			filepath.Join(hc, "parser", string(f)+"-parse", "input", name+".in.rb")}
	case "exported":
		eset, ename, _ := strings.Cut(name, "__")
		cands = []string{filepath.Join("../testdata/export", string(f), eset, ename+".sql")}
	default:
		cands = []string{filepath.Join("../testdata/import/inputs", string(f), name+".sql"),
			filepath.Join("../testdata/import/inputs", string(f), name+".rb")}
	}
	for _, c := range cands {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return cands[0]
}

func TestImportGoldens(t *testing.T) {
	for _, f := range importer.Formats {
		t.Run(string(f), func(t *testing.T) {
			var files []string
			_ = filepath.WalkDir(filepath.Join("../testdata/import", string(f)), func(p string, e os.DirEntry, err error) error {
				if err == nil && !e.IsDir() && (strings.HasSuffix(p, ".dbml") || strings.HasSuffix(p, ".err")) {
					files = append(files, p)
				}
				return nil
			})
			sort.Strings(files)
			if len(files) == 0 {
				t.Skip("no goldens")
			}
			pass := 0
			var fails, regressions, cleanFails []string
			cleanTotal := 0
			for _, gp := range files {
				set := filepath.Base(filepath.Dir(gp))
				ext := filepath.Ext(gp)
				name := strings.TrimSuffix(filepath.Base(gp), ext)
				src, err := os.ReadFile(importInput(f, set, name))
				if err != nil {
					t.Fatal(err)
				}
				want, _ := os.ReadFile(gp)
				got, err := importer.CoreDBML(string(src), f)
				id := fmt.Sprintf("%s/%s/%s", f, set, name)
				ok, why := false, ""
				switch {
				case ext == ".err":
					ok = err != nil
					why = fmt.Sprintf("want error %q, got none", strings.TrimSpace(string(want)))
				case err != nil:
					why = "unexpected error: " + err.Error()
				default:
					ok = normalize(got) == normalize(string(want))
					if !ok {
						why = firstLineDiff(normalize(string(want)), normalize(got))
					}
				}
				// clean mode (the product default) must import and compile
				cleanTotal++
				if r, err := importer.ImportWith(string(src), f, importer.ImportOptions{IncludeRecords: true}); err != nil {
					cleanFails = append(cleanFails, id+": "+err.Error())
				} else if r.Database == nil {
					cleanFails = append(cleanFails, id+": "+fmt.Sprint(r.Diagnostics))
				}
				if ok {
					pass++
					continue
				}
				fails = append(fails, id+": "+why)
				if _, known := knownImportFailures[id]; !known {
					regressions = append(regressions, id)
				}
			}
			rate := float64(pass) / float64(len(files))
			t.Logf("@dbml/core 10.3.1 %s importer: %d/%d exact (%.1f%%)", f, pass, len(files), rate*100)
			for _, x := range fails {
				t.Logf("  FAIL %s", x)
			}
			cleanOK := cleanTotal - len(cleanFails)
			t.Logf("clean mode: %d/%d inputs import and compile (%.1f%%)", cleanOK, cleanTotal, 100*float64(cleanOK)/float64(cleanTotal))
			for _, x := range cleanFails {
				t.Errorf("  CLEAN-FAIL %s", trunc(x))
			}
			if len(regressions) > 0 {
				t.Errorf("%d %s import golden(s) regressed (not in knownImportFailures)", len(regressions), f)
			}
			if rate < importThresholds[f] {
				t.Errorf("%s import match %.1f%% is below the %.0f%% gate", f, rate*100, importThresholds[f]*100)
			}
		})
	}
}

func firstLineDiff(want, got string) string {
	wl, gl := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(wl) || i < len(gl); i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return fmt.Sprintf("line %d: want %q got %q", i+1, w, g)
		}
	}
	return "equal"
}

func trunc(s string) string {
	if len(s) > 300 {
		return s[:300] + "..."
	}
	return s
}
