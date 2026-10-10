// Portions of this file are derived from holistics/dbml
// packages/dbml-core/src/parse/ANTLR/ASTGeneration/mssql/MssqlASTGen.js
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: re-implemented in Go over a hand-written tokenizer and
// recursive-descent statement parser (no ANTLR grammar); the visitor's
// semantics (what is kept, how text is rendered) are mirrored.

package importer

import (
	"strings"

	"github.com/aidenappl/go-dbml/model"
)

type mssqlImporter struct {
	raw rawDB
	err error
	cfg *config
}

func importMSSQL(src string, cfg *config) (*rawDB, error) {
	im := &mssqlImporter{cfg: cfg}
	forEachStatement(src, dMSSQL, false, cfg.skipStatement, cfg.maxStatementTokens(), func(st statement) {
		if im.err != nil {
			return
		}
		im.batch(newCursor(src, st.toks))
	})
	if im.err != nil {
		return nil, im.err
	}
	return &im.raw, nil
}

// batch walks the statements of one GO batch (T-SQL needs no separators).
func (im *mssqlImporter) batch(c *cursor) {
	// CREATE/ALTER PROCEDURE|FUNCTION|TRIGGER|VIEW must be alone in a batch.
	if c.peek().is("CREATE") || c.peek().is("ALTER") {
		i := 1
		if c.at(1).is("OR") && c.at(2).is("ALTER") {
			i = 3
		}
		switch strings.ToUpper(c.at(i).text) {
		case "PROCEDURE", "PROC", "FUNCTION", "TRIGGER", "VIEW":
			return
		}
	}
	depth := 0
	for !c.eof() {
		t := c.peek()
		if t.isPunct("(") {
			depth++
		} else if t.isPunct(")") && depth > 0 {
			depth--
		}
		if depth > 0 {
			c.next()
			continue
		}
		start := c.i
		switch {
		case t.is("CREATE") && c.at(1).is("TABLE"):
			c.i += 2
			im.createTable(c)
		case t.is("CREATE") && (c.at(1).is("INDEX") || c.at(1).is("UNIQUE") || c.at(1).is("CLUSTERED") || c.at(1).is("NONCLUSTERED")):
			c.next()
			im.createIndex(c)
		case t.is("ALTER") && c.at(1).is("TABLE"):
			c.i += 2
			im.alterTable(c)
		case t.is("INSERT"):
			c.next()
			im.insert(c)
		case t.is("EXEC") || t.is("EXECUTE"):
			c.next()
			im.execute(c)
		}
		if im.err != nil {
			return
		}
		if c.i == start {
			c.next()
		}
	}
}

// mssqlID mirrors visitId_: [x] and "x" lose their delimiters.
func mssqlID(t token) string {
	if t.kind == tQIdent {
		// [EXIST], [MODIFY], [QUERY] and [VALUE] are keyword tokens in the
		// reference lexer, so their brackets are kept
		switch strings.ToUpper(t.text) {
		case "[EXIST]", "[MODIFY]", "[QUERY]", "[VALUE]":
			return t.text
		}
		return stripQuotes(t.text)
	}
	return t.text
}

// tableName: id_ ('.' id_)*  (empty parts allowed: db..table)
func (im *mssqlImporter) tableName(c *cursor) []string {
	var names []string
	for {
		t := c.peek()
		if t.isPunct(".") {
			c.next()
			continue
		}
		if t.kind != tWord && t.kind != tQIdent {
			break
		}
		c.next()
		names = append(names, im.id(t))
		if !c.peek().isPunct(".") {
			break
		}
	}
	return names
}

// strictFind mirrors `t.name === tableName && t.schemaName === schemaName`.
func (im *mssqlImporter) strictFind(schema *string, name string) *model.Table {
	return im.raw.strict(schema, name)
}

// find is strictFind in parity mode (the reference compares schemaName
// with ===) and dbo-defaulting in clean mode.
func (im *mssqlImporter) find(schema *string, name string) *model.Table {
	if im.cfg.clean {
		return im.dboFind(schema, name)
	}
	return im.strictFind(schema, name)
}

// dboFind mirrors findTable (schema defaults to dbo).
func (im *mssqlImporter) dboFind(schema *string, name string) *model.Table {
	return im.raw.lookup("dbo", schema, name)
}

type msDefs struct {
	fields   []*model.Column
	inline   []*model.Ref
	indexes  []*model.Index
	refs     []*model.Ref
	defaults []struct {
		column string
		value  *model.Default
	}
	checks []*model.Check
}

func (im *mssqlImporter) createTable(c *cursor) {
	schema, name := lastTwo(im.tableName(c))
	if name == "" || !c.acceptPunct("(") {
		return
	}
	defs, tableIndices := im.definitions(c)
	c.acceptPunct(")")
	for _, r := range defs.refs {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
	}
	for _, r := range defs.inline {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
	}
	im.raw.refs = append(im.raw.refs, defs.inline...)
	im.raw.refs = append(im.raw.refs, defs.refs...)
	t := &model.Table{Name: name, SchemaName: schema, Fields: defs.fields, Checks: defs.checks}
	t.Indexes = append(append(t.Indexes, tableIndices...), defs.indexes...)
	im.raw.addTable(t)
	im.skipTableTail(c)
}

func (im *mssqlImporter) skipTableTail(c *cursor) {
	for c.peek().is("ON") || c.peek().is("TEXTIMAGE_ON") || c.peek().is("WITH") || c.peek().is("LOCK") || c.peek().is("FILESTREAM_ON") {
		c.next()
		if c.peek().isPunct("(") {
			c.skipGroup()
			continue
		}
		c.next()
		if c.peek().isPunct("(") {
			c.skipGroup()
		}
	}
}

// definitions parses column_def_table_constraints and table_indices up to
// the closing paren (elements separated by optional commas).
func (im *mssqlImporter) definitions(c *cursor) (*msDefs, []*model.Index) {
	defs := &msDefs{}
	var tableIndices []*model.Index
	for !c.eof() && !c.peek().isPunct(")") {
		if c.acceptPunct(",") {
			continue
		}
		start := c.i
		p := c.peek()
		switch {
		case p.is("INDEX"):
			c.next()
			if ix := im.tableIndex(c); ix != nil {
				tableIndices = append(tableIndices, ix)
			}
		case p.is("CONSTRAINT") || p.is("PRIMARY") || p.is("UNIQUE") || p.is("FOREIGN") || p.is("CHECK") || p.is("CONNECTION") || (p.is("DEFAULT") && !c.at(1).isPunct(")")):
			im.tableConstraint(c, defs)
		default:
			im.columnDefinition(c, defs)
		}
		if c.i == start {
			c.next()
		}
	}
	return defs, tableIndices
}

// tableIndex: INDEX id_ UNIQUE? clustered? '(' cols ')' ...
func (im *mssqlImporter) tableIndex(c *cursor) *model.Index {
	name := im.id(c.next())
	unique := c.accept("UNIQUE")
	c.accept("CLUSTERED")
	c.accept("NONCLUSTERED")
	if !c.peek().isPunct("(") {
		if !im.cfg.clean {
			im.err = errUnsupported("COLUMNSTORE table index")
		}
		return nil
	}
	cols := im.columnsWithOrder(c)
	im.skipIndexTail(c)
	return &model.Index{Name: strp(name), Unique: boolp(unique), Columns: cols}
}

// msStop is a definition boundary or the start of the next statement
// (T-SQL statements need no separator).
func msStop(t token) bool {
	if boundary(t) || t.isPunct(";") {
		return true
	}
	if t.kind == tWord {
		switch strings.ToUpper(t.text) {
		case "ALTER", "CREATE", "INSERT", "EXEC", "EXECUTE", "GO", "DROP", "SET", "IF", "BEGIN", "DECLARE", "PRINT", "USE", "SELECT", "UPDATE", "DELETE":
			return true
		}
	}
	return false
}

func (im *mssqlImporter) skipIndexTail(c *cursor) {
	for !c.eof() && !msStop(c.peek()) {
		if c.peek().is("INDEX") || c.peek().is("CONSTRAINT") {
			return
		}
		if c.peek().isPunct("(") {
			c.skipGroup()
			continue
		}
		c.next()
	}
}

// columnsWithOrder: '(' id_ (ASC|DESC)? (',' ...)* ')'
func (im *mssqlImporter) columnsWithOrder(c *cursor) []*model.IndexColumn {
	var cols []*model.IndexColumn
	if !c.acceptPunct("(") {
		return cols
	}
	for !c.eof() && !c.peek().isPunct(")") {
		if t := c.peek(); t.kind != tWord && t.kind != tQIdent {
			if !im.cfg.clean {
				im.err = errUnsupported("index column " + t.text)
				return cols
			}
			// clean: an expression key part
			from, to := c.skipUntil(func(t token) bool { return t.isPunct(",") || t.is("ASC") || t.is("DESC") })
			if from == to {
				c.skipUntil(func(t token) bool { return false })
				break
			}
			e := strings.TrimSpace(c.origText(from, to))
			if strings.HasPrefix(e, "(") && strings.HasSuffix(e, ")") && c.matchOf(from) == to-1 {
				e = c.origText(from+1, to-1)
			}
			cols = append(cols, &model.IndexColumn{Type: "expression", Value: e})
			c.accept("ASC")
			c.accept("DESC")
			if !c.acceptPunct(",") {
				break
			}
			continue
		}
		cols = append(cols, &model.IndexColumn{Type: "column", Value: im.id(c.next())})
		c.accept("ASC")
		c.accept("DESC")
		if !c.acceptPunct(",") {
			break
		}
	}
	c.acceptPunct(")")
	return cols
}

func (im *mssqlImporter) nameList(c *cursor) []string {
	out := []string{}
	if !c.acceptPunct("(") {
		return out
	}
	for !c.eof() && !c.peek().isPunct(")") {
		out = append(out, im.id(c.next()))
		if !c.acceptPunct(",") {
			break
		}
	}
	c.acceptPunct(")")
	return out
}

// tableConstraint mirrors visitTable_constraint.
func (im *mssqlImporter) tableConstraint(c *cursor, defs *msDefs) {
	name := ""
	hasName := false
	if c.accept("CONSTRAINT") {
		name, hasName = im.id(c.next()), true
	}
	switch {
	case c.peek().is("PRIMARY") || c.peek().is("UNIQUE"):
		pk := c.accept("PRIMARY")
		if pk {
			c.accept("KEY")
		} else {
			c.next()
		}
		c.accept("CLUSTERED")
		c.accept("NONCLUSTERED")
		ix := &model.Index{Name: strp(name), Columns: im.columnsWithOrder(c)}
		if pk {
			ix.PK = boolp(true)
		} else {
			ix.Unique = boolp(true)
		}
		im.skipIndexTail(c)
		defs.indexes = append(defs.indexes, ix)
	case c.peek().is("FOREIGN"):
		c.next()
		c.accept("KEY")
		cols := im.nameList(c)
		if !c.accept("REFERENCES") {
			return
		}
		r := im.foreignKeyOptions(c)
		r.Name = strp(name)
		r.Endpoints[0] = model.RefEndpoint{FieldNames: cols, Relation: "*"}
		r.Endpoints[1].Relation = "1"
		defs.refs = append(defs.refs, r)
	case c.peek().is("DEFAULT"):
		c.next()
		from, to := c.skipUntil(func(t token) bool { return t.is("FOR") || msStop(t) })
		val := im.expression(c, from, to)
		if c.accept("FOR") {
			col := im.id(c.next())
			if !hasName {
				name = col
			}
			c.acceptSeq("WITH", "VALUES")
			defs.defaults = append(defs.defaults, struct {
				column string
				value  *model.Default
			}{col, val})
		}
	case c.peek().is("CHECK"):
		if ch := im.checkConstraint(c); ch != nil {
			if name != "" {
				ch.Name = strp(name)
			}
			defs.checks = append(defs.checks, ch)
		}
	case c.peek().is("CONNECTION"):
		c.next()
		if c.peek().isPunct("(") {
			c.skipGroup()
		}
	}
}

// checkConstraint: CHECK (NOT FOR REPLICATION)? '(' search_condition ')'
func (im *mssqlImporter) checkConstraint(c *cursor) *model.Check {
	c.accept("CHECK")
	c.acceptSeq("NOT", "FOR", "REPLICATION")
	if !c.peek().isPunct("(") {
		return nil
	}
	from, to := c.skipGroup()
	expr := c.origText(from+1, to-1)
	return &model.Check{Expression: expr}
}

// foreignKeyOptions: table_name '(' cols ')' (on_delete | on_update)* (NOT FOR REPLICATION)?
func (im *mssqlImporter) foreignKeyOptions(c *cursor) *model.Ref {
	schema, table := lastTwo(im.tableName(c))
	cols := im.nameList(c)
	r := &model.Ref{}
	for c.peek().is("ON") && (c.at(1).is("DELETE") || c.at(1).is("UPDATE")) {
		c.next()
		which := c.next()
		a := c.next()
		action := "SET DEFAULT"
		switch {
		case a.is("NO"):
			c.next()
			action = "NO ACTION"
		case a.is("CASCADE"):
			action = "CASCADE"
		case a.is("SET"):
			if c.next().is("NULL") {
				action = "SET NULL"
			}
		default:
			if !im.cfg.clean {
				im.err = errUnsupported("referential action " + a.text)
			}
			action = strings.ToUpper(a.text)
		}
		if which.is("DELETE") {
			if r.OnDelete == nil {
				r.OnDelete = strp(action)
			}
		} else if r.OnUpdate == nil {
			r.OnUpdate = strp(action)
		}
	}
	c.acceptSeq("NOT", "FOR", "REPLICATION")
	r.Endpoints[0] = model.RefEndpoint{Relation: "0..*"}
	r.Endpoints[1] = model.RefEndpoint{TableName: table, SchemaName: schema, FieldNames: cols, Relation: "0..1"}
	return r
}

var msElementStart = map[string]bool{"CONSTRAINT": true, "NOT": true, "NULL": true, "PRIMARY": true, "UNIQUE": true,
	"FOREIGN": true, "REFERENCES": true, "CHECK": true, "IDENTITY": true, "COLLATE": true, "SPARSE": true, "FILESTREAM": true,
	"ROWGUIDCOL": true, "MASKED": true, "GENERATED": true, "ENCRYPTED": true, "INDEX": true, "DEFAULT": true, "PERSISTED": true}

func msElementBoundary(t token) bool {
	return msStop(t) || (t.kind == tWord && msElementStart[strings.ToUpper(t.text)])
}

func (im *mssqlImporter) columnDefinition(c *cursor, defs *msDefs) {
	nameTok := c.next()
	f := &model.Column{Name: im.id(nameTok)}
	if c.peek().is("AS") && im.cfg.clean {
		// computed column: T-SQL gives it no type
		c.next()
		from, to := c.skipUntil(msElementBoundary)
		e := strings.TrimSpace(c.origText(from, to))
		if strings.HasPrefix(e, "(") && c.matchOf(from) == to-1 {
			e = c.origText(from+1, to-1)
		}
		f.Type = model.ColumnType{TypeName: "sql_variant"}
		setGenerated(f, e, c.accept("PERSISTED"))
	} else if c.accept("AS") {
		from, to := c.skipUntil(msElementBoundary)
		v := im.expression(c, from, to)
		typ := "AS " + jsStr(v)
		if c.peek().is("PERSISTED") {
			typ += " " + c.next().text
		}
		f.Type = model.ColumnType{TypeName: typ}
	} else {
		f.Type = model.ColumnType{TypeName: im.dataType(c)}
		// anything but a column element after the type (e.g. the VARYING of
		// "character varying", an array suffix) is a syntax error there
		if p := c.peek(); !msStop(p) && !(p.kind == tWord && msElementStart[strings.ToUpper(p.text)]) {
			if !im.cfg.clean {
				im.err = errUnsupported("unexpected " + p.text + " after the type of column " + f.Name)
				return
			}
			// clean: multi-word / array types ("character varying(10)", int[])
			from, to := c.skipUntil(msElementBoundary)
			f.Type.TypeName += " " + c.origText(from, to)
		}
	}
	var refs []*model.Ref
	for !c.eof() && !msStop(c.peek()) {
		start := c.i
		im.columnElement(c, f, &refs)
		if c.peek().is("INDEX") {
			// column_index: INDEX id_ (CLUSTERED|NONCLUSTERED)? ... - skipped
			c.next()
			im.skipIndexTail(c)
		}
		if c.i == start {
			c.next()
		}
	}
	for _, r := range refs {
		r.Endpoints[0].FieldNames = []string{f.Name}
	}
	defs.inline = append(defs.inline, refs...)
	defs.fields = append(defs.fields, f)
}

// dataType mirrors visitData_type.
func (im *mssqlImporter) dataType(c *cursor) string {
	t := c.next()
	if t.is("DOUBLE") {
		if p := c.peek(); p.is("PRECISION") {
			c.next()
			return t.text + " " + p.text
		}
		return t.text
	}
	id := im.id(t)
	if c.peek().isPunct("(") && c.at(1).is("MAX") && c.at(2).isPunct(")") {
		c.i += 3
		return id + "(MAX)"
	}
	if !im.cfg.clean && c.accept("IDENTITY") {
		if c.peek().isPunct("(") && c.at(1).kind == tNumber && c.at(2).isPunct(",") && c.at(3).kind == tNumber && c.at(4).isPunct(")") {
			s := id + " IDENTITY(" + c.at(1).text + "," + c.at(3).text + ")"
			c.i += 5
			return s
		}
		return id + " IDENTITY"
	}
	if c.peek().isPunct("(") && c.at(1).kind == tNumber {
		if c.at(2).isPunct(")") {
			s := id + "(" + c.at(1).text + ")"
			c.i += 3
			return s
		}
		if c.at(2).isPunct(",") && c.at(3).kind == tNumber && c.at(4).isPunct(")") {
			s := id + "(" + c.at(1).text + "," + c.at(3).text + ")"
			c.i += 5
			return s
		}
	}
	return id
}

// columnElement mirrors visitColumn_definition_element / visitColumn_constraint.
func (im *mssqlImporter) columnElement(c *cursor, f *model.Column, refs *[]*model.Ref) {
	var cname *string
	if c.accept("CONSTRAINT") {
		cname = strp(im.id(c.next()))
	}
	switch {
	case c.accept("DEFAULT"):
		from, to := c.skipUntil(msElementBoundary)
		f.DBDefault = im.cleanExpr(im.expression(c, from, to))
	case c.accept("IDENTITY"):
		if c.peek().isPunct("(") {
			c.skipGroup()
		}
		f.Increment = boolp(true)
	case c.peek().is("NOT") && c.at(1).is("NULL"):
		c.i += 2
		f.NotNull = boolp(true)
	case c.peek().is("NOT") && c.at(1).is("FOR"):
		c.acceptSeq("NOT", "FOR", "REPLICATION")
	case c.accept("NULL"):
		f.NotNull = boolp(false)
	case c.peek().is("PRIMARY") || c.peek().is("UNIQUE"):
		if c.accept("PRIMARY") {
			c.accept("KEY")
			f.PK = boolp(true)
		} else {
			c.next()
			f.Unique = boolp(true)
		}
		c.accept("CLUSTERED")
		c.accept("NONCLUSTERED")
		for c.peek().is("WITH") || c.peek().is("ON") {
			c.next()
			if c.peek().isPunct("(") {
				c.skipGroup()
			} else {
				c.next()
			}
		}
	case c.peek().is("FOREIGN") || c.peek().is("REFERENCES"):
		c.accept("FOREIGN")
		c.accept("KEY")
		if c.accept("REFERENCES") {
			*refs = append(*refs, im.foreignKeyOptions(c))
		}
	case c.peek().is("CHECK"):
		if ch := im.checkConstraint(c); ch != nil {
			ch.Name = cname
			f.Checks = append(f.Checks, ch)
		}
	case c.accept("COLLATE"):
		c.next()
	case c.peek().isPunct("("):
		c.skipGroup()
	default:
		c.next()
	}
}

// builtinFunctions are T-SQL built-ins callable without parentheses.
var builtinNoParens = map[string]bool{"CURRENT_TIMESTAMP": true, "CURRENT_USER": true, "SESSION_USER": true, "SYSTEM_USER": true, "USER": true, "CURRENT_DATE": true}

// maxExprLevels bounds the paren/unary peeling of one expression; deeper
// input is kept as raw expression text.
const maxExprLevels = 256

// expression mirrors visitExpression for tokens [from, to). It is
// iterative: '(' expr ')' and unary operators are peeled in a loop using
// the cursor's precomputed paren-match table (no recursion, O(n) for
// chains like "((((1))))" or "- - - 1").
func (im *mssqlImporter) expression(c *cursor, from, to int) *model.Default {
	for to > from && c.toks[to-1].isPunct(";") {
		to--
	}
	if to <= from {
		return nil
	}
	start, end := from, to
	prefix := ""
	maxBinary := -1 // highest top-level binary operator position of the current range
	scanned := -1   // range start maxBinary was computed for
	wrap := func(d *model.Default) *model.Default {
		if prefix == "" || d == nil {
			if d == nil && prefix != "" {
				return &model.Default{Type: "expression", Value: c.origText(start, end)}
			}
			return d
		}
		return &model.Default{Type: d.Type, Value: prefix + jsStr(d)}
	}
	for level := 0; ; level++ {
		if level > maxExprLevels {
			return &model.Default{Type: "expression", Value: c.origText(start, end)}
		}
		n := to - from
		if n <= 0 {
			return wrap(nil)
		}
		t := c.toks[from]
		if n == 1 {
			return wrap(im.atom(c, t, c.origText(from, to)))
		}
		// '(' expression ')'
		if t.isPunct("(") && c.matchOf(from) == to-1 {
			if c.toks[from+1].is("SELECT") {
				return wrap(nil)
			}
			from, to = from+1, to-1
			scanned = -1 // depth changed: recompute top-level operators
			continue
		}
		// unary operator: op + inner.value, inner.type
		if t.isPunct("-") || t.isPunct("+") || t.isPunct("~") {
			if scanned < 0 || scanned > from {
				maxBinary, scanned = topLevelBinaryMax(c, from, to), from
			}
			if maxBinary <= from+1 {
				prefix += t.text
				from++
				continue
			}
		}
		orig := c.origText(from, to)
		switch {
		case t.isPunct("$"): // money
			return wrap(&model.Default{Type: "string", Value: c.getText(from, to)})
		case (t.kind == tWord || t.kind == tQIdent) && c.toks[from+1].isPunct("(") && c.matchOf(from+1) == to-1:
			// function call covering everything
			return wrap(&model.Default{Type: "expression", Value: orig})
		case isDottedName(c, from, to):
			// full_column_name: a.b.c (rendered like a JS array: "a,b,c")
			var parts []string
			for i := from; i < to; i += 2 {
				parts = append(parts, im.id(c.toks[i]))
			}
			return wrap(&model.Default{Type: "expression", Value: strings.Join(parts, ",")})
		}
		return wrap(&model.Default{Type: "expression", Value: orig})
	}
}

// msAtom mirrors visitExpression for a single token.
func (im *mssqlImporter) atom(c *cursor, t token, orig string) *model.Default {
	switch {
	case t.is("NULL"):
		return &model.Default{Type: "boolean", Value: strings.ToLower(t.text)}
	case t.kind == tString:
		return im.strConst(t.text)
	case t.kind == tNumber:
		if strings.HasPrefix(strings.ToLower(t.text), "0x") {
			// BINARY: getStringFromRawString slices the first/last char
			return &model.Default{Type: "string", Value: stripQuotes(t.text)}
		}
		return &model.Default{Type: "number", Value: t.text}
	case t.kind == tWord && builtinNoParens[strings.ToUpper(t.text)]:
		return &model.Default{Type: "expression", Value: orig}
	case t.kind == tWord || t.kind == tQIdent:
		return &model.Default{Type: "expression", Value: im.id(t)}
	case t.kind == tVar:
		return &model.Default{Type: "expression", Value: t.text}
	}
	return &model.Default{Type: "expression", Value: orig}
}

// topLevelBinaryMax returns the highest position in (from, to) holding a
// binary operator at paren depth 0, or -1.
func topLevelBinaryMax(c *cursor, from, to int) int {
	max := -1
	depth := 0
	for i := from; i < to; i++ {
		t := c.toks[i]
		if t.isPunct("(") {
			depth++
		} else if t.isPunct(")") {
			depth--
		} else if depth == 0 && i > from && t.kind == tPunct && binaryOps[t.text] {
			max = i
		}
	}
	return max
}

func isDottedName(c *cursor, from, to int) bool {
	if (to-from)%2 == 0 {
		return false
	}
	for i := from; i < to; i++ {
		t := c.toks[i]
		if (i-from)%2 == 0 {
			if t.kind != tWord && t.kind != tQIdent {
				return false
			}
		} else if !t.isPunct(".") {
			return false
		}
	}
	return true
}

// msString mirrors getStringFromRawString: N'..' is kept verbatim as an
// expression, other strings lose their first and last character.
func msString(text string) *model.Default {
	if strings.HasPrefix(text, "N'") {
		return &model.Default{Type: "expression", Value: text}
	}
	return &model.Default{Type: "string", Value: stripQuotes(text)}
}

func jsStr(d *model.Default) string {
	if d == nil {
		return "undefined"
	}
	if s, ok := d.Value.(string); ok {
		return s
	}
	return ""
}

// ---------------------------------------------------------------- ALTER / INDEX / INSERT / EXEC

func (im *mssqlImporter) alterTable(c *cursor) {
	schema, name := lastTwo(im.tableName(c))
	// (WITH CHECK)? (NOCHECK|CHECK) CONSTRAINT ALL is not in the reference
	// grammar (constraint must be an identifier): a syntax error there.
	save := c.i
	c.acceptSeq("WITH", "CHECK")
	c.acceptSeq("WITH", "NOCHECK")
	if (c.accept("NOCHECK") || c.accept("CHECK")) && c.accept("CONSTRAINT") && c.peek().is("ALL") {
		if !im.cfg.clean {
			im.err = errUnsupported("CHECK CONSTRAINT ALL")
		}
		return
	}
	c.i = save
	t := im.find(schema, name)
	if t == nil {
		return
	}
	if c.peek().is("WITH") && (c.at(1).is("CHECK") || c.at(1).is("NOCHECK")) && c.at(2).is("ADD") {
		c.i += 3
		var cname *string
		if c.accept("CONSTRAINT") {
			cname = strp(im.id(c.next()))
		}
		if !(c.peek().is("FOREIGN") && c.at(1).is("KEY")) {
			return
		}
		c.i += 2
		local := im.nameList(c)
		if !c.accept("REFERENCES") {
			return
		}
		rs, rt := lastTwo(im.tableName(c))
		refCols := local
		if c.peek().isPunct("(") {
			refCols = im.nameList(c)
		}
		r := &model.Ref{Name: cname}
		for c.peek().is("ON") && (c.at(1).is("DELETE") || c.at(1).is("UPDATE")) {
			c.next()
			which := c.next()
			a := c.next()
			action := "SET DEFAULT"
			switch {
			case a.is("NO"):
				c.next()
				action = "NO ACTION"
			case a.is("CASCADE"):
				action = "CASCADE"
			case a.is("SET"):
				if c.next().is("NULL") {
					action = "SET NULL"
				}
			}
			if which.is("DELETE") {
				if r.OnDelete == nil {
					r.OnDelete = strp(action)
				}
			} else if r.OnUpdate == nil {
				r.OnUpdate = strp(action)
			}
		}
		r.Endpoints[0] = model.RefEndpoint{TableName: name, SchemaName: schema, FieldNames: local, Relation: "*"}
		r.Endpoints[1] = model.RefEndpoint{TableName: rt, SchemaName: rs, FieldNames: refCols, Relation: "1"}
		im.raw.refs = append(im.raw.refs, r)
		return
	}
	if !c.accept("ADD") {
		return
	}
	defs := &msDefs{}
	for !c.eof() && !c.peek().is("GO") {
		if c.acceptPunct(",") {
			continue
		}
		p := c.peek()
		if p.is("ALTER") || p.is("CREATE") || p.is("INSERT") || p.is("EXEC") || p.is("EXECUTE") || p.isPunct(";") {
			break
		}
		start := c.i
		switch {
		case p.is("CONSTRAINT") || p.is("PRIMARY") || p.is("UNIQUE") || p.is("FOREIGN") || p.is("CHECK") || p.is("DEFAULT"):
			im.tableConstraint(c, defs)
		default:
			im.columnDefinition(c, defs)
		}
		if c.i == start {
			c.next()
		}
	}
	for _, r := range defs.inline {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
	}
	for _, r := range defs.refs {
		r.Endpoints[0].TableName, r.Endpoints[0].SchemaName = name, schema
	}
	im.raw.refs = append(im.raw.refs, defs.inline...)
	im.raw.refs = append(im.raw.refs, defs.refs...)
	t.Fields = append(t.Fields, defs.fields...)
	t.Indexes = append(t.Indexes, defs.indexes...)
	for _, d := range defs.defaults {
		if f := findField(t, d.column); f != nil {
			f.DBDefault = d.value
		}
	}
	for _, ch := range defs.checks {
		if ch.Name != nil && *ch.Name == "" {
			ch.Name = nil
		}
		t.Checks = append(t.Checks, ch)
	}
}

func (im *mssqlImporter) createIndex(c *cursor) {
	unique := c.accept("UNIQUE")
	c.accept("CLUSTERED")
	c.accept("NONCLUSTERED")
	if !c.accept("INDEX") {
		return
	}
	name := im.id(c.next())
	if !c.accept("ON") {
		return
	}
	schema, table := lastTwo(im.tableName(c))
	t := im.find(schema, table)
	cols := im.columnsWithOrder(c)
	if t == nil {
		return
	}
	t.Indexes = append(t.Indexes, &model.Index{Name: strp(name), Unique: boolp(unique), Columns: cols})
}

func (im *mssqlImporter) insert(c *cursor) {
	if !im.cfg.records || im.cfg.rows >= im.cfg.maxRows {
		if im.cfg.records {
			im.cfg.truncated = true
		}
		return // the batch loop skips the rest of the statement
	}
	if c.accept("TOP") && c.peek().isPunct("(") {
		c.skipGroup()
		c.accept("PERCENT")
	}
	c.accept("INTO")
	schema, table := lastTwo(im.tableName(c))
	if c.peek().is("WITH") && c.at(1).isPunct("(") {
		c.next()
		c.skipGroup()
	}
	var cols []string
	explicit := false
	if c.peek().isPunct("(") && !c.at(1).is("SELECT") {
		explicit = true
		c.next()
		for !c.eof() && !c.peek().isPunct(")") {
			t := c.next()
			if t.isPunct(".") {
				continue
			}
			cols = append(cols, im.id(t))
			c.acceptPunct(",")
		}
		c.acceptPunct(")")
	}
	if !explicit {
		if t := im.dboFind(schema, table); t != nil {
			for _, f := range t.Fields {
				cols = append(cols, f.Name)
			}
		}
	}
	if cols == nil {
		cols = []string{}
	}
	if c.accept("OUTPUT") {
		c.skipUntil(func(t token) bool {
			return t.is("VALUES") || t.is("SELECT") || t.is("DEFAULT") || t.is("EXEC") || t.is("EXECUTE")
		})
	}
	rec := &model.TableRecord{SchemaName: schema, TableName: table, Columns: cols}
	if c.accept("VALUES") {
		for c.peek().isPunct("(") {
			c.next()
			row := []*model.RecordValue{}
			for !c.eof() && !c.peek().isPunct(")") {
				from, to := c.skipUntil(boundary)
				if d := im.expression(c, from, to); d != nil {
					row = append(row, &model.RecordValue{Value: d.Value, Type: d.Type})
				} else {
					row = append(row, nil)
				}
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

// execute handles EXEC sp_addextendedproperty (descriptions as notes).
func (im *mssqlImporter) execute(c *cursor) {
	if c.peek().kind == tVar && c.at(1).isPunct("=") {
		c.i += 2
	}
	names := im.tableName(c)
	if len(names) == 0 || names[len(names)-1] != "sp_addextendedproperty" {
		return
	}
	args := map[string]string{}
	positional := []string{"name", "value", "level0type", "level0name", "level1type", "level1name", "level2type", "level2name"}
	for pos := 0; !c.eof(); pos++ {
		v := c.peek()
		key := ""
		if v.kind == tVar && c.at(1).isPunct("=") {
			c.i += 2
			key = strings.ToLower(v.text[1:])
		} else {
			if v.kind != tString && v.kind != tNumber && v.kind != tWord && v.kind != tQIdent {
				return
			}
			if !im.cfg.clean {
				// positional arguments: the reference crashes on arg.name
				im.err = errUnsupported("positional sp_addextendedproperty arguments")
				return
			}
			if pos >= len(positional) {
				break
			}
			key = positional[pos]
		}
		from, to := c.skipUntil(func(t token) bool { return t.isPunct(",") || t.isPunct(";") || msStop(t) })
		args[key] = im.param(c, from, to)
		if !c.acceptPunct(",") {
			break
		}
	}
	if !strings.Contains(args["name"], "Description") {
		return
	}
	if !strings.Contains(strings.ToLower(args["level0type"]), "schema") {
		return
	}
	var schema *string
	if l0 := args["level0name"]; l0 != "dbo" {
		schema = strp(l0)
	}
	table := ""
	if strings.Contains(strings.ToLower(args["level1type"]), "table") {
		table = args["level1name"]
	}
	t := im.find(schema, table)
	if t == nil {
		return
	}
	if _, ok := args["level2type"]; !ok || args["level2type"] == "" {
		t.Note = &model.NoteValue{Value: args["value"]}
		return
	}
	col := ""
	if strings.Contains(strings.ToLower(args["level2type"]), "column") {
		col = args["level2name"]
	}
	if f := findField(t, col); f != nil {
		f.Note = &model.NoteValue{Value: args["value"]}
	}
}

// param mirrors visitExecute_parameter: constants per visitConstant,
// anything else its getText. Clean mode decodes strings and identifiers.
func (im *mssqlImporter) param(c *cursor, from, to int) string {
	if to-from == 1 {
		t := c.toks[from]
		if t.kind == tString {
			if im.cfg.clean {
				return unescapeStdString(t.text)
			}
			return jsStr(im.strConst(t.text))
		}
		if im.cfg.clean && t.kind == tQIdent {
			return unquoteIdent(t.text)
		}
	}
	return c.getText(from, to)
}

// id decodes an identifier token (mssqlID in parity mode, a plain unquote
// in clean mode).
func (im *mssqlImporter) id(t token) string {
	if im.cfg.clean {
		if t.kind == tQIdent {
			return unquoteIdent(t.text)
		}
		return t.text
	}
	return mssqlID(t)
}

// strConst decodes a string literal: msString in parity mode, the
// unescaped value (N'..' included) in clean mode.
func (im *mssqlImporter) strConst(text string) *model.Default {
	if im.cfg.clean {
		return &model.Default{Type: "string", Value: unescapeStdString(text)}
	}
	return msString(text)
}

// cleanExpr rewrites [x] identifiers in an expression default (clean mode).
func (im *mssqlImporter) cleanExpr(d *model.Default) *model.Default {
	if im.cfg.clean && d != nil && d.Type == "expression" {
		if s, ok := d.Value.(string); ok {
			d.Value = s
		}
	}
	return d
}
