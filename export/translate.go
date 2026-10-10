package export

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// Options tunes SQLWith.
type Options struct {
	// TranslateTypes rewrites a model written for one dialect into another:
	// column types (int unsigned, tinyint(1), serial, datetime,
	// timestamptz, json/jsonb, uuid, blob/bytea, nvarchar(max), bit, ...),
	// auto-increment defaults (nextval → increment), common default and
	// CHECK functions (now()/getdate()/CURRENT_TIMESTAMP, uuid functions,
	// PostgreSQL ::casts, ~ / REGEXP) and index methods. It needs the source
	// dialect (SourceDialect, or the model's Project database_type) and does
	// nothing when the source equals the target. Fragments it cannot
	// translate safely are dropped with a warning instead of failing the
	// export. Off by default (the @dbml/core goldens are untranslated).
	TranslateTypes bool
	// SourceDialect overrides the dialect recorded in the model.
	SourceDialect Dialect
}

// Result is SQLWith's output.
type Result struct {
	SQL string
	// Warnings lists what TranslateTypes could not translate (kept as
	// written) or dropped.
	Warnings []string
}

// SourceDialectOf returns the dialect recorded in db's Project
// database_type ("" when absent or unknown).
func SourceDialectOf(db *model.Database) Dialect {
	if db == nil || db.Project == nil {
		return ""
	}
	v := db.Project.Extra["database_type"]
	switch strings.ToLower(strings.ReplaceAll(v, " ", "")) {
	case "mysql", "mariadb":
		return MySQL
	case "postgresql", "postgres", "pg":
		return Postgres
	case "sqlserver", "mssql", "microsoftsqlserver":
		return MSSQL
	}
	return ""
}

// DatabaseType is the Project database_type value naming d.
func DatabaseType(d Dialect) string {
	switch d {
	case MySQL:
		return "MySQL"
	case Postgres:
		return "PostgreSQL"
	case MSSQL:
		return "SQL Server"
	}
	return ""
}

// SQLWith is SQL with options.
func SQLWith(db *model.Database, d Dialect, o Options) (res *Result, err error) {
	if db == nil {
		return nil, fmt.Errorf("export: nil database")
	}
	defer func() {
		if r := recover(); r != nil {
			if ue, ok := r.(*UnsafeExpressionError); ok {
				res, err = nil, ue
				return
			}
			res, err = nil, &InternalError{Value: r}
		}
	}()
	c, err := buildCore(db)
	if err != nil {
		return nil, err
	}
	res = &Result{}
	src := o.SourceDialect
	if src == "" {
		src = SourceDialectOf(db)
	}
	if o.TranslateTypes && src != "" && src != d {
		tr := &translator{src: src, dst: d}
		tr.run(c)
		c.translated = true
		res.Warnings = tr.warns
	}
	var sql string
	switch d {
	case MySQL:
		sql, err = exportMySQL(c)
	case Postgres:
		sql, err = exportPostgres(c)
	case MSSQL:
		sql, err = exportMSSQL(c)
	default:
		return nil, fmt.Errorf("unsupported SQL dialect %q", d)
	}
	if err != nil {
		return nil, err
	}
	res.SQL = sql
	return res, nil
}

type translator struct {
	src, dst Dialect
	table    *cTable // the table being translated (expression identifiers)
	msPaths  map[string]bool
	warns    []string
	// keyed marks fields in a key (pk, unique, index, ref): text/blob types
	// get a bounded length where the target cannot index unbounded ones.
	keyed map[*cField]bool
}

func (tr *translator) warn(format string, a ...any) {
	tr.warns = append(tr.warns, fmt.Sprintf(format, a...))
}

// sqlType is a parsed column type.
type sqlType struct {
	base     string // lower case, single spaces, without args/array/unsigned
	args     []string
	unsigned bool
	array    bool
}

var (
	arrayDims  = regexp.MustCompile(`\s*\[[0-9]*\]$`)
	typeArgs   = regexp.MustCompile(`\(([^()]*)\)`)
	typeSpaces = regexp.MustCompile(`\s+`)
)

func parseType(t string) sqlType {
	s := strings.ToLower(strings.TrimSpace(t))
	var st sqlType
	for m := arrayDims.FindStringIndex(s); m != nil; m = arrayDims.FindStringIndex(s) {
		st.array = true
		s = strings.TrimSpace(s[:m[0]])
	}
	if m := typeArgs.FindStringSubmatch(s); m != nil {
		for _, a := range strings.Split(m[1], ",") {
			st.args = append(st.args, strings.TrimSpace(a))
		}
		s = typeArgs.ReplaceAllString(s, " ")
	}
	var words []string
	for _, w := range strings.Fields(s) {
		switch w {
		case "unsigned":
			st.unsigned = true
		case "signed", "zerofill":
		default:
			words = append(words, w)
		}
	}
	st.base = typeSpaces.ReplaceAllString(strings.Join(words, " "), " ")
	return st
}

func (st sqlType) arg(i int) (int, bool) {
	if i >= len(st.args) {
		return 0, false
	}
	n, err := strconv.Atoi(st.args[i])
	return n, err == nil
}

// kind is a dialect-neutral column type.
type kind struct {
	name     string // bool, int1..int8, numeric, float4, float8, money, char, varchar, text, binary, varbinary, blob, date, time, timetz, timestamp, timestamptz, interval, year, json, jsonb, uuid, xml, bit, varbit, inet, macaddr, set, array
	n, p, s  int    // length / precision / scale (0 = none)
	unsigned bool
	serial   bool // implies auto increment (PostgreSQL serial types)
}

// classify maps a source type to a kind (ok=false: unknown).
func (tr *translator) classify(st sqlType) (kind, bool) {
	n, hasN := st.arg(0)
	s, _ := st.arg(1)
	if st.array {
		return kind{name: "array"}, true
	}
	k := kind{unsigned: st.unsigned}
	if st.args != nil && strings.EqualFold(strings.Join(st.args, ""), "max") {
		n, hasN = -1, true // (MAX)
	}
	set := func(name string) (kind, bool) {
		k.name = name
		if hasN {
			k.n, k.p = n, n
		}
		k.s = s
		return k, true
	}
	switch st.base {
	case "bool", "boolean":
		return set("bool")
	case "bit":
		if tr.src == MSSQL || (hasN && n == 1) || !hasN {
			k.name = "bool"
			return k, true
		}
		return set("bit")
	case "bit varying", "varbit":
		return set("varbit")
	case "tinyint", "int1":
		if tr.src == MySQL && hasN && n == 1 {
			k.name = "bool"
			return k, true
		}
		if tr.src == MSSQL {
			k.unsigned = true // 0..255
		}
		k.name = "int1"
		return k, true
	case "smallint", "int2":
		k.name = "int2"
		return k, true
	case "smallserial", "serial2":
		k.name, k.serial = "int2", true
		return k, true
	case "mediumint", "int3", "middleint":
		k.name = "int3"
		return k, true
	case "int", "integer", "int4":
		k.name = "int4"
		return k, true
	case "serial", "serial4":
		k.name, k.serial = "int4", true
		return k, true
	case "bigint", "int8":
		k.name = "int8"
		return k, true
	case "bigserial", "serial8":
		k.name, k.serial = "int8", true
		return k, true
	case "decimal", "numeric", "dec", "fixed", "number":
		k.name = "numeric"
		if hasN {
			k.p, k.s = n, s
		}
		return k, true
	case "money", "smallmoney":
		k.name = "money"
		return k, true
	case "real", "float4":
		if tr.src == MySQL {
			k.name = "float8" // MySQL REAL = DOUBLE
		} else {
			k.name = "float4"
		}
		return k, true
	case "float":
		switch {
		case tr.src == MySQL && (!hasN || n <= 24):
			k.name = "float4"
		case tr.src == MSSQL && hasN && n <= 24:
			k.name = "float4"
		default:
			k.name = "float8"
		}
		return k, true
	case "double", "double precision", "float8":
		k.name = "float8"
		return k, true
	case "char", "character", "nchar", "national char", "national character", "bpchar":
		if !hasN {
			n, hasN = 1, true
		}
		return set("char")
	case "varchar", "character varying", "nvarchar", "national varchar", "national character varying", "char varying", "varchar2", "nvarchar2", "varcharacter":
		if hasN && n == -1 {
			k.name = "text"
			return k, true
		}
		return set("varchar")
	case "text", "tinytext", "mediumtext", "longtext", "ntext", "citext", "clob", "long", "long varchar", "string":
		k.name = "text"
		return k, true
	case "binary":
		return set("binary")
	case "varbinary", "binary varying":
		if hasN && n == -1 {
			k.name = "blob"
			return k, true
		}
		return set("varbinary")
	case "blob", "tinyblob", "mediumblob", "longblob", "bytea", "image", "long varbinary":
		k.name = "blob"
		return k, true
	case "date":
		k.name = "date"
		return k, true
	case "time", "time without time zone":
		return set("time")
	case "timetz", "time with time zone":
		return set("timetz")
	case "datetime", "datetime2", "smalldatetime", "timestamp without time zone":
		return set("timestamp")
	case "timestamp":
		if tr.src == MSSQL { // rowversion
			k.name, k.n = "binary", 8
			return k, true
		}
		return set("timestamp")
	case "timestamptz", "timestamp with time zone", "datetimeoffset":
		return set("timestamptz")
	case "year":
		k.name = "year"
		return k, true
	case "json":
		k.name = "json"
		return k, true
	case "jsonb":
		k.name = "jsonb"
		return k, true
	case "uuid", "uniqueidentifier":
		k.name = "uuid"
		return k, true
	case "xml":
		k.name = "xml"
		return k, true
	case "inet", "cidr":
		k.name = "inet"
		return k, true
	case "name":
		k.name, k.n = "varchar", 64
		return k, true
	case "macaddr", "macaddr8":
		k.name = "macaddr"
		return k, true
	case "set":
		k.name = "set"
		return k, true
	case "rowversion":
		k.name, k.n = "binary", 8
		return k, true
	case "sql_variant":
		k.name = "text"
		return k, true
	case "tsvector", "tsquery", "hstore", "ltree", "int4range", "int8range", "numrange", "tsrange", "tstzrange", "daterange":
		k.name = "text"
		return k, true
	case "point", "polygon", "linestring", "geometry", "multipoint", "multipolygon", "multilinestring", "geometrycollection":
		k.name = "geo:" + st.base
		return k, true
	case "circle", "box", "lseg", "line", "path":
		k.name = "text" // no equivalent outside PostgreSQL
		if tr.dst == Postgres {
			k.name = "geo:" + st.base
		}
		return k, true
	}
	if strings.HasPrefix(st.base, "interval") {
		k.name = "interval"
		return k, true
	}
	return k, false
}

// render writes kind k as a type of the target dialect.
func (tr *translator) render(k kind, keyed bool) string {
	d := tr.dst
	num := func(p, s int) string {
		switch {
		case p > 0 && s > 0:
			return "(" + strconv.Itoa(p) + "," + strconv.Itoa(s) + ")"
		case p > 0:
			return "(" + strconv.Itoa(p) + ")"
		}
		return ""
	}
	ln := func(n int) string {
		if n > 0 {
			return "(" + strconv.Itoa(n) + ")"
		}
		return ""
	}
	// unsigned integers widen outside MySQL
	if k.unsigned && d != MySQL {
		switch k.name {
		case "int1":
			if d == MSSQL {
				return "tinyint"
			}
			k.name = "int2"
		case "int2":
			k.name = "int4"
		case "int3":
			k.name = "int4"
		case "int4":
			k.name = "int8"
		case "int8":
			return "numeric(20)"
		}
	}
	us := ""
	if k.unsigned && d == MySQL {
		us = " unsigned"
	}
	switch k.name {
	case "bool":
		return map[Dialect]string{MySQL: "boolean", Postgres: "boolean", MSSQL: "bit"}[d]
	case "int1":
		if d == Postgres || (d == MSSQL && !k.unsigned) {
			return "smallint"
		}
		return "tinyint" + us
	case "int2":
		return "smallint" + us
	case "int3":
		if d == MySQL {
			return "mediumint" + us
		}
		return map[Dialect]string{Postgres: "integer", MSSQL: "int"}[d]
	case "int4":
		return map[Dialect]string{MySQL: "int", Postgres: "integer", MSSQL: "int"}[d] + us
	case "int8":
		return "bigint" + us
	case "numeric":
		p, s := k.p, k.s
		if p <= 0 {
			switch d {
			case MySQL:
				p, s = 65, 30
			case MSSQL:
				p, s = 38, 10
			default:
				return "numeric"
			}
		}
		if d == MySQL && p > 65 {
			p = 65
		}
		if d == MSSQL && p > 38 {
			p = 38
		}
		return map[Dialect]string{MySQL: "decimal", Postgres: "numeric", MSSQL: "decimal"}[d] + num(p, s)
	case "money":
		if d == MSSQL {
			return "money"
		}
		return map[Dialect]string{MySQL: "decimal(19,4)", Postgres: "numeric(19,4)"}[d]
	case "float4":
		return map[Dialect]string{MySQL: "float", Postgres: "real", MSSQL: "real"}[d]
	case "float8":
		return map[Dialect]string{MySQL: "double", Postgres: "double precision", MSSQL: "float"}[d]
	case "char":
		n := k.n
		switch {
		case d == MySQL && n > 255, d == MSSQL && n > 4000:
			return tr.render(kind{name: "varchar", n: n}, keyed)
		}
		return map[Dialect]string{MySQL: "char", Postgres: "char", MSSQL: "nchar"}[d] + ln(n)
	case "varchar":
		n := k.n
		switch d {
		case MySQL:
			if n <= 0 {
				n = 255
			}
			if n > 16383 {
				return tr.render(kind{name: "text"}, keyed)
			}
			return "varchar" + ln(n)
		case MSSQL:
			if n <= 0 {
				n = 255
			}
			if n > 4000 {
				return tr.render(kind{name: "text"}, keyed)
			}
			return "nvarchar" + ln(n)
		}
		return "varchar" + ln(n)
	case "text", "set", "xml", "interval":
		switch d {
		case MySQL:
			if keyed {
				return "varchar(255)"
			}
			if k.name == "interval" {
				return "varchar(64)"
			}
			return "longtext"
		case MSSQL:
			if k.name == "xml" {
				return "xml"
			}
			if keyed {
				return "nvarchar(450)"
			}
			if k.name == "interval" {
				return "nvarchar(64)"
			}
			return "nvarchar(max)"
		}
		if k.name == "xml" || k.name == "interval" {
			return k.name
		}
		return "text"
	case "binary":
		if d == Postgres {
			return "bytea"
		}
		return "binary" + ln(k.n)
	case "varbinary":
		if d == Postgres {
			return "bytea"
		}
		return "varbinary" + ln(k.n)
	case "blob":
		switch d {
		case MySQL:
			if keyed {
				return "varbinary(255)"
			}
			return "longblob"
		case MSSQL:
			if keyed {
				return "varbinary(900)"
			}
			return "varbinary(max)"
		}
		return "bytea"
	case "date":
		return "date"
	case "time":
		return "time" + ln(min(k.n, 6))
	case "timetz":
		if d == Postgres {
			return "time" + ln(min(k.n, 6)) + " with time zone"
		}
		return "time" + ln(min(k.n, 6))
	case "timestamp":
		return map[Dialect]string{MySQL: "datetime", Postgres: "timestamp", MSSQL: "datetime2"}[d] + ln(min(k.n, 6))
	case "timestamptz":
		return map[Dialect]string{MySQL: "datetime", Postgres: "timestamptz", MSSQL: "datetimeoffset"}[d] + ln(min(k.n, 6))
	case "year":
		if d == MySQL {
			return "year"
		}
		return "smallint"
	case "json", "jsonb":
		switch d {
		case MySQL:
			return "json"
		case MSSQL:
			return "nvarchar(max)"
		}
		return k.name
	case "uuid":
		return map[Dialect]string{MySQL: "char(36)", Postgres: "uuid", MSSQL: "uniqueidentifier"}[d]
	case "bit":
		if d == MSSQL {
			return "binary" + ln((k.n+7)/8)
		}
		return "bit" + ln(k.n)
	case "varbit":
		switch d {
		case MySQL:
			return "bit(64)"
		case MSSQL:
			return "varbinary(8)"
		}
		return "bit varying" + ln(k.n)
	case "inet":
		return map[Dialect]string{MySQL: "varchar(43)", Postgres: "inet", MSSQL: "nvarchar(43)"}[d]
	case "macaddr":
		return map[Dialect]string{MySQL: "varchar(17)", Postgres: "macaddr", MSSQL: "nvarchar(17)"}[d]
	case "array":
		return map[Dialect]string{MySQL: "json", MSSQL: "nvarchar(max)", Postgres: "text[]"}[d]
	}
	if g, ok := strings.CutPrefix(k.name, "geo:"); ok {
		switch {
		case d == MSSQL:
			return "geometry"
		case d == MySQL && (g == "point" || g == "polygon" || g == "linestring"):
			return g
		case d == MySQL && strings.HasPrefix(g, "multi") || d == MySQL && (g == "geometry" || g == "geometrycollection"):
			return g
		case d == Postgres && (g == "point" || g == "polygon" || g == "circle" || g == "box" || g == "lseg" || g == "line" || g == "path"):
			return g
		}
		return tr.render(kind{name: "text"}, keyed)
	}
	return ""
}

func (tr *translator) run(c *core) {
	tr.keyed = map[*cField]bool{}
	for _, s := range c.schemas {
		for _, t := range s.tables {
			for _, ix := range t.indexes {
				for _, ic := range ix.Columns {
					if f := t.findField(ic.Value); f != nil && ic.Type == "column" {
						tr.keyed[f] = true
					}
				}
			}
			for _, f := range t.fields {
				if f.pk || f.unique {
					tr.keyed[f] = true
				}
			}
		}
		for _, r := range s.refs {
			for _, ep := range r.endpoints {
				for _, f := range ep.fields {
					tr.keyed[f] = true
				}
			}
		}
	}
	types := map[string]string{} // source type -> target type (records)
	for _, s := range c.schemas {
		for _, t := range s.tables {
			where := t.name
			tr.table = t
			for _, f := range t.fields {
				tr.field(f, where, types)
			}
			if tr.dst == MySQL {
				// MySQL: a generated column may not refer to an AUTO_INCREMENT one
				for _, f := range t.fields {
					if f.generated == nil {
						continue
					}
					for _, inc := range t.fields {
						if inc.increment && regexp.MustCompile(`(?i)(^|[^A-Za-z0-9_$])`+"`?"+regexp.QuoteMeta(inc.name)+"`?"+`([^A-Za-z0-9_$]|$)`).MatchString(f.generated.Expression) {
							tr.warn("generated column %s.%s dropped (refers to the auto-increment column %s, which MySQL forbids)", t.name, f.name, inc.name)
							f.generated = nil
							break
						}
					}
				}
			}
			var checks []cCheck
			for _, ch := range t.checks {
				if e, ok := tr.expr(ch.expression, "check on table "+where, true); ok {
					checks = append(checks, cCheck{name: ch.name, expression: e})
				}
			}
			t.checks = checks
			t.indexes = tr.indexes(t)
			if tr.dst == Postgres {
				for i := range t.checks {
					t.checks[i].expression = boolCompare(t.checks[i].expression, t)
				}
				for _, f := range t.fields {
					for i := range f.checks {
						f.checks[i].expression = boolCompare(f.checks[i].expression, t)
					}
				}
			}
		}
		tr.indexNames(s)
	}
	for _, s := range c.schemas {
		var refs []*cRef
		for _, r := range s.refs {
			if tr.ref(r) {
				refs = append(refs, r)
			}
		}
		s.refs = refs
	}
	// record values typed by their column type
	if len(c.db.Records) > 0 {
		db := *c.db
		db.Records = nil
		for _, r := range c.db.Records {
			nr := *r
			nr.Values = make([][]*model.RecordValue, len(r.Values))
			for i, row := range r.Values {
				nr.Values[i] = make([]*model.RecordValue, len(row))
				for j, v := range row {
					if v == nil {
						continue
					}
					nv := *v
					if j < len(r.Columns) {
						if mt, ok := types[r.TableName+"\x00"+r.Columns[j]]; ok {
							nv.Type = mt
						}
					}
					nr.Values[i][j] = &nv
				}
			}
			db.Records = append(db.Records, &nr)
		}
		c.db = &db
	}
}

func (tr *translator) field(f *cField, table string, types map[string]string) {
	where := "column " + table + "." + f.name
	k := kind{}
	if f.enum == nil && (f.typ.schemaName == nil || *f.typ.schemaName == "" || *f.typ.schemaName == defaultSchema) && f.typ.originalTypeName == nil {
		st := parseType(f.typ.typeName)
		var ok bool
		k, ok = tr.classify(st)
		if ok && k.name == "bool" && tr.src == MySQL {
			if safe, why := tr.boolSafe(f); !safe {
				k = kind{name: "int1"}
				tr.warn("%s: %s kept numeric (%s), not translated to a boolean", where, f.typ.typeName, why)
			}
		}
		if ok {
			nt := tr.render(k, tr.keyed[f])
			if note := tr.lossyNote(k, nt, tr.keyed[f]); note != "" {
				tr.warn("%s: %s: %s", where, f.typ.typeName, note)
			}
			types[table+"\x00"+f.name] = nt
			if (f.increment || k.serial) && strings.HasPrefix(nt, "numeric") {
				nt = "bigint" // an identity column must be an integer type
			}
			f.typ = &cType{typeName: nt}
			if k.serial {
				f.increment = true
			}
		} else {
			tr.warn("%s: type %q kept as written (no %s equivalent known)", where, f.typ.typeName, tr.dst)
		}
	} else if f.enum != nil && tr.dst == MSSQL {
		tr.warn("%s: enum %s becomes nvarchar(255) with a CHECK constraint", where, f.enum.name)
	}
	if f.increment && tr.dst == MySQL && !f.pk {
		// MySQL: the auto column must be a key
		f.increment = false
		tr.warn("%s: auto increment dropped (not the primary key)", where)
	}
	if d := f.dbdefault; d != nil {
		f.dbdefault = tr.defaultValue(f, k, d, where)
	}
	var checks []cCheck
	for _, ch := range f.checks {
		if e, ok := tr.expr(ch.expression, "check on "+where, true); ok {
			checks = append(checks, cCheck{name: ch.name, expression: e})
		}
	}
	f.checks = checks
	if g := f.generated; g != nil {
		e, ok := tr.expr(g.Expression, "generated "+where, false)
		if ok && tr.dst == Postgres && !immutableExpr(e) {
			tr.warn("%s: generated expression %q dropped (PostgreSQL requires an immutable expression)", where, g.Expression)
			ok = false
		}
		if ok {
			f.generated = &model.GeneratedColumn{Expression: e, Stored: g.Stored}
		} else {
			f.generated = nil
		}
	}
}

var (
	pgCast        = regexp.MustCompile(`::\s*(?:(?:"[^"]*"|[A-Za-z_][A-Za-z0-9_]*)\.?)+(?:\s+(?:varying|precision|without\s+time\s+zone|with\s+time\s+zone))?(?:\s*\([0-9, ]*\))?(?:\[[0-9]*\])*`)
	nextvalRe     = regexp.MustCompile(`(?i)^\s*nextval\s*\(`)
	nowFuncs      = regexp.MustCompile(`(?i)\b(now|getdate|sysdatetime|sysutcdatetime|getutcdate|current_timestamp|localtimestamp|localtime|transaction_timestamp|statement_timestamp|clock_timestamp|utc_timestamp)\s*\(\s*[0-9]*\s*\)|\b(current_timestamp|localtimestamp)\b`)
	dateFuncs     = regexp.MustCompile(`(?i)\b(curdate|current_date)\s*\(\s*\)|\bcurrent_date\b`)
	uuidFuncs     = regexp.MustCompile(`(?i)\b(gen_random_uuid|uuid_generate_v4|uuid_generate_v1|newid|newsequentialid|uuid)\s*\(\s*\)`)
	pgRegexOp     = regexp.MustCompile(`\s!?~\*?\s`)
	myRegexOp     = regexp.MustCompile(`(?i)\s(not\s+)?(regexp|rlike)\s`)
	lenFunc       = regexp.MustCompile(`(?i)\b(len|char_length|character_length|length)\s*\(`)
	isnullFunc    = regexp.MustCompile(`(?i)\b(isnull|ifnull|nvl)\s*\(`)
	zeroDate      = regexp.MustCompile(`^0000-00-00`)
	funcCall      = regexp.MustCompile(`([A-Za-z_][A-Za-z0-9_.]*)\(`)
	pgIntervalLit = regexp.MustCompile(`(?i)\binterval\s+'\s*(\d+)\s+(year|month|week|day|hour|minute|second)s?\s*'`)
	myInterval    = regexp.MustCompile(`(?i)\binterval\s+(\d+)\s+(year|month|week|day|hour|minute|second)\b`)
	bitLiteral    = regexp.MustCompile(`(?i)^\(*\s*b'([01]+)'\s*\)*$|^\(*\s*0x0*([01])\s*\)*$`)
	strPlus       = regexp.MustCompile(`(N?'(?:[^']|'')*')\s*\+|\+\s*(N?'(?:[^']|'')*')`)
	pgOnlyOps     = regexp.MustCompile(`->|#>|@>|<@|\?\||\?&|~~|!~|\|\||\[`)
	identTok      = regexp.MustCompile(`'(?:[^']|'')*'|"(?:[^"]|"")*"|[A-Za-z_][A-Za-z0-9_]*`)
)

// portableFuncs are the functions translated expressions may call: they
// exist (possibly after the rewrites in expr) in MySQL/MariaDB,
// PostgreSQL and SQL Server.
var portableFuncs = setOf("LOWER", "UPPER", "COALESCE", "NULLIF", "ABS", "ROUND", "FLOOR", "CEILING", "TRIM", "LTRIM",
	"RTRIM", "SUBSTRING", "REPLACE", "POWER", "SQRT", "SIGN", "CAST",
	// rewritten per target by expr
	"CHAR_LENGTH", "LENGTH", "LEN", "NOW", "GETDATE", "CURDATE", "UUID", "GEN_RANDOM_UUID", "NEWID",
	// SQL keywords followed directly by "("
	"IN", "AND", "OR", "NOT", "EXISTS", "ANY", "ALL", "VALUES", "CASE", "WHEN", "THEN", "ELSE", "BETWEEN", "IS", "LIKE", "AS")

// immutableFuncs may appear in PostgreSQL generated columns and index
// expressions.
var immutableFuncs = setOf("LOWER", "UPPER", "ABS", "COALESCE", "LENGTH", "ROUND", "FLOOR", "CEIL", "CEILING", "TRIM",
	"MD5", "NULLIF", "GREATEST", "LEAST", "SUBSTR", "SUBSTRING", "IN", "AND", "OR", "NOT", "CASE", "WHEN", "THEN", "ELSE", "IS", "LIKE")

func callsOnly(e string, ok map[string]bool) (string, bool) {
	for _, m := range funcCall.FindAllStringSubmatch(stripStrings(e), -1) {
		name := strings.ToUpper(m[1])
		if i := strings.LastIndexByte(name, '.'); i >= 0 {
			name = name[i+1:]
		}
		if !ok[name] {
			return m[1], false
		}
	}
	return "", true
}

// stripStrings blanks string literals (function-like text inside them is
// not a call).
func stripStrings(e string) string {
	return regexp.MustCompile(`'(?:[^']|'')*'`).ReplaceAllString(e, "''")
}

// immutableExpr reports whether e calls only functions PostgreSQL treats
// as immutable (a conservative list).
func immutableExpr(e string) bool {
	_, ok := callsOnly(e, immutableFuncs)
	return ok
}

func (tr *translator) nowExpr() string {
	return map[Dialect]string{MySQL: "CURRENT_TIMESTAMP", Postgres: "now()", MSSQL: "getdate()"}[tr.dst]
}

// expr translates a raw SQL expression. ok=false: dropped (warned).
func (tr *translator) expr(e, where string, isCheck bool) (string, bool) {
	out := e
	if tr.dst != Postgres {
		if c := pgCast.ReplaceAllString(out, ""); c != out {
			tr.warn("%s: PostgreSQL casts removed from %q", where, e)
			out = c
		}
	}
	if isCheck && tr.dst == MySQL && (nowFuncs.MatchString(out) || dateFuncs.MatchString(out)) {
		tr.warn("%s: %q dropped (MySQL CHECK constraints must be deterministic)", where, e)
		return "", false
	}
	if tr.dst != Postgres && pgOnlyOps.MatchString(stripStrings(out)) {
		tr.warn("%s: %q dropped (PostgreSQL-only operator)", where, e)
		return "", false
	}
	out = nowFuncs.ReplaceAllString(out, tr.nowExpr())
	out = dateFuncs.ReplaceAllString(out, map[Dialect]string{MySQL: "curdate()", Postgres: "CURRENT_DATE", MSSQL: "CAST(getdate() AS date)"}[tr.dst])
	out = uuidFuncs.ReplaceAllString(out, map[Dialect]string{MySQL: "uuid()", Postgres: "gen_random_uuid()", MSSQL: "newid()"}[tr.dst])
	out = isnullFunc.ReplaceAllString(out, "COALESCE(")
	switch tr.dst {
	case MySQL:
		out = pgIntervalLit.ReplaceAllString(out, "INTERVAL $1 $2")
	case Postgres:
		out = myInterval.ReplaceAllString(out, "interval '$1 $2'")
		// T-SQL string concatenation: 'x' + col
		for prev := ""; prev != out; {
			prev = out
			out = strPlus.ReplaceAllStringFunc(out, func(m string) string { return strings.Replace(m, "+", "||", 1) })
		}
	}
	switch {
	case tr.dst == MySQL && tr.src == Postgres:
		out = pgRegexOp.ReplaceAllStringFunc(out, func(op string) string {
			if strings.Contains(op, "!") {
				return " NOT REGEXP "
			}
			return " REGEXP "
		})
	case tr.dst == Postgres && tr.src == MySQL:
		out = myRegexOp.ReplaceAllStringFunc(out, func(op string) string {
			if strings.Contains(strings.ToLower(op), "not") {
				return " !~ "
			}
			return " ~ "
		})
	}
	out = lenFunc.ReplaceAllStringFunc(out, func(string) string {
		return map[Dialect]string{MySQL: "CHAR_LENGTH(", Postgres: "length(", MSSQL: "LEN("}[tr.dst]
	})
	if tr.dst == MSSQL && myInterval.MatchString(out) || tr.dst == MSSQL && pgIntervalLit.MatchString(out) {
		tr.warn("%s: %q dropped (interval arithmetic)", where, e)
		return "", false
	}
	if f, ok := callsOnly(out, portableFuncs); !ok {
		tr.warn("%s: %q dropped (function %s is not portable to %s)", where, e, f, tr.dst)
		return "", false
	}
	out = tr.idents(out)
	if r := checkFragment(out, tr.dst); r != "" {
		tr.warn("%s: %q dropped (%s for %s)", where, e, r, tr.dst)
		return "", false
	}
	return out, true
}

// idents adapts identifier quoting: "x" is an identifier in the source
// (PostgreSQL / SQL Server, and normalised imports) and becomes `x` for
// MySQL; a bare column name with upper-case letters (case-insensitive in
// MySQL / SQL Server) becomes "Name" for PostgreSQL, which folds bare names
// to lower case while the exporter creates the column quoted.
func (tr *translator) idents(e string) string {
	var b strings.Builder
	for i := 0; i < len(e); i++ {
		c := e[i]
		switch {
		case c == '\'':
			j := scanQuoted(e, i, '\'', tr.dst == MySQL)
			if j < 0 {
				b.WriteString(e[i:])
				return b.String()
			}
			b.WriteString(e[i : j+1])
			i = j
		case c == '"':
			j := scanQuoted(e, i, '"', false)
			if j < 0 {
				b.WriteString(e[i:])
				return b.String()
			}
			if tr.dst == MySQL && tr.src != MySQL {
				b.WriteString(bq(strings.ReplaceAll(e[i+1:j], `""`, `"`)))
			} else {
				b.WriteString(e[i : j+1])
			}
			i = j
		case c == '`':
			j := scanQuoted(e, i, '`', false)
			if j < 0 {
				b.WriteString(e[i:])
				return b.String()
			}
			b.WriteString(e[i : j+1])
			i = j
		case c == '_' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z':
			j := i
			for j < len(e) && (e[j] == '_' || e[j] == '$' || e[j] >= '0' && e[j] <= '9' || e[j] >= 'a' && e[j] <= 'z' || e[j] >= 'A' && e[j] <= 'Z') {
				j++
			}
			w := e[i:j]
			if tr.dst == Postgres && tr.src != Postgres && tr.table != nil && (i == 0 || e[i-1] != '.') && !strings.HasPrefix(strings.TrimLeft(e[j:], " "), "(") {
				for _, f := range tr.table.fields {
					if strings.EqualFold(f.name, w) && f.name != strings.ToLower(f.name) {
						w = dq(f.name)
						break
					}
				}
			}
			b.WriteString(w)
			i = j - 1
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

// refersToColumn reports whether expression e names a column of t.
func refersToColumn(e string, t *cTable) bool {
	for _, m := range identTok.FindAllString(e, -1) {
		if m[0] == '\'' {
			continue
		}
		name := strings.Trim(m, `"`)
		if f := t.findField(name); f != nil {
			return true
		}
	}
	return false
}

func (tr *translator) defaultValue(f *cField, k kind, d *model.Default, where string) *model.Default {
	v := jsString(d.Value)
	switch d.Type {
	case "expression":
		if nextvalRe.MatchString(v) {
			if f.pk || tr.dst != MySQL {
				f.increment = true
			}
			return nil
		}
		if tr.dst != Postgres && regexp.MustCompile(`(?i)^\s*'[^']*'\s*::`).MatchString(v) {
			// 'x'::type → a plain string default
			lit := strings.TrimSpace(pgCast.ReplaceAllString(v, ""))
			if len(lit) >= 2 {
				return tr.defaultValue(f, k, &model.Default{Type: "string", Value: strings.ReplaceAll(lit[1:len(lit)-1], "''", "'")}, where)
			}
		}
		if bm := bitLiteral.FindStringSubmatch(v); bm != nil && (k.name == "bool" || k.name == "bit") {
			b := bm[1] + bm[2]
			return tr.boolDefault(map[bool]string{true: "1", false: "0"}[strings.Contains(b, "1")], where)
		}
		if tr.dst != MySQL && (refersToColumn(v, f.table) || hasBareIdent(v)) {
			tr.warn("default of %s: %q dropped (refers to a column)", where, v)
			return nil
		}
		if (k.name == "bool") && (v == "1" || v == "0" || strings.EqualFold(v, "true") || strings.EqualFold(v, "false") || v == "(1)" || v == "(0)" || v == "((1))" || v == "((0))") {
			return tr.boolDefault(strings.Trim(v, "()"), where)
		}
		e, ok := tr.expr(v, "default of "+where, false)
		if !ok {
			return nil
		}
		return &model.Default{Type: "expression", Value: e}
	case "string":
		if zeroDate.MatchString(v) && tr.dst != MySQL {
			tr.warn("%s: zero-date default %q dropped", where, v)
			return nil
		}
		switch k.name {
		case "int1", "int2", "int3", "int4", "int8", "numeric", "float4", "float8", "money":
			if !numberText.MatchString(strings.TrimSpace(v)) {
				tr.warn("%s: non-numeric default %q dropped", where, v)
				return nil
			}
		case "timestamptz":
			if tr.dst != Postgres {
				return &model.Default{Type: "string", Value: tzSuffix.ReplaceAllString(v, "")}
			}
		case "blob", "binary", "varbinary":
			if tr.dst == Postgres {
				tr.warn("%s: binary default %q dropped", where, v)
				return nil
			}
		}
		if k.name == "bool" {
			return tr.boolDefault(v, where)
		}
	case "number":
		if k.name == "bool" {
			return tr.boolDefault(v, where)
		}
		if (k.name == "blob" || k.name == "binary" || k.name == "varbinary") && tr.dst != MySQL {
			tr.warn("%s: numeric default %q dropped for a binary column", where, v)
			return nil
		}
	case "boolean":
		switch k.name {
		case "int1", "int2", "int3", "int4", "int8", "numeric":
			if b, ok := clearBool(v); ok {
				return &model.Default{Type: "number", Value: map[bool]float64{true: 1, false: 0}[b]}
			}
		default:
			if tr.dst == MSSQL {
				return tr.boolDefault(v, where)
			}
		}
	}
	return d
}

func (tr *translator) boolDefault(v, where string) *model.Default {
	b, ok := clearBool(v)
	if !ok {
		tr.warn("%s: default %q dropped (not a boolean value)", where, v)
		return nil
	}
	t := b
	if tr.dst == MSSQL {
		return &model.Default{Type: "number", Value: map[bool]float64{true: 1, false: 0}[t]}
	}
	return &model.Default{Type: "boolean", Value: map[bool]string{true: "true", false: "false"}[t]}
}

// indexes translates index methods and expression parts.
func (tr *translator) indexes(t *cTable) []*model.Index {
	var out []*model.Index
	for _, ix := range t.indexes {
		nx := *ix
		cols := make([]*model.IndexColumn, 0, len(ix.Columns))
		drop := false
		for _, c := range ix.Columns {
			nc := *c
			if c.Type == "expression" {
				if m := prefixPart.FindStringSubmatch(c.Value); m != nil && tr.src == MySQL && t.findField(m[1]) != nil {
					// MySQL prefix key part: the whole column elsewhere
					nc = model.IndexColumn{Type: "column", Value: m[1]}
				} else if tr.dst == MySQL || tr.dst == MSSQL {
					drop = true // no functional index parts (MariaDB, SQL Server)
				} else if e, ok := tr.expr(c.Value, "index expression on "+t.name, false); ok && (tr.dst != Postgres || immutableExpr(e)) {
					nc.Value = e
				} else {
					drop = true
				}
			}
			cols = append(cols, &nc)
		}
		if drop {
			tr.warn("index on %s with expression parts dropped (not portable to %s)", t.name, tr.dst)
			continue
		}
		nx.Columns = cols
		if c := tr.unindexable(t, cols); c != "" {
			tr.warn("index on %s dropped: column %s cannot be indexed in %s", t.name, c, tr.dst)
			continue
		}
		if ix.Type != nil {
			typ := strings.ToLower(*ix.Type)
			keep := false
			switch tr.dst {
			case MySQL:
				keep = typ == "btree" || typ == "hash"
			case Postgres:
				keep = typ == "btree" || (typ == "hash" && len(cols) == 1 && !truthy(ix.Unique)) || typ == "gin" || typ == "gist" || typ == "brin" || typ == "spgist"
				keep = keep && (tr.src == Postgres || typ == "btree" || typ == "hash")
			}
			if !keep {
				tr.warn("index on %s: method %s dropped (not available in %s)", t.name, *ix.Type, tr.dst)
				nx.Type = nil
			}
		}
		out = append(out, &nx)
	}
	return out
}

var boolCmp = regexp.MustCompile(`(?i)("?)([A-Za-z_][A-Za-z0-9_]*)("?)\s*(=|<>|!=)\s*([01])\b`)

// boolCompare rewrites `bool_column = 1` (MySQL / SQL Server style) into
// `bool_column = true` for PostgreSQL.
func boolCompare(e string, t *cTable) string {
	return boolCmp.ReplaceAllStringFunc(e, func(m string) string {
		sm := boolCmp.FindStringSubmatch(m)
		f := t.findField(sm[2])
		if f == nil || f.typ == nil || !isBooleanType(f.typ.typeName) {
			return m
		}
		return sm[1] + sm[2] + sm[3] + " " + sm[4] + " " + map[string]string{"1": "true", "0": "false"}[sm[5]]
	})
}

// isKey reports whether fields are the primary key or covered by a unique
// constraint of their table.
func isKey(fs []*cField) bool {
	if len(fs) == 0 {
		return false
	}
	t := fs[0].table
	set := map[string]bool{}
	for _, f := range fs {
		set[f.name] = true
	}
	same := func(cols []*model.IndexColumn) bool {
		if len(cols) != len(set) {
			return false
		}
		for _, c := range cols {
			if c.Type != "column" || !set[c.Value] {
				return false
			}
		}
		return true
	}
	if len(fs) == 1 && (fs[0].pk || fs[0].unique) {
		return true
	}
	var pks []string
	for _, f := range t.fields {
		if f.pk {
			pks = append(pks, f.name)
		}
	}
	if len(pks) == len(set) && len(pks) > 0 {
		all := true
		for _, n := range pks {
			all = all && set[n]
		}
		if all {
			return true
		}
	}
	for _, ix := range t.indexes {
		if (truthy(ix.PK) || truthy(ix.Unique)) && same(ix.Columns) {
			return true
		}
	}
	return false
}

// indexed reports whether fs lead some index of their table (MySQL needs an
// index on referenced columns).
func indexed(fs []*cField) bool {
	if isKey(fs) {
		return true
	}
	t := fs[0].table
	for _, ix := range t.indexes {
		if len(ix.Columns) < len(fs) {
			continue
		}
		ok := true
		for i, f := range fs {
			ok = ok && ix.Columns[i].Type == "column" && ix.Columns[i].Value == f.name
		}
		if ok {
			return true
		}
	}
	return false
}

// ref adapts a foreign key to the target: the referenced columns must be
// a key (PostgreSQL / SQL Server) or indexed (MySQL), and the referencing
// columns get the referenced columns' types. false = drop the ref.
func (tr *translator) ref(r *cRef) bool {
	refEp, forEp, m2m := refSides(r)
	sides := []*cEndpoint{refEp}
	if m2m {
		sides = append(sides, forEp)
	}
	for _, ep := range sides {
		if len(ep.fields) == 0 {
			continue
		}
		t := ep.fields[0].table
		switch {
		case tr.dst == MySQL && !indexed(ep.fields):
			var cols []*model.IndexColumn
			for _, f := range ep.fields {
				cols = append(cols, &model.IndexColumn{Type: "column", Value: f.name})
			}
			t.indexes = append(t.indexes, &model.Index{Columns: cols})
		case tr.dst != MySQL && !isKey(ep.fields):
			tr.warn("foreign key to %s(%s) dropped: the referenced columns are not unique in %s", t.name, fieldList(ep.fields), tr.dst)
			return false
		}
	}
	if tr.dst == MSSQL {
		for _, a := range []**string{&r.onDelete, &r.onUpdate} {
			if *a != nil && strings.EqualFold(**a, "restrict") {
				na := "no action"
				*a = &na
			}
		}
		// SQL Server refuses cascading actions on a self reference or on a
		// second path between the same tables
		parent, child := refEp.fields[0].table, forEp.fields[0].table
		key := parent.name + "\x00" + child.name
		cascading := func(a *string) bool { return a != nil && !strings.EqualFold(*a, "no action") }
		if cascading(r.onDelete) || cascading(r.onUpdate) {
			if parent == child || tr.msPaths[key] {
				tr.warn("foreign key %s -> %s: cascading action dropped (SQL Server refuses cycles / multiple cascade paths)", child.name, parent.name)
				r.onDelete, r.onUpdate = nil, nil
			}
			if tr.msPaths == nil {
				tr.msPaths = map[string]bool{}
			}
			tr.msPaths[key] = true
		}
	}
	if (tr.dst == MySQL || tr.dst == MSSQL) && !m2m {
		for _, f := range forEp.fields {
			if f.notNull {
				for _, a := range []**string{&r.onDelete, &r.onUpdate} {
					if *a != nil && strings.EqualFold(**a, "set null") {
						tr.warn("foreign key %s.%s: ON … SET NULL on a NOT NULL column dropped (%s rejects it)", f.table.name, f.name, tr.dst)
						*a = nil
					}
				}
			}
		}
	}
	if !m2m && len(refEp.fields) == len(forEp.fields) {
		for i, f := range forEp.fields {
			p := refEp.fields[i]
			if f.enum != nil || p.enum != nil || f.typ == nil || p.typ == nil || f.typ.typeName == p.typ.typeName {
				continue
			}
			ck, _ := tr.classify(parseType(f.typ.typeName))
			pk, _ := tr.classify(parseType(p.typ.typeName))
			if tr.dst == MySQL || ck.name != pk.name || ck.n != pk.n || ck.p != pk.p || ck.unsigned != pk.unsigned {
				tr.warn("column %s.%s: type %s changed to %s to match the referenced %s.%s", f.table.name, f.name, f.typ.typeName, p.typ.typeName, p.table.name, p.name)
				f.typ = &cType{typeName: p.typ.typeName}
			}
		}
	}
	return true
}

func fieldList(fs []*cField) string {
	var n []string
	for _, f := range fs {
		n = append(n, f.name)
	}
	return strings.Join(n, ", ")
}

// indexNames makes index names unique where the target requires it:
// PostgreSQL per schema (tables included), MySQL and SQL Server per table.
func (tr *translator) indexNames(s *cSchema) {
	used := map[string]bool{}
	if tr.dst == Postgres {
		for _, t := range s.tables {
			used[strings.ToLower(t.name)] = true
		}
	}
	for _, t := range s.tables {
		if tr.dst != Postgres {
			used = map[string]bool{}
		}
		for _, ix := range t.indexes {
			if ix.Name == nil || *ix.Name == "" {
				continue
			}
			name := *ix.Name
			if used[strings.ToLower(name)] {
				base := t.name + "_" + name
				name = base
				for n := 2; used[strings.ToLower(name)]; n++ {
					name = base + "_" + strconv.Itoa(n)
				}
				tr.warn("index %s on %s renamed to %s (names are unique per %s in %s)", *ix.Name, t.name, name,
					map[bool]string{true: "schema", false: "table"}[tr.dst == Postgres], tr.dst)
				ix.Name = &name
			}
			used[strings.ToLower(name)] = true
		}
	}
}

var tzSuffix = regexp.MustCompile(`\s*([+-]\d\d(:?\d\d)?|Z)$`)

// defaultKeywords may appear bare in a default expression.
var defaultKeywords = setOf("TRUE", "FALSE", "NULL", "AND", "OR", "NOT", "IS", "IN", "LIKE", "BETWEEN", "CASE", "WHEN",
	"THEN", "ELSE", "END", "AS", "INTERVAL", "CURRENT_DATE", "CURRENT_TIMESTAMP", "CURRENT_TIME", "LOCALTIMESTAMP",
	"LOCALTIME", "YEAR", "MONTH", "WEEK", "DAY", "HOUR", "MINUTE", "SECOND", "DATE", "CHAR", "VARCHAR", "INT", "INTEGER",
	"TEXT", "SIGNED", "UNSIGNED", "DECIMAL", "NUMERIC")

// hasBareIdent reports whether expression e has an identifier that is not
// a function name, a keyword, a string prefix or a CAST target type: a
// column reference, which PostgreSQL and SQL Server refuse in a DEFAULT.
func hasBareIdent(e string) bool {
	locs := identTok.FindAllStringIndex(e, -1)
	for i, l := range locs {
		tok := e[l[0]:l[1]]
		if tok[0] == '\'' {
			continue
		}
		rest := strings.TrimLeft(e[l[1]:], " \t")
		if strings.HasPrefix(rest, "(") || strings.HasPrefix(e[l[1]:], "'") || defaultKeywords[strings.ToUpper(strings.Trim(tok, `"`))] {
			continue
		}
		if i > 0 {
			p := e[locs[i-1][0]:locs[i-1][1]]
			if strings.EqualFold(p, "AS") {
				continue
			}
		}
		return true
	}
	return false
}

// unindexable returns a column of cols whose target type the target cannot
// index with a plain index ("" = none).
func (tr *translator) unindexable(t *cTable, cols []*model.IndexColumn) string {
	for _, c := range cols {
		if c.Type != "column" {
			continue
		}
		f := t.findField(c.Value)
		if f == nil || f.typ == nil {
			continue
		}
		ty := strings.ToLower(f.typ.typeName)
		switch tr.dst {
		case MSSQL:
			if strings.Contains(ty, "(max)") || ty == "xml" || ty == "geometry" || ty == "text" || ty == "ntext" || ty == "image" {
				return f.name
			}
		case MySQL:
			if ty == "json" || ty == "longtext" || ty == "longblob" || ty == "geometry" || ty == "point" || ty == "polygon" || ty == "linestring" || strings.HasPrefix(ty, "multi") {
				return f.name
			}
		case Postgres:
			switch ty {
			case "point", "polygon", "circle", "box", "lseg", "line", "path", "json", "xml":
				return f.name
			}
		}
	}
	return ""
}

// clearBool parses a value that unambiguously means a boolean.
func clearBool(v string) (bool, bool) {
	switch strings.ToLower(strings.Trim(strings.TrimSpace(v), "'")) {
	case "1", "true", "t", "y", "yes", "b'1'":
		return true, true
	case "0", "false", "f", "n", "no", "b'0'":
		return false, true
	}
	return false, false
}

// boolSafe reports whether a MySQL tinyint(1) / bit(1) column can become a
// boolean: its default (if any) is 0/1/true/false and no CHECK, index
// expression or generated column of its table uses it (they would compare
// it with numbers).
func (tr *translator) boolSafe(f *cField) (bool, string) {
	if d := f.dbdefault; d != nil && !isNullDefault(d) {
		v := jsString(d.Value)
		if bm := bitLiteral.FindStringSubmatch(v); bm != nil {
			v = bm[1] + bm[2]
		}
		if _, ok := clearBool(strings.Trim(v, "()")); !ok {
			return false, "default " + v
		}
	}
	col := `["` + "`" + `]?` + regexp.QuoteMeta(f.name) + `["` + "`" + `]?`
	re := regexp.MustCompile(`(?i)(^|[^A-Za-z0-9_$])` + col + `([^A-Za-z0-9_$]|$)`)
	// boolean uses: col = TRUE / 1 / FALSE / 0 (either side), col IS [NOT]
	// TRUE|FALSE|NULL, NOT col, col alone between AND / OR / parentheses
	lit := `(?:TRUE|FALSE|0|1)`
	boolUses := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(^|[^A-Za-z0-9_$])` + col + `\s*(?:=|<>|!=)\s*` + lit + `\b`),
		regexp.MustCompile(`(?i)\b` + lit + `\s*(?:=|<>|!=)\s*` + col + `([^A-Za-z0-9_$]|$)`),
		regexp.MustCompile(`(?i)(^|[^A-Za-z0-9_$])` + col + `\s+IS\s+(?:NOT\s+)?(?:TRUE|FALSE|NULL|UNKNOWN)\b`),
		regexp.MustCompile(`(?i)\bNOT\s+` + col + `([^A-Za-z0-9_$]|$)`),
		regexp.MustCompile(`(?i)(^|\(|\bAND\b|\bOR\b)\s*` + col + `\s*(\)|\bAND\b|\bOR\b|$)`),
	}
	t := f.table
	// a use is numeric when the column is still referenced once its
	// boolean uses are removed
	uses := func(e string) bool {
		e = stripStrings(e)
		for _, b := range boolUses {
			for prev := ""; prev != e; {
				prev = e
				e = b.ReplaceAllString(e, " ")
			}
		}
		return re.MatchString(e)
	}
	for _, ch := range t.checks {
		if uses(ch.expression) {
			return false, "used numerically in a CHECK"
		}
	}
	for _, g := range t.fields {
		for _, ch := range g.checks {
			if uses(ch.expression) {
				return false, "used numerically in a CHECK"
			}
		}
		if g.generated != nil && uses(g.generated.Expression) {
			return false, "used in a generated column"
		}
	}
	for _, ix := range t.indexes {
		for _, c := range ix.Columns {
			if c.Type == "expression" && uses(c.Value) {
				return false, "used in an index expression"
			}
		}
	}
	return true, ""
}

// lossyNote describes what rendering kind k for the target loses or
// changes ("" = nothing).
func (tr *translator) lossyNote(k kind, nt string, keyed bool) string {
	d := tr.dst
	switch {
	case k.unsigned && d != MySQL:
		return "unsigned dropped, widened to " + nt
	case k.name == "set":
		return "SET values become free text (" + nt + ")"
	case k.name == "year" && d != MySQL:
		return "YEAR becomes " + nt
	case (k.name == "json" || k.name == "jsonb") && d == MSSQL:
		return "JSON stored as " + nt
	case k.name == "jsonb" && d == MySQL:
		return "jsonb becomes json (no binary representation)"
	case k.name == "xml" && d == MySQL, k.name == "interval" && d != Postgres, k.name == "inet" && d != Postgres,
		k.name == "macaddr" && d != Postgres, k.name == "array" && d != Postgres:
		return k.name + " stored as text " + nt
	case k.name == "uuid" && d == MySQL:
		return "uuid stored as " + nt
	case k.name == "money" && d != MSSQL:
		return "money becomes " + nt
	case (k.name == "timestamptz" || k.name == "timetz") && d == MySQL:
		return "time zone dropped (" + nt + ")"
	case k.name == "timetz" && d == MSSQL:
		return "time zone dropped (" + nt + ")"
	case (k.name == "time" || k.name == "timestamp" || k.name == "timestamptz" || k.name == "timetz") && k.n > 6:
		return "fractional seconds precision reduced to 6"
	case k.name == "numeric" && (k.p <= 0 && d != Postgres || d == MySQL && k.p > 65 || d == MSSQL && k.p > 38):
		return "precision becomes " + nt
	case (k.name == "text" || k.name == "blob") && keyed && (d == MySQL || d == MSSQL):
		return "length bounded to " + nt + " (the column is part of a key)"
	case k.name == "char" && (d == MySQL && k.n > 255 || d == MSSQL && k.n > 4000),
		k.name == "varchar" && (d == MySQL && k.n > 16383 || d == MSSQL && k.n > 4000):
		return "fixed/limited length becomes " + nt
	case k.name == "bit" && d == MSSQL, k.name == "varbit":
		return "bit string becomes " + nt
	case strings.HasPrefix(k.name, "geo:") && !strings.Contains(nt, strings.TrimPrefix(k.name, "geo:")):
		return "geometric type stored as " + nt
	case k.name == "int1" && d == MSSQL && !k.unsigned:
		return ""
	}
	return ""
}
