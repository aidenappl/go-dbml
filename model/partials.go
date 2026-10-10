// Portions of this file are derived from holistics/dbml packages/dbml-core/src/model_structure/table.ts (processPartials)
// (Apache-2.0, Copyright 2019 Holistics Software Pte Ltd; commit 34108bb).
// Modified: ported to Go.

package model

import "sort"

// ResolvedColumn is a column of a table after partial injection.
type ResolvedColumn struct {
	*Column
	// InjectedFrom names the TablePartial the column came from ("" for own columns).
	InjectedFrom string `json:"injectedFrom,omitempty"`
}

// ResolvedIndex is an index of a table after partial injection.
type ResolvedIndex struct {
	*Index
	InjectedFrom string `json:"injectedFrom,omitempty"`
}

// ResolvedCheck is a check of a table after partial injection. Checks
// injected from a partial have their name prefixed with "<table>." (as
// @dbml/core does) so names stay unique across tables.
type ResolvedCheck struct {
	Token        *Token  `json:"token"`
	Expression   string  `json:"expression"`
	Name         *string `json:"name,omitempty"`
	InjectedFrom string  `json:"injectedFrom,omitempty"`
}

// ResolvedTable is a table with every `~partial` injection applied.
type ResolvedTable struct {
	*Table
	Columns     []*ResolvedColumn `json:"columns"`
	AllIndexes  []*ResolvedIndex  `json:"allIndexes"`
	AllChecks   []*ResolvedCheck  `json:"allChecks"`
	Note        *NoteValue        `json:"resolvedNote,omitempty"`
	HeaderColor *string           `json:"resolvedHeaderColor,omitempty"`
	// InjectedRefs are refs created by `[ref: ...]` settings on injected
	// partial columns (endpoint 0 is this table's column).
	InjectedRefs []*Ref `json:"injectedRefs,omitempty"`
}

// FindTablePartial returns the partial with the given name, or nil.
func (db *Database) FindTablePartial(name string) *TablePartial {
	for _, p := range db.TablePartials {
		if p.Name == name {
			return p
		}
	}
	return nil
}

// ResolveTable applies the table's partial injections using the precedence
// documented by DBML (and implemented by @dbml/core):
//
//  1. Columns, note and headerColor defined locally on the table win.
//  2. Otherwise the last injected partial (in source order) that defines a
//     conflicting column/setting wins.
//
// Injected columns appear at the position of their `~partial` line. Indexes
// and checks of the injected partials come first (last partial first) and
// the table's own follow — @dbml/core's Table runs processPartials before
// processIndexes/processChecks.
func (db *Database) ResolveTable(t *Table) *ResolvedTable {
	rt := &ResolvedTable{Table: t, Note: t.Note, HeaderColor: t.HeaderColor}
	// own appends the table's own indexes/checks after the injected ones.
	own := func() *ResolvedTable {
		for _, ix := range t.Indexes {
			rt.AllIndexes = append(rt.AllIndexes, &ResolvedIndex{Index: ix})
		}
		for _, c := range t.Checks {
			rt.AllChecks = append(rt.AllChecks, &ResolvedCheck{Token: c.Token, Expression: c.Expression, Name: c.Name})
		}
		return rt
	}
	fields := make([]*ResolvedColumn, 0, len(t.Fields))
	for _, c := range t.Fields {
		fields = append(fields, &ResolvedColumn{Column: c})
	}
	if len(t.Partials) == 0 {
		rt.Columns = fields
		return own()
	}
	existing := map[string]bool{}
	for _, c := range t.Fields {
		existing[c.Name] = true
	}
	hasNote, hasColor := t.Note != nil, t.HeaderColor != nil

	parts := append([]*PartialInject{}, t.Partials...)
	sort.SliceStable(parts, func(i, j int) bool { return parts[i].Order > parts[j].Order })
	// placeholders, ascending order
	for i := len(parts) - 1; i >= 0; i-- {
		at := clamp(parts[i].Order, len(fields))
		fields = append(fields[:at], append([]*ResolvedColumn{nil}, fields[at:]...)...)
	}
	for _, inj := range parts {
		at := clamp(inj.Order, len(fields)-1)
		p := db.FindTablePartial(inj.Name)
		var injected []*ResolvedColumn
		if p != nil {
			for _, f := range p.Fields {
				if existing[f.Name] {
					continue
				}
				existing[f.Name] = true
				injected = append(injected, &ResolvedColumn{Column: f, InjectedFrom: p.Name})
				for _, ir := range f.InlineRefs {
					target, targetSchema := ir.TableName, ir.SchemaName
					// A bare column ref inside a partial (`[ref: > other_col]`)
					// names the partial itself; it means the table being built.
					if targetSchema == nil && target == p.Name && db.findTable(nil, p.Name) == nil {
						target, targetSchema = t.Name, t.SchemaName
					}
					m, ok := relationshipOps[ir.Relation]
					if !ok {
						m = [2]string{"*", "*"}
					}
					rt.InjectedRefs = append(rt.InjectedRefs, &Ref{
						Token: f.Token,
						Endpoints: [2]RefEndpoint{
							{TableName: t.Name, SchemaName: t.SchemaName, FieldNames: []string{f.Name}, Relation: m[0], Token: f.Token},
							{TableName: target, SchemaName: targetSchema, FieldNames: ir.FieldNames, Relation: m[1], Token: ir.Token},
						},
					})
				}
			}
		}
		fields = append(fields[:at], append(injected, fields[at+1:]...)...)
		if p == nil {
			continue
		}
		if !hasNote && p.Note != nil {
			rt.Note = p.Note
			hasNote = true
		}
		if !hasColor && p.HeaderColor != nil {
			rt.HeaderColor = p.HeaderColor
			hasColor = true
		}
		for _, ix := range p.Indexes {
			rt.AllIndexes = append(rt.AllIndexes, &ResolvedIndex{Index: ix, InjectedFrom: p.Name})
		}
		for _, c := range p.Checks {
			rc := &ResolvedCheck{Token: c.Token, Expression: c.Expression, InjectedFrom: p.Name}
			if c.Name != nil {
				rc.Name = strPtr(t.Name + "." + *c.Name)
			}
			rt.AllChecks = append(rt.AllChecks, rc)
		}
	}
	rt.Columns = fields
	return own()
}

func clamp(i, n int) int {
	if i < 0 {
		return 0
	}
	if i > n {
		return n
	}
	return i
}

func (db *Database) findTable(schema *string, name string) *Table {
	for _, t := range db.Tables {
		if t.Name == name && ((t.SchemaName == nil && schema == nil) || (t.SchemaName != nil && schema != nil && *t.SchemaName == *schema)) {
			return t
		}
	}
	return nil
}
