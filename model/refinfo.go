// Portions of this file are derived from holistics/dbml packages/dbml-parse/src/core/global_modules/ref/constraint_fixes.ts, src/core/utils/diagnostics_reporter/ref.ts
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import (
	"strings"

	"github.com/aidenappl/go-dbml/ast"
	"github.com/aidenappl/go-dbml/diag"
)

// ---------------------------------------------------------------- ref infos

// validateRefCardinalities reports INVALID_REF_RELATIONSHIP infos when a ref
// operator disagrees with column nullability/uniqueness.
func (in *interpreter) validateRefCardinalities(db *Database) {
	file := in.u.path
	for _, e := range in.u.prog.Declarations() {
		in.c.poll(false)
		switch {
		case e.IsKind(kRef):
			inx := refInfix(e)
			if inx == nil || inx.Op == nil {
				continue
			}
			l, _ := in.c.resolveEndpoint(in.sc, inx.Left, modeRef, nil)
			r, _ := in.c.resolveEndpoint(in.sc, inx.Right, modeRef, nil)
			if !l.ok || !r.ok {
				continue
			}
			in.cardinality(inx.Op.Value, l, r, inx.Left, inx.Right, file)
		case e.IsKind(kTable), e.IsKind(kTablePartial):
			t := in.u.symOfDecl(e)
			if t == nil {
				continue
			}
			fields, _ := partitionBody(elementBody(e))
			for _, f := range fields {
				cname, ok := extractVariable(f.Callee)
				if !ok {
					continue
				}
				lst := trailingList(f.Args)
				if lst == nil {
					continue
				}
				set, _ := aggregateSettings(lst, file)
				for _, a := range set.get(settingRef) {
					pre, ok := a.Value.(*ast.PrefixExpression)
					if !ok || pre.Op == nil || !isRelationshipOp(pre.Op.Value) {
						continue
					}
					r, _ := in.c.resolveEndpoint(in.sc, pre.Expression, modeInline, t)
					cl := in.c.lookupColumn(t, t.unit.scope, cname)
					if !r.ok || cl == nil {
						continue
					}
					l := resolved{table: in.sc.visibleOf(t), cols: []*col{cl}, ok: true}
					in.cardinality(pre.Op.Value, l, r, f, a, file)
				}
			}
		}
	}
}

type card struct {
	min  int
	many bool
}

func parseCard(c string) card {
	switch c {
	case "1":
		return card{1, false}
	case "*":
		return card{1, true}
	case "0..1":
		return card{0, false}
	case "0..*":
		return card{0, true}
	}
	return card{1, false}
}

func (in *interpreter) cardinality(op string, l, r resolved, leftNode, rightNode ast.Node, file string) {
	m, ok := relationshipOps[op]
	if !ok {
		return
	}
	for _, side := range []string{"right", "left"} {
		own, other := parseCard(m[0]), parseCard(m[1])
		ownRes, otherRes := l, r
		ownNode, otherNode := leftNode, rightNode
		if side == "right" {
			own, other = other, own
			ownRes, otherRes = r, l
			ownNode, otherNode = rightNode, leftNode
		}
		otherIsSource := other.many || (!own.many && !other.many && side == "left")
		push := func(node ast.Node, res resolved, msg string) {
			in.diags = append(in.diags, diag.NewWithSeverity(diag.SeverityInfo, diag.INVALID_REF_RELATIONSHIP, msg, file, nodeRange(node)))
			for _, c := range res.cols {
				if c.fa != nil && ast.Node(c.fa) != node {
					in.diags = append(in.diags, diag.NewWithSeverity(diag.SeverityInfo, diag.INVALID_REF_RELATIONSHIP, msg, res.table.sym.unit.path, nodeRange(c.fa)))
				}
			}
		}
		names := in.qualifiedNames(otherRes)
		if own.min >= 1 && otherIsSource && len(otherRes.cols) > 0 && in.allNullable(otherRes) {
			push(otherNode, otherRes, nullableMsg(names, op, len(otherRes.cols) > 1))
		}
		if own.min == 0 && otherIsSource && len(otherRes.cols) > 0 && in.allNotNull(otherRes) {
			push(otherNode, otherRes, "Column `"+names+"` is NOT NULL but operator `"+op+"` allows it to be optional")
		}
		ownNames := in.qualifiedNames(ownRes)
		if own.many && in.columnsUnique(ownRes) {
			push(ownNode, ownRes, "Column `"+ownNames+"` is unique but operator `"+op+"` allows many")
		}
		if !own.many && !in.columnsUnique(ownRes) {
			push(ownNode, ownRes, "Column `"+ownNames+"` should be unique or primary key for operator `"+op+"`")
		}
	}
}

func nullableMsg(names, op string, composite bool) string {
	if composite {
		return "Columns (`" + names + "`) are all nullable but operator `" + op + "` requires at least one to be NOT NULL"
	}
	return "Column `" + names + "` is nullable but operator `" + op + "` requires it to be NOT NULL"
}

func (in *interpreter) qualifiedNames(r resolved) string {
	var parts []string
	for _, c := range r.cols {
		schema, name := in.sc.interpretedName(r.table.sym)
		q := name + "." + c.name
		if schema != nil {
			q = *schema + "." + q
		}
		parts = append(parts, q)
	}
	return strings.Join(parts, ", ")
}

func (in *interpreter) colSettings(c *col) *settings {
	if c.fa == nil {
		return newSettings()
	}
	var rest []ast.Node
	if len(c.fa.Args) > 1 {
		rest = c.fa.Args[1:]
	}
	s, errs := columnSettings(rest, "", c.partial != nil)
	if len(errs) > 0 {
		return newSettings()
	}
	return s
}

func (in *interpreter) colPK(c *col) bool {
	s := in.colSettings(c)
	return s.has(settingPK) || s.has(settingPrimaryKey)
}

func (in *interpreter) colIncrement(c *col, owner *sym) bool {
	if in.colSettings(c).has(settingIncrement) {
		return true
	}
	return serialTypes[normalizeTypeName(in.columnType(c, owner).base)]
}

func (in *interpreter) tableIndexes(t *sym) []*Index {
	var idx []*Index
	for _, s := range append([]*sym{t}, in.c.injectedPartials(t, t.unit.scope)...) {
		_, subs := partitionBody(elementBody(s.decl))
		for _, sub := range subs {
			if sub.IsKind(kIndexes) {
				idx = append(idx, interpretIndexes(sub, s.unit.path)...)
			}
		}
	}
	return idx
}

func (in *interpreter) nullable(c *col, owner *sym) bool {
	if in.colPK(c) || in.colIncrement(c, owner) {
		return false
	}
	for _, ix := range in.tableIndexes(owner) {
		if ix.PK == nil || !*ix.PK {
			continue
		}
		for _, ic := range ix.Columns {
			if ic.Value == c.name {
				return false
			}
		}
	}
	return !in.colSettings(c).has(settingNotNull)
}

func (in *interpreter) allNullable(r resolved) bool {
	for _, c := range r.cols {
		if !in.nullable(c, r.table.sym) {
			return false
		}
	}
	return true
}

func (in *interpreter) allNotNull(r resolved) bool {
	for _, c := range r.cols {
		if in.nullable(c, r.table.sym) {
			return false
		}
	}
	return true
}

func (in *interpreter) columnsUnique(r resolved) bool {
	if len(r.cols) == 0 {
		return false
	}
	if len(r.cols) == 1 {
		c := r.cols[0]
		if in.colSettings(c).has(settingUnique) || in.colPK(c) || in.colIncrement(c, r.table.sym) {
			return true
		}
	}
	names := map[string]bool{}
	for _, c := range r.cols {
		names[c.name] = true
	}
	for _, ix := range in.tableIndexes(r.table.sym) {
		if (ix.Unique == nil || !*ix.Unique) && (ix.PK == nil || !*ix.PK) {
			continue
		}
		if len(ix.Columns) != len(names) {
			continue
		}
		all := true
		for _, ic := range ix.Columns {
			if !names[ic.Value] {
				all = false
			}
		}
		if all {
			return true
		}
	}
	return false
}
