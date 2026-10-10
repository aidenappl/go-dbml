package export

import "testing"

func TestQuoting(t *testing.T) {
	for _, tc := range []struct {
		name string
		f    func(string) string
		in   string
		want string
	}{
		{"mysql ident", bq, "a`b", "`a``b`"},
		{"pg ident", dq, `a"b`, `"a""b"`},
		{"mssql ident", sb, "a]b[", "[a]]b[]"},
		{"mysql literal", myLit, `it's \ ok`, `'it''s \\ ok'`},
		{"std literal", stdLit, `it's \ ok`, `'it''s \ ok'`},
		{"number kept", func(s string) string { return safeNumber(s, stdLit) }, "-1.5e3", "-1.5e3"},
		{"non-number quoted", func(s string) string { return safeNumber(s, stdLit) }, "1); DROP TABLE x; --", "'1); DROP TABLE x; --'"},
	} {
		if got := tc.f(tc.in); got != tc.want {
			t.Errorf("%s(%q) = %s, want %s", tc.name, tc.in, got, tc.want)
		}
	}
}
