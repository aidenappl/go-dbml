package export

import "testing"

func TestCheckFragment(t *testing.T) {
	type res struct{ my, pg, ms bool } // true = accepted
	cases := []struct {
		frag string
		want res
	}{
		{"a > 0", res{true, true, true}},
		{"now()", res{true, true, true}},
		{"lower(\"Name\")", res{true, true, true}},
		{"a <> 'it''s'", res{true, true, true}},
		{"a ~ '^\\\\d'", res{true, true, true}}, // '\\d': same extent with and without escapes
		{"a <> 'x\\'", res{false, false, true}}, // trailing backslash: escapes the quote in MySQL / PG scs=off
		{"a <> 'x\\' OR a = ')); DROP TABLE v; -- '", res{false, false, true}},
		{"a <> E'x\\''", res{false, true, false}}, // E-string: PostgreSQL only
		{"`a` > 0", res{true, true, true}},        // backtick: plain punctuation outside MySQL
		{"[a] > 0", res{true, true, true}},        // brackets balanced as punctuation and as a T-SQL name
		{"a[1", res{true, true, false}},           // unterminated T-SQL bracket name
		{"$$x$$", res{true, true, true}},
		{"$$x", res{true, false, true}}, // unterminated dollar quote
		{"$1", res{true, true, true}},
		{"a > 0; DROP TABLE v", res{false, false, false}},
		{"a > 0 -- c", res{false, false, false}},
		{"a > 0 /* c */", res{false, false, false}},
		{"a # b", res{false, true, true}}, // MySQL comment
		{"(a > 0", res{false, false, false}},
		{"a > 0)", res{false, false, false}},
		{"a > 0\nGO\nDROP TABLE v", res{true, true, false}},
		{"a > 0\n  go 2 \nx", res{true, true, false}},
		{"a\nDELIMITER //\nx", res{false, true, true}},
		{"a \\! touch x", res{false, false, false}},
		{"'a\x00'", res{false, false, false}},
		{"\"a\\\" ; x\"", res{false, false, false}}, // MySQL: \" escapes the quote; elsewhere the ; is outside
		{"int(--1)", res{false, false, false}},
		{"varchar(255)", res{true, true, true}},
		{"enum('a','b')", res{true, true, true}},
		{"{fn now()}", res{true, true, false}},
	}
	for _, c := range cases {
		got := res{checkFragment(c.frag, MySQL) == "", checkFragment(c.frag, Postgres) == "", checkFragment(c.frag, MSSQL) == ""}
		if got != c.want {
			t.Errorf("%q: accepted (mysql, postgres, mssql) = %v, want %v (reasons %q %q %q)", c.frag, got, c.want,
				checkFragment(c.frag, MySQL), checkFragment(c.frag, Postgres), checkFragment(c.frag, MSSQL))
		}
	}
}

func TestMySQLQuoteReserved(t *testing.T) {
	for in, want := range map[string]string{
		`concat("key", '"key"', name)`: "concat(`key`, '\"key\"', name)",
		`"active" = 'x'`:               `"active" = 'x'`,
		"`order` > 0":                  "`order` > 0",
		`a = 'it''s "order"'`:          `a = 'it''s "order"'`,
	} {
		if got := mysqlQuoteReserved(in); got != want {
			t.Errorf("%s: %s, want %s", in, got, want)
		}
	}
}
