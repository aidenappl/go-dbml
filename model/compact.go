package model

import "reflect"

// StripTokens removes every source token from db (in place) for a compact
// JSON encoding — tokens are most of the model JSON's size. Tables,
// columns, refs, enums, table groups, partials and sticky notes keep their
// 1-based start line in Line (for click-to-source). Use it on a model you
// own; it is not safe while other goroutines read db.
func StripTokens(db *Database) {
	if db == nil {
		return
	}
	setLine := func(l *int, t *Token) {
		if t != nil { // a second StripTokens keeps the recorded lines
			*l = t.Start.Line
		}
	}
	for _, t := range db.Tables {
		if t == nil {
			continue
		}
		setLine(&t.Line, t.Token)
		for _, c := range t.Fields {
			if c != nil {
				setLine(&c.Line, c.Token)
			}
		}
	}
	for _, r := range db.Refs {
		if r != nil {
			setLine(&r.Line, r.Token)
		}
	}
	for _, e := range db.Enums {
		if e != nil {
			setLine(&e.Line, e.Token)
		}
	}
	for _, g := range db.TableGroups {
		if g != nil {
			setLine(&g.Line, g.Token)
		}
	}
	for _, p := range db.TablePartials {
		if p == nil {
			continue
		}
		setLine(&p.Line, p.Token)
		for _, c := range p.Fields {
			if c != nil {
				setLine(&c.Line, c.Token)
			}
		}
	}
	for _, n := range db.Notes {
		if n != nil {
			setLine(&n.Line, n.Token)
		}
	}
	nilTokens(reflect.ValueOf(db), map[uintptr]bool{})
}

var tokenPtrType = reflect.TypeOf((*Token)(nil))

// nilTokens sets every *Token reachable from v to nil.
func nilTokens(v reflect.Value, seen map[uintptr]bool) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		if seen[v.Pointer()] {
			return
		}
		seen[v.Pointer()] = true
		nilTokens(v.Elem(), seen)
	case reflect.Interface:
		if !v.IsNil() {
			nilTokens(v.Elem(), seen)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if !f.CanSet() {
				continue
			}
			if f.Type() == tokenPtrType {
				f.Set(reflect.Zero(tokenPtrType))
				continue
			}
			nilTokens(f, seen)
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			nilTokens(v.Index(i), seen)
		}
	case reflect.Map:
		it := v.MapRange()
		for it.Next() {
			nilTokens(it.Value(), seen) // map values are not addressable: only pointers inside are reached
		}
	}
}
