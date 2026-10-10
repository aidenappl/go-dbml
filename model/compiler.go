// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/compiler/queries/{canonicalName,files/usableMembers}.ts, src/core/global_modules/{schema,table}/index.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
	"github.com/aidenappl/go-dbml/lexer"
	"github.com/aidenappl/go-dbml/parser"
)

// SourceLoader returns the source of a file by absolute path, or ok=false if
// the file does not exist. Paths are slash-separated and absolute ("/main.dbml").
type SourceLoader func(filepath string) (src string, ok bool)

// Result is the outcome of compiling one entry file.
type Result struct {
	// Database is nil when Diagnostics contains at least one error.
	Database *Database
	// Partial is a best-effort model built even when there are errors (for
	// editors: outline, diagram preview). It is nil only when go-dbml hit an
	// internal error. Treat it as untrusted: references may be missing.
	Partial     *Database
	Diagnostics diag.List
	// Program is the entry file's syntax tree (nil only on internal error).
	Program *ast.Program
}

// Compile parses, validates, binds and interprets a single in-memory file
// with DefaultOptions.
func Compile(src, filepath string) Result {
	return CompileWithOptions(src, filepath, DefaultOptions())
}

// CompileWithOptions is Compile with explicit limits.
func CompileWithOptions(src, filepath string, opts Options) Result {
	return CompileProjectWithOptions(filepath, func(p string) (string, bool) {
		if p == filepath {
			return src, true
		}
		return "", false
	}, opts)
}

// CompileProject compiles entry and every file it (transitively) imports via
// `use` / `reuse`, with DefaultOptions.
func CompileProject(entry string, load SourceLoader) Result {
	return CompileProjectWithOptions(entry, load, DefaultOptions())
}

// CompileProjectWithOptions is CompileProject with explicit limits. It never
// panics: an unexpected internal failure is returned as an INTERNAL_ERROR
// diagnostic.
func CompileProjectWithOptions(entry string, load SourceLoader, opts Options) (res Result) {
	opts = opts.withDefaults()
	defer func() {
		if r := recover(); r != nil {
			if r == errCanceled {
				res = Result{Diagnostics: diag.List{canceledDiag(entry)}}
				return
			}
			res = Result{Diagnostics: diag.List{diag.New(diag.INTERNAL_ERROR, fmt.Sprintf("internal error in go-dbml: %v", r), entry,
				diag.Range{Start: Position{Line: 1, Column: 1}, End: Position{Line: 1, Column: 1}})}}
		}
	}()
	c := &compiler{load: load, units: map[string]*unit{}, opts: opts}
	return c.run(entry)
}

type compiler struct {
	opts  Options
	load  SourceLoader
	units map[string]*unit
	order []*unit

	// memoisation (all results are pure functions of the parsed files)
	colCache       map[colCacheKey]*colSet
	reachableCache map[*unit][]*unit
	// limitDiags are resource-limit diagnostics raised while loading files.
	limitDiags []diag.Diagnostic
	polls      int
}

// errCanceled aborts a canceled compile (recovered in
// CompileProjectWithOptions).
var errCanceled = &struct{ s string }{"canceled"}

func canceledDiag(file string) diag.Diagnostic {
	return diag.New(diag.CANCELED, "Compilation was canceled", file, diag.Range{Start: Position{Line: 1, Column: 1}, End: Position{Line: 1, Column: 1}})
}

// poll aborts the compile when Options.Cancel is closed. It checks the
// channel every 64 calls (always when now is set).
func (c *compiler) poll(now bool) {
	if c.opts.Cancel == nil {
		return
	}
	c.polls++
	if !now && c.polls&63 != 0 {
		return
	}
	select {
	case <-c.opts.Cancel:
		panic(errCanceled)
	default:
	}
}

type colCacheKey struct {
	t  *sym
	sc *scope
}

// colSet is a memoised merged column list with a name index.
type colSet struct {
	list   []*col
	byName map[string]*col
}

type unit struct {
	path    string
	src     string
	prog    *ast.Program
	diags   []diag.Diagnostic // lex + parse
	decls   []*sym
	imports []*importDecl
	scope   *scope
	// declIndex maps declarations to their symbols (built lazily).
	declIndex map[*ast.ElementDeclaration]*sym
	// unresolved lists selective imports whose element was not found.
	unresolved []unresolvedImport
}

type importDecl struct {
	decl     *ast.UseDeclaration
	path     string // resolved absolute path ("" when not relative)
	reuse    bool
	wildcard bool
	specs    []*ast.UseSpecifier
	target   *unit
}

func (c *compiler) unit(p string) *unit {
	if u, ok := c.units[p]; ok {
		return u
	}
	if len(c.units) >= c.opts.MaxFiles {
		c.limitDiags = append(c.limitDiags, diag.New(diag.TOO_MANY_FILES, fmt.Sprintf("Project imports more than %d files; '%s' was not loaded", c.opts.MaxFiles, p), p, diag.Range{Start: Position{Line: 1, Column: 1}, End: Position{Line: 1, Column: 1}}))
		c.units[p] = nil
		return nil
	}
	src, ok := c.load(p)
	if !ok && len(c.units) > 0 {
		c.units[p] = nil
		return nil
	}
	u := &unit{path: p, src: src}
	c.units[p] = u
	var tooLarge bool
	if len(src) > c.opts.MaxBytes {
		tooLarge = true
		src = ""
	}
	res := parser.ParseWithOptions(src, p, parser.Options{MaxDepth: c.opts.MaxDepth, MaxErrors: c.opts.MaxDiagnostics, Cancel: c.opts.Cancel,
		MaxTokens: c.opts.tokenLimit(len(src)), MaxElements: c.opts.MaxElements, MaxRecordBytes: c.opts.MaxRecordBytes, MaxRecordRows: c.opts.MaxRecordRows})
	u.prog = res.Program
	u.diags = res.Errors
	if tooLarge {
		u.diags = append(u.diags, diag.New(diag.INPUT_TOO_LARGE, fmt.Sprintf("File is larger than the %d-byte limit and was not parsed", c.opts.MaxBytes), p, diag.Range{Start: Position{Line: 1, Column: 1}, End: Position{Line: 1, Column: 1}}))
	}
	c.order = append(c.order, u)
	c.collectDecls(u)
	for _, ud := range u.prog.Uses() {
		if ud.ImportPath == nil || ast.IsNil(ud.Specifiers) {
			continue
		}
		target := resolveImportPath(p, ud.ImportPath.Value)
		imp := &importDecl{decl: ud, reuse: ud.IsReuse()}
		if _, ok := ud.Specifiers.(*ast.Wildcard); ok {
			imp.wildcard = true
		} else if l, ok := ud.Specifiers.(*ast.UseSpecifierList); ok {
			imp.specs = l.Specifiers
		}
		imp.path = target
		if target != "" {
			imp.target = c.unit(target)
		}
		u.imports = append(u.imports, imp)
	}
	return u
}

func resolveImportPath(from, spec string) string {
	if !strings.HasPrefix(spec, "./") && !strings.HasPrefix(spec, "../") && !strings.HasPrefix(spec, "/") {
		return ""
	}
	var p string
	if strings.HasPrefix(spec, "/") {
		p = path.Clean(spec)
	} else {
		p = path.Join(path.Dir(from), spec)
	}
	if !strings.HasSuffix(p, ".dbml") {
		p += ".dbml"
	}
	return p
}

// reachable returns u and every unit reachable through imports.
func (c *compiler) reachable(u *unit) []*unit {
	if r, ok := c.reachableCache[u]; ok {
		return r
	}
	if c.reachableCache == nil {
		c.reachableCache = map[*unit][]*unit{}
	}
	seen := map[*unit]bool{}
	var out []*unit
	var visit func(*unit)
	visit = func(x *unit) {
		if x == nil || seen[x] {
			return
		}
		seen[x] = true
		out = append(out, x)
		for _, imp := range x.imports {
			visit(imp.target)
		}
	}
	visit(u)
	c.reachableCache[u] = out
	return out
}

func (c *compiler) run(entry string) Result {
	u := c.unit(entry)
	files := c.reachable(u)
	for _, f := range files {
		c.buildScope(f)
	}
	var diags []diag.Diagnostic
	for _, f := range files {
		for _, d := range f.diags {
			if d.Code == diag.CANCELED {
				panic(errCanceled)
			}
		}
		c.poll(true)
		diags = append(diags, f.diags...)
		v := &validator{c: c, u: f}
		v.validateProgram()
		diags = append(diags, v.diags...)
		b := &binder{c: c, u: f}
		b.bind()
		diags = append(diags, b.diags...)
	}
	diags = append(diags, c.limitDiags...)
	res := Result{Program: u.prog}
	earlyErrors := diag.List(diags).HasErrors()

	// Always interpret, for the best-effort Partial model. Like @dbml/parse,
	// interpreter diagnostics only count when the earlier stages are clean.
	db, interpDiags := c.interpretSafely(u)
	c.poll(true)
	res.Partial = db
	if !earlyErrors {
		diags = append(diags, interpDiags...)
	}
	res.Diagnostics = capDiagnostics(diags, c.opts.MaxDiagnostics)
	if !diag.List(diags).HasErrors() {
		res.Database = db
	}
	return res
}

// interpretSafely runs the interpreter, converting a panic on malformed
// (already diagnosed) input into a nil partial model.
func (c *compiler) interpretSafely(u *unit) (db *Database, diags []diag.Diagnostic) {
	defer func() {
		if r := recover(); r != nil {
			if r == errCanceled {
				panic(r)
			}
			db, diags = nil, []diag.Diagnostic{diag.New(diag.INTERNAL_ERROR, fmt.Sprintf("internal error in go-dbml: %v", r), u.path,
				diag.Range{Start: Position{Line: 1, Column: 1}, End: Position{Line: 1, Column: 1}})}
		}
	}()
	in := &interpreter{c: c, u: u}
	db = in.interpret()
	return db, in.diags
}

// capDiagnostics keeps the first max diagnostics and appends one
// TOO_MANY_DIAGNOSTICS entry (an error if any dropped entry was an error).
func capDiagnostics(ds []diag.Diagnostic, max int) diag.List {
	if max <= 0 || len(ds) <= max {
		return ds
	}
	sev := diag.SeverityWarning
	for _, d := range ds[max:] {
		if d.Severity == diag.SeverityError {
			sev = diag.SeverityError
		}
	}
	last := ds[max-1]
	out := append(diag.List{}, ds[:max]...)
	return append(out, diag.NewWithSeverity(sev, diag.TOO_MANY_DIAGNOSTICS, fmt.Sprintf("Too many diagnostics; %d more were omitted", len(ds)-max), last.Filepath, last.Range))
}

// ---------------------------------------------------------------- positions

func tokenPos(p lexer.Position) Position { return lexer.ToPos(p) }

func nodeToken(n ast.Node, file string) *Token {
	b := n.Base()
	if !b.HasPos {
		return &Token{Filepath: file}
	}
	return &Token{Start: tokenPos(b.Start), End: tokenPos(b.End), Filepath: file}
}

func nodeRange(n ast.Node) diag.Range {
	b := n.Base()
	return diag.Range{Start: tokenPos(b.Start), End: tokenPos(b.End)}
}

func nodeDiag(code diag.Code, msg string, n ast.Node, file string) diag.Diagnostic {
	return diag.New(code, msg, file, nodeRange(n))
}

func tokDiag(code diag.Code, msg string, t *lexer.Token, file string) diag.Diagnostic {
	return diag.New(code, msg, file, diag.Range{Start: tokenPos(t.Start), End: tokenPos(t.End)})
}

// ---------------------------------------------------------------- symbols

type symKind int

const (
	symTable symKind = iota
	symEnum
	symPartial
	symGroup
	symNote
	symView
)

func (k symKind) String() string {
	return [...]string{"Table", "Enum", "TablePartial", "TableGroup", "Note", "DiagramView"}[k]
}

// sym is a top-level named element.
type sym struct {
	kind   symKind
	name   string
	schema string // dotted schema chain; "public" when omitted
	alias  string
	decl   *ast.ElementDeclaration
	unit   *unit
	// declared by `Table a.b.c` (nested schema)
	nested bool
	// placeholder marks an unresolved import (no declaration).
	placeholder bool
}

func (s *sym) qualified() string { return s.schema + "." + s.name }

// col is a column visible in a table (own or injected from a partial).
type col struct {
	name    string
	fa      *ast.FunctionApplication
	partial *sym // non-nil when injected from a partial
}

func (c *compiler) collectDecls(u *unit) {
	for _, d := range u.prog.Declarations() {
		var kind symKind
		switch {
		case d.IsKind(kTable):
			kind = symTable
		case d.IsKind(kEnum):
			kind = symEnum
		case d.IsKind(kTablePartial):
			kind = symPartial
		case d.IsKind(kTableGroup):
			kind = symGroup
		case d.IsKind(kNote):
			kind = symNote
		case d.IsKind(kDiagramView):
			kind = symView
		default:
			continue
		}
		parts, ok := destructureComplexVariable(d.Name)
		if !ok || len(parts) == 0 {
			continue
		}
		if _, isArr := d.Name.(*ast.Array); isArr {
			continue
		}
		s := &sym{kind: kind, decl: d, unit: u, name: parts[len(parts)-1], schema: defaultSchema}
		if len(parts) > 1 {
			s.schema = strings.Join(parts[:len(parts)-1], ".")
			s.nested = len(parts) > 2
		}
		if kind == symTable && !ast.IsNil(d.Alias) && isSimpleName(d.Alias) {
			s.alias, _ = extractVariable(d.Alias)
		}
		u.decls = append(u.decls, s)
	}
}

// visible is a symbol as seen from a scope.
type visible struct {
	sym    *sym
	schema string
	name   string
	spec   ast.Node // use specifier that imported it (nil for local)
}

type scope struct {
	u        *unit
	entries  []*visible
	byKey    map[string][]*visible // kind|schema|name
	bySym    map[*sym]*visible     // first visible entry of each symbol
	aliases  map[string]*sym       // table aliases (program level)
	schemas  map[string]bool
	partials map[*sym][]*sym
}

func key(k symKind, schema, name string) string {
	return k.String() + "|" + schema + "|" + name
}

func (sc *scope) add(v *visible) {
	k := key(v.sym.kind, v.schema, v.name)
	for _, e := range sc.byKey[k] {
		if e.sym == v.sym {
			return
		}
	}
	sc.entries = append(sc.entries, v)
	sc.byKey[k] = append(sc.byKey[k], v)
	if _, ok := sc.bySym[v.sym]; !ok {
		sc.bySym[v.sym] = v
	}
	sc.addSchema(v.schema)
}

func (sc *scope) addSchema(s string) {
	parts := strings.Split(s, ".")
	for i := range parts {
		sc.schemas[strings.Join(parts[:i+1], ".")] = true
	}
}

func (c *compiler) buildScope(u *unit) {
	if u.scope != nil {
		return
	}
	sc := &scope{u: u, byKey: map[string][]*visible{}, bySym: map[*sym]*visible{}, aliases: map[string]*sym{}, schemas: map[string]bool{defaultSchema: true}}
	u.scope = sc
	for _, s := range u.decls {
		sc.add(&visible{sym: s, schema: s.schema, name: s.name})
		if s.alias != "" {
			if _, exists := sc.aliases[s.alias]; !exists {
				sc.aliases[s.alias] = s
			}
		}
	}
	for _, imp := range u.imports {
		c.applyImport(sc, imp)
	}
}

// exported returns the symbols a unit exposes to importers: its own
// declarations plus everything it re-exports via `reuse`.
func (c *compiler) exported(u *unit, seen map[*unit]bool) []*visible {
	if u == nil || seen[u] {
		return nil
	}
	seen[u] = true
	var out []*visible
	for _, s := range u.decls {
		if s.kind == symView {
			continue
		}
		out = append(out, &visible{sym: s, schema: s.schema, name: s.name})
	}
	for _, imp := range u.imports {
		if !imp.reuse || imp.target == nil {
			continue
		}
		out = append(out, c.selectImport(imp, seen)...)
	}
	return out
}

func (c *compiler) selectImport(imp *importDecl, seen map[*unit]bool) []*visible {
	avail := c.exported(imp.target, seen)
	if imp.wildcard {
		return avail
	}
	var out []*visible
	for _, spec := range imp.specs {
		if spec.ImportKind == nil || ast.IsNil(spec.Name) {
			continue
		}
		parts, ok := destructureComplexVariable(spec.Name)
		if !ok || len(parts) == 0 {
			continue
		}
		kind := strings.ToLower(spec.ImportKind.Value)
		aliasName := ""
		if !ast.IsNil(spec.Alias) {
			aliasName, _ = extractVariable(spec.Alias)
		}
		name := parts[len(parts)-1]
		schema := defaultSchema
		if len(parts) > 1 {
			schema = strings.Join(parts[:len(parts)-1], ".")
		}
		switch kind {
		case "schema":
			full := strings.Join(parts, ".")
			for _, v := range avail {
				if v.schema == full || strings.HasPrefix(v.schema, full+".") {
					nv := *v
					if aliasName != "" {
						nv.schema = aliasName + strings.TrimPrefix(v.schema, full)
					}
					nv.spec = spec
					out = append(out, &nv)
				}
			}
		case "tablegroup":
			for _, v := range avail {
				if v.sym.kind == symGroup && v.name == name && v.schema == schema {
					nv := *v
					if aliasName != "" {
						nv.name = aliasName
					}
					nv.spec = spec
					out = append(out, &nv)
					for _, t := range groupTableRefs(v.sym, avail) {
						tv := *t
						tv.spec = spec
						out = append(out, &tv)
					}
				}
			}
		default:
			want := map[string]symKind{"table": symTable, "enum": symEnum, "tablepartial": symPartial, "note": symNote}
			k, ok := want[kind]
			if !ok {
				continue
			}
			for _, v := range avail {
				if v.sym.kind == k && v.name == name && v.schema == schema {
					nv := *v
					if aliasName != "" {
						nv.name = aliasName
						nv.schema = defaultSchema
					}
					nv.spec = spec
					out = append(out, &nv)
				}
			}
		}
	}
	return out
}

func groupTableRefs(g *sym, avail []*visible) []*visible {
	var out []*visible
	fields, _ := partitionBody(elementBody(g.decl))
	for _, f := range fields {
		parts, ok := destructureComplexVariable(f.Callee)
		if !ok || len(parts) == 0 {
			continue
		}
		name := parts[len(parts)-1]
		schema := defaultSchema
		if len(parts) > 1 {
			schema = strings.Join(parts[:len(parts)-1], ".")
		}
		for _, v := range avail {
			if v.sym.kind == symTable && ((v.name == name && v.schema == schema) || (len(parts) == 1 && v.sym.alias == name)) {
				out = append(out, v)
			}
		}
	}
	return out
}

func (c *compiler) applyImport(sc *scope, imp *importDecl) {
	var selected []*visible
	if imp.target != nil {
		selected = c.selectImport(imp, map[*unit]bool{})
	}
	for _, v := range selected {
		sc.add(v)
		if v.sym.kind == symTable && v.sym.alias != "" && v.name == v.sym.name {
			if _, exists := sc.aliases[v.sym.alias]; !exists {
				sc.aliases[v.sym.alias] = v.sym
			}
		}
	}
	for _, spec := range imp.specs {
		if spec.ImportKind == nil {
			continue
		}
		kind := strings.ToLower(spec.ImportKind.Value)
		if kind == "schema" {
			if parts, ok := destructureComplexVariable(spec.Name); ok {
				name := strings.Join(parts, ".")
				if !ast.IsNil(spec.Alias) {
					name, _ = extractVariable(spec.Alias)
				}
				sc.addSchema(name)
			}
			continue
		}
		k, ok := importSymKinds[kind]
		if !ok {
			continue
		}
		found := false
		for _, v := range selected {
			if v.spec == ast.Node(spec) && v.sym.kind == k {
				found = true
			}
		}
		if found {
			continue
		}
		// Like @dbml/parse's UseSymbol with an unresolved original: the name
		// is still declared in this file (so later references bind to it and
		// report only the deeper failure), and binding reports the import.
		parts, ok := destructureComplexVariable(spec.Name)
		if !ok || len(parts) == 0 {
			continue
		}
		ph := &sym{kind: k, name: parts[len(parts)-1], schema: defaultSchema, unit: sc.u, placeholder: true}
		if len(parts) > 1 {
			ph.schema = strings.Join(parts[:len(parts)-1], ".")
		}
		name := ph.name
		schema := ph.schema
		if !ast.IsNil(spec.Alias) {
			name, _ = extractVariable(spec.Alias)
			schema = defaultSchema
		}
		sc.add(&visible{sym: ph, schema: schema, name: name, spec: spec})
		sc.u.unresolved = append(sc.u.unresolved, unresolvedImport{spec: spec, kind: k, display: strings.Join(parts, ".")})
	}
}

var importSymKinds = map[string]symKind{"table": symTable, "enum": symEnum, "tablepartial": symPartial, "note": symNote, "tablegroup": symGroup}

type unresolvedImport struct {
	spec    *ast.UseSpecifier
	kind    symKind
	display string
}

// lookup finds a symbol of kind in schema (schema "" = program level).
func (sc *scope) lookup(k symKind, schema, name string) *visible {
	if schema == "" {
		schema = defaultSchema
	}
	if vs := sc.byKey[key(k, schema, name)]; len(vs) > 0 {
		return vs[0]
	}
	if k == symTable && schema == defaultSchema {
		if s, ok := sc.aliases[name]; ok {
			return sc.visibleOf(s)
		}
	}
	return nil
}

func (sc *scope) visibleOf(s *sym) *visible {
	if v, ok := sc.bySym[s]; ok {
		return v
	}
	return &visible{sym: s, schema: s.schema, name: s.name}
}

func (sc *scope) hasSchema(s string) bool { return sc.schemas[s] }

// interpretedName is the (schema, name) of s as seen from this scope, with
// the public schema reported as nil.
func (sc *scope) interpretedName(s *sym) (*string, string) {
	if s.unit == sc.u {
		return schemaPtr(s.schema), s.name
	}
	if v, ok := sc.bySym[s]; ok {
		if v.name != s.name {
			return strPtr(""), v.name
		}
		return schemaPtr(v.schema), v.name
	}
	return schemaPtr(s.schema), s.name
}

func schemaPtr(s string) *string {
	if s == defaultSchema || s == "" {
		return nil
	}
	return strPtr(s)
}

// columns returns the merged columns of a table: own columns plus columns
// injected from partials. Own columns win over injected ones; among
// injected columns the last injection wins (its position is kept).
func (c *compiler) columns(t *sym, sc *scope) []*col { return c.columnSet(t, sc).list }

func (c *compiler) columnSet(t *sym, sc *scope) *colSet {
	k := colCacheKey{t, sc}
	if cs, ok := c.colCache[k]; ok {
		return cs
	}
	if c.colCache == nil {
		c.colCache = map[colCacheKey]*colSet{}
	}
	// guard against partial-injection cycles while computing
	c.colCache[k] = &colSet{byName: map[string]*col{}}
	list := c.computeColumns(t, sc)
	cs := &colSet{list: list, byName: make(map[string]*col, len(list))}
	for _, cl := range list {
		cs.byName[cl.name] = cl
	}
	c.colCache[k] = cs
	return cs
}

func (c *compiler) computeColumns(t *sym, sc *scope) []*col {
	type member struct {
		c        *col
		inject   bool
		injected bool
	}
	var members []member
	fields, _ := partitionBody(elementBody(t.decl))
	for _, f := range fields {
		if isValidPartialInjection(f.Callee) {
			pname, _ := extractVariable(f.Callee.(*ast.PrefixExpression).Expression)
			members = append(members, member{inject: true, c: &col{name: pname, fa: f}})
			continue
		}
		name, ok := extractVariable(f.Callee)
		if !ok {
			continue
		}
		members = append(members, member{c: &col{name: name, fa: f}})
	}
	if t.kind == symTable {
		var expanded []member
		for _, m := range members {
			if m.inject {
				if p := c.lookupPartial(sc, m.c.name); p != nil {
					for _, pc := range c.columns(p.sym, p.sym.unit.scope) {
						expanded = append(expanded, member{c: &col{name: pc.name, fa: pc.fa, partial: p.sym}, injected: true})
					}
				}
			}
			expanded = append(expanded, m)
		}
		members = expanded
	}
	direct := map[string]bool{}
	for _, m := range members {
		if !m.inject && !m.injected {
			direct[m.c.name] = true
		}
	}
	var filtered []member
	for _, m := range members {
		if m.inject {
			continue
		}
		if m.injected && direct[m.c.name] {
			continue
		}
		filtered = append(filtered, m)
	}
	// keep last occurrence per name
	seen := map[string]bool{}
	var rev []*col
	for i := len(filtered) - 1; i >= 0; i-- {
		n := filtered[i].c.name
		if seen[n] {
			continue
		}
		seen[n] = true
		rev = append(rev, filtered[i].c)
	}
	out := make([]*col, len(rev))
	for i := range rev {
		out[i] = rev[len(rev)-1-i]
	}
	return out
}

func (c *compiler) lookupColumn(t *sym, sc *scope, name string) *col {
	return c.columnSet(t, sc).byName[name]
}

func (c *compiler) lookupPartial(sc *scope, name string) *visible {
	if v := sc.lookup(symPartial, defaultSchema, name); v != nil {
		return v
	}
	for _, v := range sc.entries {
		if v.sym.kind == symPartial && v.name == name {
			return v
		}
	}
	return nil
}

// injectedPartials returns the partials injected into table t, in order.
func (c *compiler) injectedPartials(t *sym, sc *scope) []*sym {
	var out []*sym
	fields, _ := partitionBody(elementBody(t.decl))
	for _, f := range fields {
		if !isValidPartialInjection(f.Callee) {
			continue
		}
		pname, _ := extractVariable(f.Callee.(*ast.PrefixExpression).Expression)
		if p := c.lookupPartial(sc, pname); p != nil {
			out = append(out, p.sym)
		}
	}
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Validate lexes, parses and structurally validates a single file without
// resolving references. It returns syntax and per-element diagnostics only
// (the "validate" stage of @dbml/parse) and is cheap enough to run on every
// keystroke.
func Validate(src, filepath string) diag.List {
	c := &compiler{opts: DefaultOptions(), load: func(p string) (string, bool) {
		if p == filepath {
			return src, true
		}
		return "", false
	}, units: map[string]*unit{}}
	u := c.unit(filepath)
	v := &validator{c: c, u: u}
	v.validateProgram()
	return capDiagnostics(append(append(diag.List{}, u.diags...), v.diags...), c.opts.MaxDiagnostics)
}
