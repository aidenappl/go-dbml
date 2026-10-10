package emit

import (
	"errors"
	"fmt"

	"github.com/aidenappl/go-dbml/model"
)

// ErrInvalidModel is wrapped by errors for structurally invalid models
// (nil elements, refs without endpoints, ...), e.g. hand-written or
// truncated model JSON.
var ErrInvalidModel = errors.New("emit: invalid model")

func invalid(what string, i int) error {
	return fmt.Errorf("%w: nil %s at index %d", ErrInvalidModel, what, i)
}

// validate rejects nil elements the emitter would dereference.
func validate(db *model.Database) error {
	for i, t := range db.Tables {
		if t == nil {
			return invalid("table", i)
		}
		if err := validateColumns(t.Fields, t.Indexes, t.Checks); err != nil {
			return fmt.Errorf("table %q: %w", t.Name, err)
		}
		for j, p := range t.Partials {
			if p == nil {
				return invalid("partial injection", j)
			}
		}
	}
	for i, p := range db.TablePartials {
		if p == nil {
			return invalid("table partial", i)
		}
		if err := validateColumns(p.Fields, p.Indexes, p.Checks); err != nil {
			return fmt.Errorf("partial %q: %w", p.Name, err)
		}
	}
	for i, r := range db.Refs {
		if r == nil {
			return invalid("ref", i)
		}
		for _, e := range r.Endpoints {
			if e.TableName == "" {
				return fmt.Errorf("%w: ref %d has an endpoint without a table", ErrInvalidModel, i)
			}
		}
	}
	for i, d := range db.Deps {
		if d == nil {
			return invalid("dep", i)
		}
		if len(d.Edges) == 0 {
			return fmt.Errorf("%w: dep %d has no edges", ErrInvalidModel, i)
		}
		for j, e := range d.Edges {
			if e == nil {
				return invalid("dep edge", j)
			}
		}
	}
	for i, e := range db.Enums {
		if e == nil {
			return invalid("enum", i)
		}
		for j, v := range e.Values {
			if v == nil {
				return invalid("enum value", j)
			}
		}
	}
	for i, g := range db.TableGroups {
		if g == nil {
			return invalid("table group", i)
		}
		for j, it := range g.Tables {
			if it == nil {
				return invalid("table group item", j)
			}
		}
	}
	for i, n := range db.Notes {
		if n == nil {
			return invalid("note", i)
		}
	}
	for i, r := range db.Records {
		if r == nil {
			return invalid("records", i)
		}
	}
	for i, v := range db.DiagramViews {
		if v == nil {
			return invalid("diagram view", i)
		}
	}
	for _, list := range [][]*model.ElementRef{db.Externals.Tables, db.Externals.Enums, db.Externals.TablePartials, db.Externals.TableGroups, db.Externals.Notes} {
		for i, x := range list {
			if x == nil {
				return invalid("external", i)
			}
		}
	}
	return nil
}

func validateColumns(cols []*model.Column, ixs []*model.Index, checks []*model.Check) error {
	for i, c := range cols {
		if c == nil {
			return invalid("column", i)
		}
		for j, r := range c.InlineRefs {
			if r == nil {
				return invalid("inline ref", j)
			}
		}
		for j, d := range c.InlineDeps {
			if d == nil {
				return invalid("inline dep", j)
			}
		}
		for j, ch := range c.Checks {
			if ch == nil {
				return invalid("column check", j)
			}
		}
	}
	for i, ix := range ixs {
		if ix == nil {
			return invalid("index", i)
		}
		for j, c := range ix.Columns {
			if c == nil {
				return invalid("index column", j)
			}
		}
	}
	for i, ch := range checks {
		if ch == nil {
			return invalid("check", i)
		}
	}
	return nil
}
