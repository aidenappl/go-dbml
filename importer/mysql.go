// Portions of this file are derived from holistics/dbml
// packages/dbml-core/src/parse/ANTLR/ASTGeneration/mysql/MySQLASTGen.js
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: re-implemented in Go over a hand-written tokenizer and
// recursive-descent statement parser (no ANTLR grammar); the visitor's
// semantics (what is kept, how text is rendered) are mirrored.

package importer

import (
	"regexp"
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

// mysqlImporter mirrors MySQLASTGen.
type mysqlImporter struct {
	lastIndexed *model.Table
	raw         rawDB
	err         error
	cfg         *config
}

func importMySQL(src string, cfg *config) (*rawDB, error) {
	im := &mysqlImporter{cfg: cfg}
	forEachStatement(src, dMySQL, true, cfg.skipStatement, cfg.maxStatementTokens(), func(st statement) {
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

func (im *mysqlImporter) statement(c *cursor) {
	switch {
	case c.peek().is("CREATE"):
		c.next()
		c.accept("TEMPORARY")
		save := c.i
		if c.accept("TABLE") {
			im.createTable(c)
			return
		}
		c.i = save
		c.accept("ONLINE")
		c.accept("OFFLINE")
		unique := false
		kind := ""
		if c.accept("UNIQUE") {
			unique = true
		} else if c.peek().is("FULLTEXT") || c.peek().is("SPATIAL") {
			kind = strings.ToLower(c.next().text)
		}
		if c.accept("INDEX") {
			im.createIndex(c, unique)
			if t := im.lastIndexed; im.cfg.clean && kind != "" && t != nil && len(t.Indexes) > 0 {
				t.Indexes[len(t.Indexes)-1].Type = strp(kind)
			}
		}
	case c.peek().is("ALTER"):
		c.next()
		c.accept("ONLINE")
		c.accept("OFFLINE")
		c.accept("IGNORE")
		if c.accept("TABLE") {
			im.alterTable(c)
		}
	case c.peek().is("INSERT"):
		c.next()
		im.insert(c)
	}
}

// fullId: uid (DOT_ID | '.' uid)?  -> name parts
func (im *mysqlImporter) fullID(c *cursor) []string {
	var names []string
	for {
		t := c.peek()
		if t.kind != tWord && t.kind != tQIdent && t.kind != tString {
			break
		}
		c.next()
		if t.text == "``" && !im.cfg.clean {
			// an empty quoted identifier is a syntax error in the reference
			im.err = errUnsupported("empty identifier")
		}
		names = append(names, im.uid(t))
		if !c.peek().isPunct(".") {
			break
		}
		c.next()
	}
	return names
}

func (im *mysqlImporter) uid(t token) string {
	if im.cfg.clean {
		if t.kind == tQIdent {
			return unquoteIdent(t.text)
		}
		if t.kind == tString {
			return unescapeMySQLString(t.text)
		}
		return t.text
	}
	if t.kind == tQIdent && t.text == "``" {
		return ""
	}
	if t.kind == tQIdent || t.kind == tString {
		return stripQuotes(t.text)
	}
	return t.text
}

func (im *mysqlImporter) createTable(c *cursor) {
	if c.acceptSeq("IF", "NOT", "EXISTS") {
	}
	names := im.fullID(c)
	schema, name := lastTwo(names)
	if name == "" {
		return
	}
	if c.peek().is("LIKE") || (c.peek().isPunct("(") && c.at(1).is("LIKE")) {
		return // copyCreateTable: not supported
	}
	if !c.peek().isPunct("(") {
		return // queryCreateTable without definitions
	}
	c.next()
	t := &model.Table{Name: name, SchemaName: schema}
	var inlineRefs, tableRefs []*model.Ref
	var singlePK *model.Index
	for !c.eof() && !c.peek().isPunct(")") {
		start := c.i
		im.createDefinition(c, t, &inlineRefs, &tableRefs, &singlePK)
		if c.i == start {
			c.next()
		}
		c.skipUntil(boundary)
		if !c.acceptPunct(",") {
			break
		}
	}
	c.acceptPunct(")")
	// table options
	for !c.eof() {
		if c.accept("COMMENT") {
			c.acceptPunct("=")
			if s := c.peek(); s.kind == tString {
				c.next()
				t.Note = &model.NoteValue{Value: im.str(s.text)}
			}
			continue
		}
		if c.peek().is("PARTITION") || c.peek().is("AS") || c.peek().is("SELECT") || c.peek().is("IGNORE") || c.peek().is("REPLACE") {
			break
		}
		c.next()
	}

	for _, f := range t.Fields {
		if strings.EqualFold(f.Type.TypeName, "enum") && f.Type.Args != nil {
			im.inlineEnum(t, f)
		}
	}
	im.raw.refs = append(im.raw.refs, inlineRefs...)
	for _, r := range tableRefs {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
	}
	im.raw.refs = append(im.raw.refs, tableRefs...)
	if singlePK != nil && len(singlePK.Columns) > 0 {
		if f := findField(t, singlePK.Columns[0].Value); f != nil {
			f.PK = boolp(true)
		}
	}
	im.raw.addTable(t)
}

var quoteChars = regexp.MustCompile("'|\"|`")

// inlineEnum turns ENUM('a','b') into a `<table>_<column>_enum` Enum.
func (im *mysqlImporter) inlineEnum(t *model.Table, f *model.Column) {
	en := &model.Enum{Name: t.Name + "_" + f.Name + "_enum", SchemaName: t.SchemaName}
	for _, a := range splitArgs(*f.Type.Args) {
		name := strings.TrimSpace(quoteChars.ReplaceAllString(a, ""))
		if im.cfg.clean {
			name = unescapeMySQLString(a)
		}
		en.Values = append(en.Values, &model.EnumValue{Name: name})
	}
	f.Type = model.ColumnType{TypeName: en.Name, SchemaName: en.SchemaName}
	im.raw.enums = append(im.raw.enums, en)
}

// splitArgs splits the raw collection options (stored NUL-separated).
func splitArgs(s string) []string { return strings.Split(s, "\x00") }

func (im *mysqlImporter) createDefinition(c *cursor, t *model.Table, inlineRefs, tableRefs *[]*model.Ref, singlePK **model.Index) {
	p := c.peek()
	switch {
	case p.is("CONSTRAINT") || p.is("PRIMARY") || p.is("UNIQUE") || p.is("FOREIGN") || p.is("CHECK"):
		var name *string
		if c.accept("CONSTRAINT") {
			if n := c.peek(); !(n.is("PRIMARY") || n.is("UNIQUE") || n.is("FOREIGN") || n.is("CHECK")) {
				c.next()
				name = strp(im.uid(n))
			}
		}
		switch {
		case c.accept("PRIMARY"):
			c.accept("KEY")
			ix := im.indexTail(c, name, true)
			ix.PK = boolp(true)
			if len(ix.Columns) > 1 {
				t.Indexes = append(t.Indexes, ix)
			} else {
				*singlePK = ix
			}
		case c.accept("UNIQUE"):
			if !c.accept("INDEX") {
				c.accept("KEY")
			}
			ix := im.indexTail(c, name, true)
			ix.Unique = boolp(true)
			t.Indexes = append(t.Indexes, ix)
		case c.accept("FOREIGN"):
			c.accept("KEY")
			if r := im.foreignKey(c, name); r != nil {
				*tableRefs = append(*tableRefs, r)
			}
		case c.accept("CHECK"):
			if ch := im.check(c, name); ch != nil {
				t.Checks = append(t.Checks, ch)
			}
		}
	case p.is("INDEX") || p.is("KEY"):
		c.next()
		t.Indexes = append(t.Indexes, im.indexTail(c, nil, true))
	case p.is("FULLTEXT") || p.is("SPATIAL"):
		kind := strings.ToLower(c.next().text)
		if !c.accept("INDEX") {
			c.accept("KEY")
		}
		ix := im.indexTail(c, nil, false)
		if im.cfg.clean {
			ix.Type = strp(kind)
		}
		t.Indexes = append(t.Indexes, ix)
	default:
		im.column(c, t, inlineRefs)
	}
}

// indexTail parses `[uid] [USING x] (cols) indexOption*`. The uid is used
// as the name when name is nil. allowType controls the leading USING.
func (im *mysqlImporter) indexTail(c *cursor, name *string, allowType bool) *model.Index {
	ix := &model.Index{}
	if t := c.peek(); (t.kind == tWord || t.kind == tQIdent) && !t.is("USING") && !c.peek().isPunct("(") {
		c.next()
		if name == nil {
			name = strp(im.uid(t))
		}
	}
	ix.Name = name
	var typ *string
	if allowType && c.peek().is("USING") {
		c.next()
		typ = strp(c.next().text)
	}
	ix.Columns = im.indexColumnNames(c)
	if ot := im.indexOptions(c); ot != nil {
		typ = ot
	}
	ix.Type = typ
	return ix
}

func (im *mysqlImporter) indexOptions(c *cursor) *string {
	var typ *string
	for !c.eof() && !boundary(c.peek()) {
		switch {
		case c.peek().is("USING"):
			c.next()
			typ = strp(c.next().text)
		case c.peek().is("COMMENT"):
			c.next()
			c.next()
		case c.peek().isPunct("("):
			c.skipGroup()
		default:
			if c.peek().is("ALGORITHM") || c.peek().is("LOCK") {
				return typ
			}
			c.next()
		}
	}
	return typ
}

// indexColumnNames: '(' indexColumnName (',' indexColumnName)* ')'
func (im *mysqlImporter) indexColumnNames(c *cursor) []*model.IndexColumn {
	var cols []*model.IndexColumn
	if !c.acceptPunct("(") {
		return cols
	}
	for !c.eof() && !c.peek().isPunct(")") {
		t := c.peek()
		next := c.at(1)
		switch {
		case (t.kind == tWord || t.kind == tQIdent) && (boundary(next) || next.is("ASC") || next.is("DESC") ||
			(next.isPunct("(") && c.at(2).kind == tNumber && c.at(3).isPunct(")"))):
			c.next()
			col := &model.IndexColumn{Type: "column", Value: im.uid(t)}
			if c.peek().isPunct("(") {
				from, _ := c.skipGroup()
				if im.cfg.clean && simpleIdent.MatchString(col.Value) {
					// clean: a prefix key part `col(n)` (MySQL exporter
					// writes it back as one)
					col = &model.IndexColumn{Type: "expression", Value: col.Value + "(" + c.toks[from+1].text + ")"}
				}
			}
			cols = append(cols, col)
		case t.kind == tString && (boundary(next) || next.is("ASC") || next.is("DESC") || next.isPunct("(")):
			c.next()
			cols = append(cols, &model.IndexColumn{Type: "string", Value: im.str(t.text)})
			if c.peek().isPunct("(") {
				c.skipGroup()
			}
		default:
			from, to := c.skipUntil(func(t token) bool { return t.isPunct(",") || t.is("ASC") || t.is("DESC") })
			cols = append(cols, &model.IndexColumn{Type: "expression", Value: im.expr(c, from, to)})
		}
		c.accept("ASC")
		c.accept("DESC")
		if !c.acceptPunct(",") {
			break
		}
	}
	c.acceptPunct(")")
	return cols
}

func (im *mysqlImporter) check(c *cursor, name *string) *model.Check {
	if !c.peek().isPunct("(") {
		return nil
	}
	from, to := c.skipGroup()
	ch := &model.Check{Expression: c.origText(from+1, to-1), Name: name}
	c.accept("NOT")
	c.accept("ENFORCED")
	return ch
}

// foreignKey: [index] (cols) referenceDefinition
func (im *mysqlImporter) foreignKey(c *cursor, name *string) *model.Ref {
	if t := c.peek(); (t.kind == tWord || t.kind == tQIdent) && !t.isPunct("(") {
		c.next()
		if name == nil {
			name = strp(im.uid(t))
		}
	}
	cols := im.indexColumnNames(c)
	if !c.accept("REFERENCES") {
		return nil
	}
	r := im.referenceDefinition(c)
	if r == nil {
		return nil
	}
	r.Name = name
	r.Endpoints[0].FieldNames = colValues(cols)
	return r
}

func colValues(cols []*model.IndexColumn) []string {
	out := []string{}
	for _, c := range cols {
		out = append(out, c.Value)
	}
	return out
}

// referenceDefinition: tableName indexColumnNames? (MATCH x)? referenceAction?
func (im *mysqlImporter) referenceDefinition(c *cursor) *model.Ref {
	schema, table := lastTwo(im.fullID(c))
	if !c.peek().isPunct("(") {
		return nil
	}
	cols := im.indexColumnNames(c)
	if c.accept("MATCH") {
		c.next()
	}
	r := &model.Ref{}
	for c.peek().is("ON") {
		c.next()
		which := c.next()
		action := im.referenceControlType(c)
		if which.is("DELETE") {
			r.OnDelete = strp(action)
		} else {
			r.OnUpdate = strp(action)
		}
	}
	r.Endpoints[0] = model.RefEndpoint{Relation: "0..*"}
	r.Endpoints[1] = model.RefEndpoint{TableName: table, SchemaName: schema, FieldNames: colValues(cols), Relation: "0..1"}
	return r
}

// referenceControlType: RESTRICT | CASCADE | SET NULL | NO ACTION | SET DEFAULT
func (im *mysqlImporter) referenceControlType(c *cursor) string {
	t := c.next()
	if t.is("SET") || t.is("NO") {
		return t.text + " " + c.next().text
	}
	return t.text
}

// ---------------------------------------------------------------- columns

func (im *mysqlImporter) column(c *cursor, t *model.Table, inlineRefs *[]*model.Ref) {
	names := im.fullColumnName(c)
	if len(names) == 0 {
		return
	}
	f := &model.Column{Name: names[len(names)-1]}
	f.Type = im.dataType(c)
	var refs []*model.Ref
	for !c.eof() && !boundary(c.peek()) {
		if im.cfg.clean && (c.peek().is("FIRST") || c.peek().is("AFTER")) {
			break // ALTER TABLE ... ADD/MODIFY column position
		}
		start := c.i
		im.columnConstraint(c, f, &refs)
		if c.i == start {
			c.next()
		}
	}
	for _, r := range refs {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName, r.Endpoints[0].FieldNames = t.Name, t.SchemaName, []string{f.Name}
	}
	*inlineRefs = append(*inlineRefs, refs...)
	t.Fields = append(t.Fields, f)
}

func (im *mysqlImporter) fullColumnName(c *cursor) []string {
	return im.fullID(c)
}

// dataType mirrors the dataType visitors: the type keyword (as written)
// plus its length/precision dimension (getText), ENUM args kept aside.
func (im *mysqlImporter) dataType(c *cursor) model.ColumnType {
	t := c.next()
	if t.kind == tEOF || boundary(t) {
		return model.ColumnType{}
	}
	name := t.text
	up := strings.ToUpper(name)
	if !im.cfg.clean && (t.kind != tWord || !mysqlDataTypes[up]) {
		// the reference grammar only accepts MySQL type keywords
		im.err = errUnsupported("MySQL data type " + t.text)
		return model.ColumnType{}
	}
	switch up {
	case "NATIONAL":
		name = c.next().text
		c.accept("VARYING")
	case "NCHAR":
		if c.peek().is("VARCHAR") {
			name = c.next().text
		} else {
			c.accept("VARYING")
		}
	case "LONG":
		if c.peek().is("VARBINARY") {
			name = t.text + " " + c.next().text
		} else {
			c.accept("VARCHAR")
		}
	case "DOUBLE":
		c.accept("PRECISION")
	case "CHAR", "CHARACTER", "VARCHAR", "NVARCHAR":
		c.accept("VARYING")
	}
	if c.peek().isPunct("(") {
		from, to := c.skipGroup()
		if !im.cfg.clean && up != "ENUM" && up != "SET" && !numericDims(c.inner(from, to)) {
			im.err = errUnsupported("MySQL type length " + c.getText(from, to))
			return model.ColumnType{}
		}
		if up == "ENUM" {
			var args []string
			for i := from + 1; i < to-1; i++ {
				if c.toks[i].kind == tString {
					args = append(args, c.toks[i].text)
				}
			}
			a := strings.Join(args, "\x00")
			return model.ColumnType{TypeName: name, Args: &a}
		}
		name += c.getText(from, to)
	}
	if im.cfg.clean {
		// clean: keep the numeric attributes that change the column's range
		for {
			switch {
			case c.accept("UNSIGNED"):
				name += " unsigned"
			case c.accept("SIGNED"):
			case c.accept("ZEROFILL"):
				name += " zerofill"
			default:
				return model.ColumnType{TypeName: name}
			}
		}
	}
	return model.ColumnType{TypeName: name}
}

func (im *mysqlImporter) columnConstraint(c *cursor, f *model.Column, refs *[]*model.Ref) {
	switch {
	case c.peek().is("NOT") && c.at(1).is("NULL"):
		c.i += 2
		f.NotNull = boolp(true)
	case c.peek().is("NULL"):
		c.next()
		f.NotNull = boolp(false)
	case c.peek().is("SERIAL") && c.at(1).is("DEFAULT") && c.at(2).is("VALUE"):
		c.i += 3
	case c.peek().is("DEFAULT"):
		c.next()
		f.DBDefault = im.defaultValue(c)
	case c.accept("AUTO_INCREMENT"):
		f.Increment = boolp(true)
	case c.peek().is("ON") && c.at(1).is("UPDATE"):
		c.i += 2
		im.currentTimestamp(c)
	case c.peek().is("PRIMARY"):
		c.next()
		c.accept("KEY")
		f.PK = boolp(true)
	case c.peek().is("KEY"):
		c.next()
	case c.peek().is("UNIQUE"):
		c.next()
		c.accept("KEY")
		f.Unique = boolp(true)
	case c.peek().is("COMMENT"):
		c.next()
		if s := c.peek(); s.kind == tString {
			c.next()
			f.Note = &model.NoteValue{Value: im.str(s.text)}
		}
	case c.peek().is("REFERENCES"):
		c.next()
		if r := im.referenceDefinition(c); r != nil {
			*refs = append(*refs, r)
		}
	case c.peek().is("COLLATE") || c.peek().is("CHARSET") || c.peek().is("COLUMN_FORMAT") || c.peek().is("STORAGE") || c.peek().is("SRID"):
		c.next()
		c.acceptPunct("=")
		c.next()
	case c.peek().is("CHARACTER") && c.at(1).is("SET"):
		c.i += 2
		c.acceptPunct("=")
		c.next()
	case c.peek().is("GENERATED") || c.peek().is("AS"):
		c.accept("GENERATED")
		c.accept("ALWAYS")
		c.accept("AS")
		if c.peek().isPunct("(") {
			from, to := c.skipGroup()
			if im.cfg.clean {
				stored := c.accept("STORED") || c.accept("PERSISTENT")
				if !stored {
					c.accept("VIRTUAL")
				}
				setGenerated(f, c.origText(from+1, to-1), stored)
			}
		}
	case c.peek().is("CONSTRAINT") || c.peek().is("CHECK"):
		var name *string
		if c.accept("CONSTRAINT") {
			if n := c.peek(); !n.is("CHECK") {
				c.next()
				name = strp(im.uid(n))
			}
		}
		if c.accept("CHECK") {
			if ch := im.check(c, name); ch != nil {
				f.Checks = append(f.Checks, ch)
			}
		}
	case c.peek().isPunct("("):
		c.skipGroup()
	case c.peek().kind == tWord || c.peek().kind == tQIdent:
		c.next() // other column attributes (UNSIGNED, VISIBLE, STORED, ...)
	default:
		// a literal or operator here is a syntax error in the reference grammar
		if !im.cfg.clean {
			im.err = errUnsupported("unexpected " + c.peek().text + " in column definition")
		}
		c.next()
	}
}

// currentTimestamp: (CURRENT_TIMESTAMP | LOCALTIME | LOCALTIMESTAMP) ('(' n? ')')? | NOW '(' n? ')'
func (im *mysqlImporter) currentTimestamp(c *cursor) (string, bool) {
	t := c.peek()
	if !(t.is("CURRENT_TIMESTAMP") || t.is("LOCALTIME") || t.is("LOCALTIMESTAMP") || t.is("NOW")) {
		return "", false
	}
	from := c.i
	c.next()
	if c.peek().isPunct("(") {
		c.skipGroup()
	}
	return c.getText(from, c.i), true
}

func (im *mysqlImporter) defaultValue(c *cursor) *model.Default {
	t := c.peek()
	switch {
	case t.is("NULL"):
		c.next()
		return &model.Default{Type: "boolean", Value: strings.ToLower(t.text)}
	case t.is("CAST") && c.at(1).isPunct("("):
		c.next()
		from, to := c.skipGroup()
		inner := c.inner(from, to)
		// the expression before the top-level AS
		depth, end := 0, len(inner)
		for i, x := range inner {
			if x.isPunct("(") {
				depth++
			} else if x.isPunct(")") {
				depth--
			} else if depth == 0 && x.is("AS") {
				end = i
				break
			}
		}
		return &model.Default{Type: "expression", Value: im.expr(c, from+1, from+1+end)}
	case t.isPunct("-") || t.isPunct("+") || t.isPunct("~") || t.isPunct("!"):
		from := c.i
		c.next()
		c.next()
		return &model.Default{Type: "expression", Value: im.expr(c, from, c.i)}
	case t.kind == tString:
		v, typ := im.stringConst(c)
		return &model.Default{Type: typ, Value: v}
	case t.kind == tNumber:
		c.next()
		if im.cfg.clean && strings.HasPrefix(strings.ToLower(t.text), "0x") {
			return &model.Default{Type: "expression", Value: t.text}
		}
		if strings.HasPrefix(strings.ToLower(t.text), "0x") || strings.ContainsAny(t.text, ".eE") || isAllDigits(t.text) {
			return &model.Default{Type: "number", Value: t.text}
		}
		return &model.Default{Type: "number", Value: t.text}
	case t.is("TRUE") || t.is("FALSE"):
		c.next()
		return &model.Default{Type: "boolean", Value: t.text}
	case t.isPunct("("):
		from, to := c.skipGroup()
		return &model.Default{Type: "expression", Value: im.expr(c, from+1, to-1)}
	}
	if v, ok := im.currentTimestamp(c); ok {
		if c.peek().is("ON") && c.at(1).is("UPDATE") {
			c.i += 2
			im.currentTimestamp(c)
		}
		return &model.Default{Type: "expression", Value: v}
	}
	// anything else (b'01', x'..', identifiers, function calls): its text
	from := c.i
	c.next()
	if c.peek().isPunct("(") {
		c.skipGroup()
	}
	typ := "expression"
	if t.kind == tString {
		typ = "string"
	}
	if typ == "expression" {
		return &model.Default{Type: typ, Value: im.expr(c, from, c.i)}
	}
	return &model.Default{Type: typ, Value: c.getText(from, c.i)}
}

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if !isDigit(s[i]) {
			return false
		}
	}
	return s != ""
}

// ---------------------------------------------------------------- ALTER / INDEX / INSERT

func (im *mysqlImporter) alterTable(c *cursor) {
	schema, name := lastTwo(im.fullID(c))
	t := im.raw.findTable(schema, name)
	if t == nil {
		return
	}
	for !c.eof() {
		start := c.i
		if im.cfg.clean && im.alterSpecClean(c, t, schema, name) {
			c.skipUntil(func(t token) bool { return t.isPunct(",") })
			if !c.acceptPunct(",") {
				break
			}
			continue
		}
		if c.accept("ADD") {
			var cname *string
			if c.accept("CONSTRAINT") {
				if n := c.peek(); !(n.is("PRIMARY") || n.is("FOREIGN") || n.is("CHECK") || n.is("UNIQUE")) {
					c.next()
					cname = strp(im.uid(n))
				}
			}
			switch {
			case c.peek().is("PRIMARY") && c.at(1).is("KEY"):
				c.i += 2
				ix := im.indexTail(c, cname, true)
				ix.PK = boolp(true)
				if len(ix.Columns) > 1 {
					t.Indexes = append(t.Indexes, ix)
				} else if len(ix.Columns) == 1 {
					if f := findField(t, ix.Columns[0].Value); f != nil {
						f.PK = boolp(true)
					}
				}
			case c.peek().is("FOREIGN") && c.at(1).is("KEY"):
				c.i += 2
				if r := im.foreignKey(c, cname); r != nil {
					r.Endpoints[0].SchemaName, r.Endpoints[0].TableName = schema, name
					im.raw.refs = append(im.raw.refs, r)
				}
			case c.peek().is("CHECK"):
				c.next()
				if ch := im.check(c, cname); ch != nil {
					t.Checks = append(t.Checks, ch)
				}
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

func (im *mysqlImporter) createIndex(c *cursor, unique bool) {
	nameTok := c.next()
	if nameTok.text == "``" && !im.cfg.clean {
		im.err = errUnsupported("empty index name")
		return
	}
	var typ *string
	if c.accept("USING") {
		typ = strp(c.next().text)
	}
	if !c.accept("ON") {
		return
	}
	schema, tname := lastTwo(im.fullID(c))
	t := im.raw.findTable(schema, tname)
	if t == nil {
		return
	}
	cols := im.indexColumnNames(c)
	if ot := im.indexOptions(c); ot != nil {
		typ = ot
	}
	t.Indexes = append(t.Indexes, &model.Index{Name: strp(im.uid(nameTok)), Columns: cols, Unique: boolp(unique), Type: typ})
	im.lastIndexed = t
}

func (im *mysqlImporter) insert(c *cursor) {
	for c.accept("LOW_PRIORITY") || c.accept("DELAYED") || c.accept("HIGH_PRIORITY") || c.accept("IGNORE") {
	}
	c.accept("INTO")
	names := im.fullID(c)
	schema, table := lastTwo(names)
	if c.accept("PARTITION") && c.peek().isPunct("(") {
		c.skipGroup()
	}
	var cols []string
	explicit := false
	if c.peek().isPunct("(") && !c.at(1).is("SELECT") {
		explicit = true
		c.next()
		for !c.eof() && !c.peek().isPunct(")") {
			cols = append(cols, im.fullID(c)...)
			if !c.acceptPunct(",") {
				break
			}
		}
		c.acceptPunct(")")
	}
	if !explicit {
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
	if c.peek().is("SET") {
		// the reference crashes on INSERT ... SET (no insertStatementValue)
		if !im.cfg.clean {
			im.err = errUnsupported("INSERT ... SET")
		}
		return
	}
	if !(c.accept("VALUES") || c.accept("VALUE")) {
		// INSERT ... SELECT: a records block without rows
		if c.peek().is("SELECT") || c.peek().is("WITH") || c.peek().is("TABLE") || (c.peek().isPunct("(") && c.at(1).is("SELECT")) {
			im.raw.records = append(im.raw.records, rec)
		}
		return
	}
	if !c.peek().isPunct("(") {
		if !im.cfg.clean {
			im.err = errUnsupported("VALUES without a row")
		}
		return
	}
	for c.peek().isPunct("(") {
		c.next()
		var row []*model.RecordValue
		for !c.eof() && !c.peek().isPunct(")") {
			if c.peek().isPunct(",") {
				if !im.cfg.clean {
					im.err = errUnsupported("empty value in VALUES")
				}
				return
			}
			row = append(row, im.recordValues(c)...)
			c.skipUntil(boundary)
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
	im.raw.records = append(im.raw.records, rec)
}

// recordValues mirrors what flattenDepth(expressionOrDefault.accept(), 3)
// yields for one INSERT value: one value object for a constant or a
// function call; for other expressions the visitor's nested child arrays,
// which the DBML exporter renders as null (a column reference is one such
// entry, a binary or parenthesised expression three).
func (im *mysqlImporter) recordValues(c *cursor) []*model.RecordValue {
	from := c.i
	_, to := c.skipUntil(boundary)
	c.i = from
	n := to - from
	placeholders := func(k int) []*model.RecordValue {
		out := make([]*model.RecordValue, k)
		for i := range out {
			out[i] = &model.RecordValue{Type: "expression"}
		}
		c.i = to
		return out
	}
	if n > 1 && c.toks[from].isPunct("(") {
		if _, end := c.skipGroup(); end == to {
			return placeholders(3)
		}
		c.i = from
	}
	// fullColumnName: uid ('.' uid)*
	isName := n > 0
	for i := from; i < to; i++ {
		t := c.toks[i]
		if (i-from)%2 == 0 {
			if !(t.kind == tWord || t.kind == tQIdent) || t.is("NULL") || t.is("TRUE") || t.is("FALSE") || t.is("DEFAULT") {
				isName = false
			}
		} else if !t.isPunct(".") {
			isName = false
		}
	}
	if isName && n%2 == 1 {
		return placeholders(1)
	}
	// a top-level binary operator makes a 3-child expression atom
	depth := 0
	for i := from + 1; i < to; i++ {
		t := c.toks[i]
		if t.isPunct("(") {
			depth++
		} else if t.isPunct(")") {
			depth--
		} else if depth == 0 && t.kind == tPunct && binaryOps[t.text] || depth == 0 && (t.is("DIV") || t.is("MOD")) {
			return placeholders(3)
		}
	}
	return []*model.RecordValue{im.recordValue(c)}
}

var binaryOps = map[string]bool{"+": true, "-": true, "*": true, "/": true, "%": true, "&": true, "|": true, "^": true, "<<": true, ">>": true}

func (im *mysqlImporter) recordValue(c *cursor) *model.RecordValue {
	t := c.peek()
	from := c.i
	switch {
	case t.is("NULL"):
		c.next()
		return &model.RecordValue{Value: strings.ToLower(t.text), Type: "boolean"}
	case t.is("TRUE") || t.is("FALSE"):
		c.next()
		return &model.RecordValue{Value: t.text, Type: "boolean"}
	case t.kind == tString:
		v, typ := im.stringConst(c)
		return &model.RecordValue{Value: v, Type: typ}
	case t.kind == tNumber && boundary(c.at(1)):
		c.next()
		return &model.RecordValue{Value: t.text, Type: "number"}
	}
	_, to := c.skipUntil(boundary)
	return &model.RecordValue{Value: im.expr(c, from, to), Type: "expression"}
}

// mysqlStringConst mirrors visitConstant for string-like tokens: a bit
// string (b'01') keeps its text as a string, a hex string (X'0F') its text
// as a number, anything else (incl. N'..', _charset'..' and adjacent
// literals) is the getText with its first and last character dropped.
func (im *mysqlImporter) stringConst(c *cursor) (string, string) {
	if im.cfg.clean {
		t := c.peek()
		if t.text != "" && strings.ContainsRune("bBxX", rune(t.text[0])) && strings.HasPrefix(t.text[1:], "'") {
			c.next()
			return t.text, "expression"
		}
		var b strings.Builder
		for c.peek().kind == tString {
			b.WriteString(unescapeMySQLString(c.next().text))
		}
		return b.String(), "string"
	}
	t := c.peek()
	if t.text != "" && (t.text[0] == 'b' || t.text[0] == 'B') && strings.HasPrefix(t.text[1:], "'") {
		c.next()
		return t.text, "string"
	}
	if t.text != "" && (t.text[0] == 'x' || t.text[0] == 'X') && strings.HasPrefix(t.text[1:], "'") {
		c.next()
		return t.text, "number"
	}
	from := c.i
	for c.peek().kind == tString {
		c.next()
	}
	return stripQuotes(c.getText(from, c.i)), "string"
}

// mysqlDataTypes are the type keywords of MySqlParser.g4's dataType rule.
var mysqlDataTypes = setOf("CHAR", "CHARACTER", "VARCHAR", "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT", "NCHAR", "NVARCHAR",
	"LONG", "NATIONAL", "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT", "MIDDLEINT", "INT1", "INT2", "INT3",
	"INT4", "INT8", "REAL", "DOUBLE", "DECIMAL", "DEC", "FIXED", "NUMERIC", "FLOAT", "FLOAT4", "FLOAT8", "DATE", "TINYBLOB",
	"MEDIUMBLOB", "LONGBLOB", "BOOL", "BOOLEAN", "SERIAL", "BIT", "TIME", "TIMESTAMP", "DATETIME", "BINARY", "VARBINARY",
	"BLOB", "YEAR", "ENUM", "SET", "GEOMETRYCOLLECTION", "GEOMCOLLECTION", "LINESTRING", "MULTILINESTRING", "MULTIPOINT",
	"MULTIPOLYGON", "POINT", "POLYGON", "JSON", "GEOMETRY")

func setOf(xs ...string) map[string]bool {
	m := map[string]bool{}
	for _, x := range xs {
		m[x] = true
	}
	return m
}

// numericDims reports whether tokens are `n` or `n, n` (decimalLiteral).
func numericDims(ts []token) bool {
	if len(ts) != 1 && len(ts) != 3 {
		return false
	}
	if ts[0].kind != tNumber {
		return false
	}
	return len(ts) == 1 || (ts[1].isPunct(",") && ts[2].kind == tNumber)
}

// str decodes a string literal token (clean) or slices its quotes off like
// the reference (parity).
func (im *mysqlImporter) str(text string) string {
	if im.cfg.clean {
		return unescapeMySQLString(text)
	}
	return stripQuotes(text)
}

// expr renders expression tokens [from, to): getText (no whitespace) in
// parity mode, the source text with quoted identifiers rewritten in clean
// mode.
func (im *mysqlImporter) expr(c *cursor, from, to int) string {
	if im.cfg.clean {
		return c.origText(from, to)
	}
	return c.getText(from, to)
}

// alterSpecClean applies the ALTER TABLE specifications the parity port
// ignores (clean mode): ADD [COLUMN] / ADD (...) / ADD INDEX|KEY|UNIQUE|
// FULLTEXT|SPATIAL, MODIFY, CHANGE, DROP [COLUMN] / INDEX / KEY / FOREIGN
// KEY, RENAME COLUMN. It reports whether it consumed the specification.
func (im *mysqlImporter) alterSpecClean(c *cursor, t *model.Table, schema *string, name string) bool {
	def := func() {
		var inline, tableRefs []*model.Ref
		var pk *model.Index
		im.createDefinition(c, t, &inline, &tableRefs, &pk)
		for _, r := range tableRefs {
			r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
		}
		im.raw.refs = append(append(im.raw.refs, inline...), tableRefs...)
		if pk != nil && len(pk.Columns) > 0 {
			if f := findField(t, pk.Columns[0].Value); f != nil {
				f.PK = boolp(true)
			}
		}
		for _, f := range t.Fields {
			if strings.EqualFold(f.Type.TypeName, "enum") && f.Type.Args != nil {
				im.inlineEnum(t, f)
			}
		}
	}
	// position clauses after a column definition
	place := func(f *model.Column) {
		for i, x := range t.Fields {
			if x == f {
				t.Fields = append(t.Fields[:i], t.Fields[i+1:]...)
				break
			}
		}
		switch {
		case c.accept("FIRST"):
			t.Fields = append([]*model.Column{f}, t.Fields...)
		case c.accept("AFTER"):
			after := im.uid(c.next())
			for i, x := range t.Fields {
				if x.Name == after {
					t.Fields = append(t.Fields[:i+1], append([]*model.Column{f}, t.Fields[i+1:]...)...)
					return
				}
			}
			t.Fields = append(t.Fields, f)
		default:
			t.Fields = append(t.Fields, f)
		}
	}
	rename := func(old, nn string) {
		if old == nn {
			return
		}
		for _, ix := range t.Indexes {
			for _, ic := range ix.Columns {
				if ic.Type == "column" && ic.Value == old {
					ic.Value = nn
				}
			}
		}
		re := regexp.MustCompile("(^|[^A-Za-z0-9_$`])`?" + regexp.QuoteMeta(old) + "`?([^A-Za-z0-9_$`]|$)")
		fix := func(e string) string {
			return re.ReplaceAllString(e, "${1}`"+strings.ReplaceAll(nn, "`", "``")+"`${2}")
		}
		for _, ch := range t.Checks {
			ch.Expression = fix(ch.Expression)
		}
		for _, f := range t.Fields {
			for _, ch := range f.Checks {
				ch.Expression = fix(ch.Expression)
			}
			if g := f.Metadata[model.MetaGenerated]; g != "" {
				f.Metadata[model.MetaGenerated] = fix(g)
			}
		}
		for _, r := range im.raw.refs {
			for k := range r.Endpoints {
				ep := &r.Endpoints[k]
				if ep.TableName == name && str(ep.SchemaName) == str(schema) {
					for j, fn := range ep.FieldNames {
						if fn == old {
							ep.FieldNames[j] = nn
						}
					}
				}
			}
		}
	}
	replace := func(old string) {
		i := -1
		for j, f := range t.Fields {
			if f.Name == old {
				i = j
			}
		}
		n := len(t.Fields)
		def()
		if len(t.Fields) == n+1 {
			f := t.Fields[n]
			if i >= 0 {
				t.Fields[i] = f
				t.Fields = t.Fields[:n]
				rename(old, f.Name)
			}
			if c.peek().is("FIRST") || c.peek().is("AFTER") {
				place(f)
			}
		}
	}
	switch {
	case c.accept("ADD"):
		if c.peek().is("CONSTRAINT") || c.peek().is("PRIMARY") || c.peek().is("FOREIGN") || c.peek().is("CHECK") {
			c.i-- // the parity path handles these (and names refs)
			return false
		}
		if c.acceptPunct("(") {
			for !c.eof() && !c.peek().isPunct(")") {
				def()
				c.skipUntil(boundary)
				if !c.acceptPunct(",") {
					break
				}
			}
			c.acceptPunct(")")
			return true
		}
		if c.accept("COLUMN") || !(c.peek().is("INDEX") || c.peek().is("KEY") || c.peek().is("UNIQUE") || c.peek().is("FULLTEXT") || c.peek().is("SPATIAL")) {
			n := len(t.Fields)
			def()
			if len(t.Fields) == n+1 && (c.peek().is("FIRST") || c.peek().is("AFTER")) {
				place(t.Fields[n])
			}
			return true
		}
		if c.accept("UNIQUE") {
			if !c.accept("INDEX") {
				c.accept("KEY")
			}
			ix := im.indexTail(c, nil, true)
			ix.Unique = boolp(true)
			t.Indexes = append(t.Indexes, ix)
			return true
		}
		def() // INDEX / KEY / FULLTEXT / SPATIAL
		return true
	case c.accept("MODIFY"):
		c.accept("COLUMN")
		replace(im.uid(c.peek()))
		return true
	case c.accept("CHANGE"):
		c.accept("COLUMN")
		replace(im.uid(c.next()))
		return true
	case c.peek().is("RENAME") && c.at(1).is("COLUMN"):
		c.i += 2
		old := im.uid(c.next())
		c.accept("TO")
		nn := im.uid(c.next())
		for _, f := range t.Fields {
			if f.Name == old {
				f.Name = nn
			}
		}
		rename(old, nn)
		return true
	case c.accept("DROP"):
		switch {
		case c.accept("INDEX") || c.accept("KEY"):
			n := im.uid(c.next())
			ixs := t.Indexes[:0]
			for _, ix := range t.Indexes {
				if ix.Name == nil || *ix.Name != n {
					ixs = append(ixs, ix)
				}
			}
			t.Indexes = ixs
		case c.peek().is("FOREIGN") && c.at(1).is("KEY"):
			c.i += 2
			n := im.uid(c.next())
			refs := im.raw.refs[:0]
			for _, r := range im.raw.refs {
				if r.Name == nil || *r.Name != n {
					refs = append(refs, r)
				}
			}
			im.raw.refs = refs
		case c.peek().is("PRIMARY"):
			for _, f := range t.Fields {
				f.PK = nil
			}
			ixs := t.Indexes[:0]
			for _, ix := range t.Indexes {
				if !truthy(ix.PK) {
					ixs = append(ixs, ix)
				}
			}
			t.Indexes = ixs
			return true
		case c.peek().is("CONSTRAINT") || c.peek().is("CHECK"):
			return true
		default:
			c.accept("COLUMN")
			n := im.uid(c.next())
			for i, f := range t.Fields {
				if f.Name == n {
					t.Fields = append(t.Fields[:i], t.Fields[i+1:]...)
					break
				}
			}
		}
		return true
	}
	return false
}
