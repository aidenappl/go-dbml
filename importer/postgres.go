// Portions of this file are derived from holistics/dbml
// packages/dbml-core/src/parse/ANTLR/ASTGeneration/postgres/PostgresASTGen.js
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: re-implemented in Go over a hand-written tokenizer and
// recursive-descent statement parser (no ANTLR grammar); the visitor's
// semantics (what is kept, how text is rendered) are mirrored.

package importer

import (
	"strconv"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

type pgImporter struct {
	raw rawDB
	err error
	cfg *config
}

func importPostgres(src string, cfg *config) (*rawDB, error) {
	im := &pgImporter{cfg: cfg}
	forEachStatement(src, dPostgres, false, cfg.skipStatement, cfg.maxStatementTokens(), func(st statement) {
		if im.err != nil {
			return
		}
		c := newCursor(src, st.toks)
		c.truncated = st.truncated
		im.statement(c)
	})
	if im.err != nil {
		return nil, im.err
	}
	return &im.raw, nil
}

func (im *pgImporter) statement(c *cursor) {
	switch {
	case c.peek().is("CREATE"):
		c.next()
		c.acceptSeq("OR", "REPLACE")
		for c.accept("GLOBAL") || c.accept("LOCAL") || c.accept("TEMPORARY") || c.accept("TEMP") || c.accept("UNLOGGED") {
		}
		switch {
		case c.accept("TABLE"):
			im.createTable(c)
		case c.peek().is("UNIQUE") && c.at(1).is("INDEX"):
			c.i += 2
			im.createIndex(c, true)
		case c.accept("INDEX"):
			im.createIndex(c, false)
		case c.accept("TYPE"):
			im.createType(c)
		}
	case c.peek().is("ALTER") && c.at(1).is("TABLE"):
		c.i += 2
		im.alterTable(c)
	case c.peek().is("COMMENT") && c.at(1).is("ON"):
		c.i += 2
		im.comment(c)
	case c.peek().is("INSERT") && c.at(1).is("INTO"):
		c.i += 2
		im.insert(c)
	}
}

// pgIdent mirrors visitColid/visitIdentifier: quoted identifiers lose
// their quotes (no unescaping), words keep their text.
func pgIdent(t token) string {
	if t.kind == tQIdent {
		if strings.HasPrefix(t.text, "U&") {
			return stripQuotes(t.text[2:])
		}
		return stripQuotes(t.text)
	}
	return t.text
}

func isNameTok(t token) bool { return t.kind == tWord || t.kind == tQIdent }

// qualifiedName: colid ('.' attr_name)*
func (im *pgImporter) qualifiedName(c *cursor) []string {
	var names []string
	for isNameTok(c.peek()) {
		names = append(names, im.ident(c.next()))
		if !c.peek().isPunct(".") {
			break
		}
		c.next()
	}
	return names
}

func (im *pgImporter) createTable(c *cursor) {
	c.acceptSeq("IF", "NOT", "EXISTS")
	schema, name := lastTwo(im.qualifiedName(c))
	if name == "" || !c.peek().isPunct("(") || c.at(1).isPunct(")") {
		return // PARTITION OF / OF type / AS SELECT / (): no table elements
	}
	c.next()
	t := &model.Table{Name: name, SchemaName: schema}
	var inline, tableRefs []*model.Ref
	for !c.eof() && !c.peek().isPunct(")") {
		start := c.i
		im.tableElement(c, t, &inline, &tableRefs)
		c.skipUntil(boundary)
		if !c.acceptPunct(",") {
			if c.i == start {
				c.next()
			}
			break
		}
	}
	for _, r := range inline {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
	}
	for _, r := range tableRefs {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
	}
	im.raw.refs = append(im.raw.refs, inline...)
	im.raw.refs = append(im.raw.refs, tableRefs...)
	im.raw.addTable(t)
}

func (im *pgImporter) tableElement(c *cursor, t *model.Table, inline, tableRefs *[]*model.Ref) {
	p := c.peek()
	switch {
	case p.is("LIKE"):
		return
	case p.is("CONSTRAINT") || p.is("CHECK") || p.is("UNIQUE") || p.is("PRIMARY") || p.is("EXCLUDE") || p.is("FOREIGN"):
		k := im.tableConstraint(c)
		switch {
		case k.index != nil:
			t.Indexes = append(t.Indexes, k.index)
		case k.ref != nil:
			*tableRefs = append(*tableRefs, k.ref)
		case k.check != nil:
			t.Checks = append(t.Checks, k.check)
		}
	default:
		im.columnDef(c, t, inline)
	}
}

// pgConstraint is one parsed table constraint; exactly one of index, ref,
// check is set (none for unsupported or malformed constraints).
type pgConstraint struct {
	index *model.Index
	ref   *model.Ref
	check *model.Check
}

// tableConstraint: [CONSTRAINT name] constraintelem
func (im *pgImporter) tableConstraint(c *cursor) pgConstraint {
	var name *string
	if c.accept("CONSTRAINT") {
		name = strp(im.ident(c.next()))
	}
	switch {
	case c.accept("CHECK"):
		if !c.peek().isPunct("(") {
			return pgConstraint{}
		}
		from, to := c.skipGroup()
		return pgConstraint{check: &model.Check{Expression: c.origText(from+1, to-1), Name: name}}
	case c.accept("UNIQUE"):
		c.acceptSeq("NULLS", "NOT", "DISTINCT")
		c.acceptSeq("NULLS", "DISTINCT")
		if !c.peek().isPunct("(") {
			return pgConstraint{}
		}
		return pgConstraint{index: &model.Index{Name: name, Unique: boolp(true), Columns: im.columnList(c)}}
	case c.peek().is("PRIMARY") && c.at(1).is("KEY"):
		c.i += 2
		if !c.peek().isPunct("(") {
			return pgConstraint{}
		}
		return pgConstraint{index: &model.Index{Name: name, PK: boolp(true), Columns: im.columnList(c)}}
	case c.peek().is("FOREIGN") && c.at(1).is("KEY"):
		c.i += 2
		first := colValues(im.columnList(c))
		if !c.accept("REFERENCES") {
			return pgConstraint{}
		}
		r := im.references(c)
		r.Name = name
		r.Endpoints[0].FieldNames = first
		return pgConstraint{ref: r}
	}
	return pgConstraint{}
}

// columnList: '(' colid (',' colid)* ')'
func (im *pgImporter) columnList(c *cursor) []*model.IndexColumn {
	var cols []*model.IndexColumn
	if !c.acceptPunct("(") {
		return cols
	}
	for !c.eof() && !c.peek().isPunct(")") {
		t := c.next()
		cols = append(cols, &model.IndexColumn{Type: "column", Value: im.ident(t)})
		c.skipUntil(boundary)
		if !c.acceptPunct(",") {
			break
		}
	}
	c.acceptPunct(")")
	return cols
}

// references: qualified_name [(cols)] [MATCH x] [key_actions]
func (im *pgImporter) references(c *cursor) *model.Ref {
	schema, table := lastTwo(im.qualifiedName(c))
	var second []string
	if c.peek().isPunct("(") {
		second = colValues(im.columnList(c))
	}
	if c.accept("MATCH") {
		c.next()
	}
	r := &model.Ref{}
	for c.peek().is("ON") && (c.at(1).is("DELETE") || c.at(1).is("UPDATE")) {
		c.next()
		which := c.next()
		a := c.next()
		action := a.text
		if a.is("SET") || a.is("NO") {
			action += " " + c.next().text
		}
		if which.is("DELETE") {
			r.OnDelete = strp(action)
		} else {
			r.OnUpdate = strp(action)
		}
	}
	r.Endpoints[0] = model.RefEndpoint{Relation: "0..*"}
	r.Endpoints[1] = model.RefEndpoint{TableName: table, SchemaName: schema, FieldNames: second, Relation: "0..1"}
	return r
}

func (im *pgImporter) columnDef(c *cursor, t *model.Table, inline *[]*model.Ref) {
	nameTok := c.next()
	if !isNameTok(nameTok) {
		return
	}
	f := &model.Column{Name: im.ident(nameTok)}
	f.Type = im.typename(c)
	switch strings.ToLower(f.Type.TypeName) {
	case "serial", "smallserial", "bigserial":
		f.Increment = boolp(true)
	}
	for !c.eof() && !boundary(c.peek()) {
		start := c.i
		im.colConstraint(c, f, inline)
		if c.i == start {
			c.next()
		}
	}
	t.Fields = append(t.Fields, f)
}

func (im *pgImporter) colConstraint(c *cursor, f *model.Column, inline *[]*model.Ref) {
	if c.accept("CONSTRAINT") {
		c.next() // the name is not kept for column constraints
	}
	switch {
	case c.peek().is("NOT") && c.at(1).is("NULL"):
		c.i += 2
		f.NotNull = boolp(true)
	case c.accept("NULL"):
		f.NotNull = boolp(false)
	case c.accept("UNIQUE"):
		c.acceptSeq("NULLS", "NOT", "DISTINCT")
		c.acceptSeq("NULLS", "DISTINCT")
		f.Unique = boolp(true)
	case c.peek().is("PRIMARY") && c.at(1).is("KEY"):
		c.i += 2
		f.PK = boolp(true)
	case c.accept("CHECK"):
		if c.peek().isPunct("(") {
			from, to := c.skipGroup()
			f.Checks = append(f.Checks, &model.Check{Expression: c.origText(from+1, to-1)})
		}
		c.acceptSeq("NO", "INHERIT")
	case c.accept("DEFAULT"):
		f.DBDefault = im.bExpr(c)
	case c.accept("GENERATED"):
		if !c.accept("ALWAYS") {
			c.acceptSeq("BY", "DEFAULT")
		}
		c.accept("AS")
		identity := c.accept("IDENTITY")
		if identity {
			f.Increment = boolp(true)
		}
		if c.peek().isPunct("(") {
			from, to := c.skipGroup()
			if im.cfg.clean && !identity {
				setGenerated(f, c.origText(from+1, to-1), true)
			}
		}
		c.accept("STORED")
	case c.accept("REFERENCES"):
		r := im.references(c)
		r.Endpoints[0].FieldNames = []string{f.Name}
		*inline = append(*inline, r)
	case c.accept("COLLATE"):
		im.qualifiedName(c)
	case c.peek().isPunct("("):
		c.skipGroup()
	default:
		c.next()
	}
}

var bExprStop = func(t token) bool {
	switch {
	case t.isPunct(","), t.kind == tTerm:
		return true
	case t.kind == tWord:
		switch strings.ToUpper(t.text) {
		case "NOT", "NULL", "UNIQUE", "PRIMARY", "CHECK", "DEFAULT", "GENERATED", "REFERENCES", "COLLATE", "CONSTRAINT",
			"DEFERRABLE", "INITIALLY":
			return true
		}
	}
	return false
}

// bExpr mirrors visitB_expr / visitC_expr_expr / visitAexprconst.
func (im *pgImporter) bExpr(c *cursor) *model.Default {
	from, to := c.skipUntil(bExprStop)
	n := to - from
	if n == 0 {
		return nil
	}
	first := c.toks[from]
	// a single c_expr: constant, (a_expr), function call or column ref
	if n == 1 {
		return pgConst(first, im.text(c, from, to), im.cfg.clean)
	}
	if first.isPunct("(") {
		c.i = from
		if _, end := c.skipGroup(); end == to {
			c.i = to
			return &model.Default{Type: "expression", Value: im.text(c, from+1, to-1)}
		}
		c.i = to
	}
	return &model.Default{Type: "expression", Value: im.text(c, from, to)}
}

// pgConst mirrors visitAexprconst for a single token.
func pgConst(t token, text string, clean bool) *model.Default {
	switch {
	case t.kind == tString && !strings.ContainsRune("BbXx", rune(t.text[0])):
		if clean {
			return &model.Default{Type: "string", Value: unescapeStdString(t.text)}
		}
		return &model.Default{Type: "string", Value: anysconst(t.text)}
	case t.is("TRUE") || t.is("FALSE") || t.is("NULL"):
		return &model.Default{Type: "boolean", Value: strings.ToLower(t.text)}
	case t.kind == tNumber:
		if v, ok := iconst(t.text); ok {
			return &model.Default{Type: "number", Value: v}
		}
		return &model.Default{Type: "number", Value: t.text}
	}
	return &model.Default{Type: "expression", Value: text}
}

// anysconst mirrors visitAnysconst: the token text minus its first and last
// character (E'..' and U&'..' keep their prefix's quote, like the reference).
func anysconst(s string) string {
	if strings.HasPrefix(s, "$") {
		// $tag$...$tag$: the reference slices one character off each end
		return stripQuotes(s)
	}
	return stripQuotes(s)
}

// iconst mirrors Number(text).toString() for an integer literal; ok is
// false for decimal/float literals (fconst keeps its text).
func iconst(s string) (string, bool) {
	low := strings.ToLower(strings.ReplaceAll(s, "_", ""))
	base := 10
	digits := low
	switch {
	case strings.HasPrefix(low, "0x"):
		base, digits = 16, low[2:]
	case strings.HasPrefix(low, "0o"):
		base, digits = 8, low[2:]
	case strings.HasPrefix(low, "0b"):
		base, digits = 2, low[2:]
	default:
		if strings.ContainsAny(low, ".e") {
			return "", false
		}
	}
	n, err := strconv.ParseUint(digits, base, 64)
	if err != nil {
		f, err := strconv.ParseFloat(digits, 64)
		if err != nil || base != 10 {
			return s, true
		}
		return jsNumberString(f), true
	}
	if n > 1<<53 {
		return jsNumberString(float64(n)), true
	}
	return strconv.FormatUint(n, 10), true
}

func jsNumberString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// ---------------------------------------------------------------- types

var pgNumericWords = map[string]bool{"INT": true, "INTEGER": true, "SMALLINT": true, "BIGINT": true, "REAL": true, "FLOAT": true,
	"DOUBLE": true, "DECIMAL": true, "DEC": true, "NUMERIC": true, "BOOLEAN": true}

// typename mirrors visitTypename / visitSimpletypename.
func (im *pgImporter) typename(c *cursor) model.ColumnType {
	c.accept("SETOF")
	start := c.i
	t := c.peek()
	up := t.upper()
	var typ string
	var schema *string
	switch {
	case t.kind == tWord && pgNumericWords[up]:
		c.next()
		if up == "DOUBLE" {
			c.accept("PRECISION")
		}
		if c.peek().isPunct("(") {
			if !im.cfg.clean && up != "FLOAT" && up != "DECIMAL" && up != "DEC" && up != "NUMERIC" {
				// INT(11) etc. are syntax errors in the reference grammar
				im.err = errUnsupported("type modifier on " + t.text)
				return model.ColumnType{}
			}
			c.skipGroup()
		}
		typ = im.typeText(c, start, c.i)
	case t.kind == tWord && (up == "CHARACTER" || up == "CHAR" || up == "NCHAR" || up == "VARCHAR" || up == "NATIONAL"):
		var words []string
		words = append(words, c.next().text)
		if up == "NATIONAL" {
			words = append(words, c.next().text)
		}
		if up != "VARCHAR" && c.peek().is("VARYING") {
			words = append(words, c.next().text)
		}
		typ = strings.Join(words, " ")
		if c.peek().isPunct("(") && c.at(1).kind == tNumber && c.at(2).isPunct(")") {
			n, _ := iconst(c.at(1).text)
			typ += "(" + n + ")"
			c.i += 3
		}
	case t.kind == tWord && up == "BIT":
		c.next()
		c.accept("VARYING")
		if c.peek().isPunct("(") {
			c.skipGroup()
		}
		typ = im.typeText(c, start, c.i)
	case t.kind == tWord && (up == "TIMESTAMP" || up == "TIME"):
		c.next()
		typ = t.text
		if c.peek().isPunct("(") && c.at(1).kind == tNumber && c.at(2).isPunct(")") {
			n, _ := iconst(c.at(1).text)
			typ += "(" + n + ")"
			c.i += 3
		}
		if c.peek().is("WITH") || c.peek().is("WITHOUT") {
			with := c.next().is("WITH")
			c.accept("TIME")
			c.accept("ZONE")
			if im.cfg.clean && with {
				typ += " with time zone"
			}
		}
	case t.kind == tWord && up == "INTERVAL":
		c.next()
		for c.peek().is("YEAR") || c.peek().is("MONTH") || c.peek().is("DAY") || c.peek().is("HOUR") || c.peek().is("MINUTE") || c.peek().is("SECOND") || c.peek().is("TO") {
			c.next()
		}
		if c.peek().isPunct("(") {
			c.skipGroup()
		}
		typ = im.typeText(c, start, c.i)
	case t.kind == tWord && pgReserved[up] && !im.cfg.clean:
		im.err = errUnsupported("reserved keyword " + t.text + " as a type name")
		return model.ColumnType{}
	case isNameTok(t):
		c.next()
		if c.peek().isPunct(".") {
			// generictype with attrs: the first part is getText (quotes kept)
			names := []string{t.text}
			if im.cfg.clean {
				names[0] = im.ident(t)
			}
			for c.peek().isPunct(".") && isNameTok(c.at(1)) {
				c.next()
				names = append(names, im.ident(c.next()))
			}
			schema, typ = lastTwo(names)
		} else {
			typ = im.ident(t)
			if im.cfg.clean && t.kind == tWord {
				typ = strings.ToLower(typ) // unquoted type names fold to lower case
			}
		}
		if c.peek().isPunct("(") {
			from, to := c.skipGroup()
			if schema == nil {
				typ += c.getText(from, to)
			}
		}
	default:
		return model.ColumnType{}
	}
	// array bounds: ([n]?)* | ARRAY [n]?
	if c.peek().isPunct("[") {
		from := c.i
		for c.peek().isPunct("[") {
			c.next()
			if c.peek().kind == tNumber {
				c.next()
			}
			c.acceptPunct("]")
		}
		typ += c.getText(from, c.i)
	} else if c.accept("ARRAY") {
		n := ""
		if c.acceptPunct("[") {
			if c.peek().kind == tNumber {
				n, _ = iconst(c.next().text)
			}
			c.acceptPunct("]")
		}
		typ += "[" + n + "]"
	}
	return model.ColumnType{TypeName: typ, SchemaName: schema}
}

// ---------------------------------------------------------------- other statements

func (im *pgImporter) createIndex(c *cursor, unique bool) {
	c.accept("CONCURRENTLY")
	var name *string
	c.acceptSeq("IF", "NOT", "EXISTS")
	if !c.peek().is("ON") {
		name = strp(im.ident(c.next()))
	}
	if !c.accept("ON") {
		return
	}
	c.accept("ONLY")
	schema, table := lastTwo(im.qualifiedName(c))
	c.acceptPunct("*")
	var typ *string
	if c.accept("USING") {
		typ = strp(im.ident(c.next()))
	}
	if !c.peek().isPunct("(") {
		return
	}
	c.next()
	var cols []*model.IndexColumn
	for !c.eof() && !c.peek().isPunct(")") {
		from, to := c.skipUntil(func(t token) bool { return t.isPunct(",") })
		if to > from && isNameTok(c.toks[from]) && (to-from == 1 || !c.toks[from+1].isPunct("(") && !c.toks[from+1].isPunct(".")) {
			cols = append(cols, &model.IndexColumn{Type: "string", Value: im.ident(c.toks[from])})
		} else if to > from {
			cols = append(cols, &model.IndexColumn{Type: "expression", Value: im.text(c, from, to)})
		}
		if !c.acceptPunct(",") {
			break
		}
	}
	c.acceptPunct(")")
	t := im.raw.findTable(schema, table)
	if t == nil {
		return
	}
	t.Indexes = append(t.Indexes, &model.Index{Name: name, Unique: boolp(unique), Type: typ, Columns: cols})
}

func (im *pgImporter) createType(c *cursor) {
	schema, name := lastTwo(im.qualifiedName(c))
	if !c.accept("AS") || !c.accept("ENUM") || !c.peek().isPunct("(") {
		return
	}
	c.next()
	en := &model.Enum{Name: name, SchemaName: schema}
	for !c.eof() && !c.peek().isPunct(")") {
		if t := c.next(); t.kind == tString {
			en.Values = append(en.Values, &model.EnumValue{Name: im.sconst(t.text)})
		}
		c.acceptPunct(",")
	}
	im.raw.enums = append(im.raw.enums, en)
}

func (im *pgImporter) alterTable(c *cursor) {
	c.acceptSeq("IF", "EXISTS")
	c.accept("ONLY")
	names := im.qualifiedName(c)
	c.acceptPunct("*")
	schema, name := lastTwo(names)
	for !c.eof() {
		start := c.i
		if c.accept("ADD") {
			k := im.tableConstraint(c)
			switch {
			case k.ref != nil:
				r := k.ref
				r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
				im.raw.refs = append(im.raw.refs, r)
			case k.check != nil:
				if t := im.raw.findTable(schema, name); t != nil {
					t.Checks = append(t.Checks, k.check)
				}
			case k.index != nil:
				ix := k.index
				if len(ix.Columns) == 0 || !(truthy(ix.PK) || truthy(ix.Unique)) {
					break
				}
				t := im.raw.findTable(schema, name)
				if t == nil {
					break
				}
				if len(ix.Columns) == 1 {
					if f := findField(t, ix.Columns[0].Value); f != nil {
						if truthy(ix.Unique) {
							f.Unique = boolp(true)
						} else {
							f.PK = boolp(true)
						}
						break
					}
				}
				t.Indexes = append(t.Indexes, ix)
			}
		}
		c.skipUntil(func(t token) bool { return t.isPunct(",") })
		if !c.acceptPunct(",") {
			if c.i == start {
				c.next()
			}
			if c.eof() {
				break
			}
		}
	}
}

func truthy(b *bool) bool { return b != nil && *b }

func (im *pgImporter) comment(c *cursor) {
	switch {
	case c.accept("TABLE"):
		schema, name := lastTwo(im.qualifiedName(c))
		t := im.raw.findTable(schema, name)
		if t == nil || !c.accept("IS") {
			return
		}
		if v := c.next(); v.kind == tString {
			t.Note = &model.NoteValue{Value: im.commentText(v.text)}
		} else {
			t.Note = nil
		}
	case c.accept("COLUMN"):
		names := im.qualifiedName(c)
		if len(names) == 0 {
			return
		}
		field := names[len(names)-1]
		var table string
		var schema *string
		if len(names) > 1 {
			table = names[len(names)-2]
		}
		if len(names) > 2 {
			schema = strp(names[len(names)-3])
		}
		t := im.raw.findTable(schema, table)
		if t == nil {
			// the reference dereferences the missing table and throws
			if !im.cfg.clean {
				im.err = errUnsupported("COMMENT ON COLUMN of an unknown table " + table)
			}
			return
		}
		f := findField(t, field)
		if f == nil || !c.accept("IS") {
			return
		}
		if v := c.next(); v.kind == tString {
			f.Note = &model.NoteValue{Value: im.commentText(v.text)}
		} else {
			f.Note = nil
		}
	}
}

func (im *pgImporter) insert(c *cursor) {
	schema, table := lastTwo(im.qualifiedName(c))
	if c.accept("AS") {
		c.next()
	}
	var cols []string
	if c.peek().isPunct("(") && !c.at(1).is("SELECT") && !c.at(1).is("VALUES") {
		c.next()
		for !c.eof() && !c.peek().isPunct(")") {
			cols = append(cols, im.ident(c.next()))
			c.skipUntil(boundary)
			if !c.acceptPunct(",") {
				break
			}
		}
		c.acceptPunct(")")
	}
	if c.accept("OVERRIDING") {
		c.next()
		c.accept("VALUE")
	}
	if c.peek().is("DEFAULT") && c.at(1).is("VALUES") {
		if !im.cfg.clean {
			im.err = errUnsupported("INSERT ... DEFAULT VALUES")
		}
		return
	}
	if len(cols) == 0 {
		if t := im.raw.findTable(schema, table); t != nil {
			for _, f := range t.Fields {
				cols = append(cols, f.Name)
			}
		}
	}
	if cols == nil {
		cols = []string{}
	}
	rec := &model.TableRecord{SchemaName: schema, TableName: table, Columns: cols}
	if c.accept("VALUES") {
		for c.peek().isPunct("(") {
			c.next()
			row := []*model.RecordValue{}
			for !c.eof() && !c.peek().isPunct(")") {
				if c.peek().isPunct(",") {
					if !im.cfg.clean {
						im.err = errUnsupported("empty value in VALUES")
					}
					return
				}
				from, to := c.skipUntil(boundary)
				row = append(row, im.exprValues(c, from, to)...)
				if !c.acceptPunct(",") {
					break
				}
			}
			c.acceptPunct(")")
			if c.truncated && c.eof() {
				im.cfg.truncated = true // the statement was cut mid-row
				break
			}
			if !im.cfg.addRow() {
				break
			}
			rec.Values = append(rec.Values, row)
			if !c.acceptPunct(",") {
				break
			}
		}
	}
	im.raw.records = append(im.raw.records, rec)
}

var pgBinaryOps = map[string]bool{"+": true, "-": true, "*": true, "/": true, "%": true, "^": true, "||": true, "<": true, ">": true,
	"=": true, "<=": true, ">=": true, "<>": true, "!=": true, "->": true, "->>": true, "#>": true, "#>>": true, "@>": true, "<@": true}

// exprValues mirrors flattenDepth(a_expr, 21) + sanitize for one VALUES
// item: one value per operand of top-level binary operators (unary signs
// and ::casts dropped), undefined entries removed.
func (im *pgImporter) exprValues(c *cursor, from, to int) []*model.RecordValue {
	var out []*model.RecordValue
	depth := 0
	start := from
	flush := func(end int) {
		// strip leading unary signs
		for start < end && (c.toks[start].isPunct("-") || c.toks[start].isPunct("+")) {
			start++
		}
		// strip a trailing ::type cast (and COLLATE x)
		e := end
		d := 0
		for i := start; i < end; i++ {
			t := c.toks[i]
			if t.isPunct("(") || t.isPunct("[") {
				d++
			} else if t.isPunct(")") || t.isPunct("]") {
				d--
			} else if d == 0 && (t.isPunct("::") || t.is("COLLATE")) {
				e = i
				break
			}
		}
		if start >= e {
			return
		}
		out = append(out, im.value(c, start, e))
	}
	for i := from; i < to; i++ {
		t := c.toks[i]
		if t.isPunct("(") || t.isPunct("[") || t.is("CASE") {
			depth++
			continue
		}
		if t.isPunct(")") || t.isPunct("]") || t.is("END") {
			depth--
			continue
		}
		if depth == 0 && i > start && t.kind == tPunct && pgBinaryOps[t.text] && !c.toks[i-1].isPunct("::") {
			flush(i)
			start = i + 1
		}
	}
	flush(to)
	return out
}

func (im *pgImporter) value(c *cursor, from, to int) *model.RecordValue {
	text := im.text(c, from, to)
	if to-from == 1 {
		d := pgConst(c.toks[from], text, im.cfg.clean)
		return &model.RecordValue{Value: d.Value, Type: d.Type}
	}
	if c.toks[from].isPunct("(") && c.matchOf(from) == to-1 {
		return &model.RecordValue{Value: im.text(c, from+1, to-1), Type: "expression"}
	}
	return &model.RecordValue{Value: text, Type: "expression"}
}

// pgReserved are PostgreSQL reserved keywords (not usable as type names).
var pgReserved = func() map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(`ALL ANALYSE ANALYZE AND ANY ARRAY AS ASC ASYMMETRIC BOTH CASE CAST CHECK COLLATE COLUMN
		CONSTRAINT CREATE CURRENT_CATALOG CURRENT_DATE CURRENT_ROLE CURRENT_TIME CURRENT_TIMESTAMP CURRENT_USER DEFAULT
		DEFERRABLE DESC DISTINCT DO ELSE END EXCEPT FALSE FETCH FOR FOREIGN FROM GRANT GROUP HAVING IN INITIALLY INTERSECT
		INTO LATERAL LEADING LIMIT LOCALTIME LOCALTIMESTAMP NOT NULL OFFSET ON ONLY OR ORDER PLACING PRIMARY REFERENCES
		RETURNING SELECT SESSION_USER SOME SYMMETRIC TABLE THEN TO TRAILING TRUE UNION UNIQUE USER USING VARIADIC WHEN
		WHERE WINDOW WITH`) {
		m[w] = true
	}
	return m
}()

// ident decodes an identifier token (pgIdent in parity mode; clean mode
// also un-doubles "" inside quoted identifiers).
func (im *pgImporter) ident(t token) string {
	if im.cfg.clean && t.kind == tQIdent && !strings.HasPrefix(t.text, "U&") {
		return unquoteIdent(t.text)
	}
	return pgIdent(t)
}

// sconst decodes a string constant (anysconst in parity mode).
func (im *pgImporter) sconst(text string) string {
	if im.cfg.clean {
		return unescapeStdString(text)
	}
	return anysconst(text)
}

// comment decodes a COMMENT ON text (the reference un-doubles ” only).
func (im *pgImporter) commentText(text string) string {
	if im.cfg.clean {
		return unescapeStdString(text)
	}
	return strings.ReplaceAll(anysconst(text), "''", "'")
}

// text renders expression tokens: getText (parity) or source text (clean).
func (im *pgImporter) text(c *cursor, from, to int) string {
	if im.cfg.clean {
		return c.origText(from, to)
	}
	return c.getText(from, to)
}

// typeText renders type tokens [from, to): getText (parity, words run
// together as the reference does) or the words separated by one space
// (clean: "double precision", "bit varying(8)").
func (im *pgImporter) typeText(c *cursor, from, to int) string {
	if !im.cfg.clean {
		return c.getText(from, to)
	}
	var b strings.Builder
	for i := from; i < to; i++ {
		t := c.toks[i]
		if i > from && t.kind == tWord && c.toks[i-1].kind == tWord {
			b.WriteByte(' ')
		}
		b.WriteString(strings.ToLower(t.text))
	}
	return b.String()
}
