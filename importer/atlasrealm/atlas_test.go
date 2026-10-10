//go:build !js

package atlasrealm_test

import (
	"strings"
	"testing"

	"ariga.io/atlas/sql/schema"

	"github.com/aidenappl/go-dbml/emit"
	"github.com/aidenappl/go-dbml/importer/atlasrealm"
)

// autoInc stands in for mysql.AutoIncrement (matched by type name).
type AutoIncrement struct{ schema.Attr }

func TestFromRealm(t *testing.T) {
	s := &schema.Schema{Name: "shop"}
	users := &schema.Table{Name: "users", Schema: s}
	id := &schema.Column{Name: "id", Type: &schema.ColumnType{Type: &schema.IntegerType{T: "bigint"}, Raw: "bigint unsigned"},
		Attrs: []schema.Attr{&AutoIncrement{}}}
	email := &schema.Column{Name: "email", Type: &schema.ColumnType{Type: &schema.StringType{T: "varchar", Size: 255}},
		Attrs: []schema.Attr{&schema.Comment{Text: "login"}}}
	status := &schema.Column{Name: "status", Type: &schema.ColumnType{Type: &schema.EnumType{Values: []string{"on", "off"}}},
		Default: &schema.Literal{V: "'on'"}}
	users.Columns = []*schema.Column{id, email, status}
	users.PrimaryKey = &schema.Index{Parts: []*schema.IndexPart{{C: id}}}
	users.Indexes = []*schema.Index{{Name: "ux_email", Unique: true, Parts: []*schema.IndexPart{{C: email}}}}
	users.Attrs = []schema.Attr{&schema.Comment{Text: "people"}, &schema.Check{Name: "ck_email", Expr: "email <> ''"}}

	orders := &schema.Table{Name: "orders", Schema: s}
	oid := &schema.Column{Name: "id", Type: &schema.ColumnType{Type: &schema.IntegerType{T: "int"}}}
	uid := &schema.Column{Name: "user_id", Type: &schema.ColumnType{Type: &schema.IntegerType{T: "bigint"}}}
	line := &schema.Column{Name: "line", Type: &schema.ColumnType{Type: &schema.IntegerType{T: "int"}}}
	created := &schema.Column{Name: "created", Type: &schema.ColumnType{Type: &schema.TimeType{T: "timestamp"}, Null: true},
		Default: &schema.RawExpr{X: "CURRENT_TIMESTAMP"}}
	orders.Columns = []*schema.Column{oid, uid, line, created}
	orders.PrimaryKey = &schema.Index{Parts: []*schema.IndexPart{{C: oid}, {C: line}}}
	orders.ForeignKeys = []*schema.ForeignKey{{Symbol: "fk_user", Table: orders, Columns: []*schema.Column{uid},
		RefTable: users, RefColumns: []*schema.Column{id}, OnDelete: schema.Cascade}}
	s.Tables = []*schema.Table{users, orders}

	db, err := atlasrealm.FromRealm(&schema.Realm{Schemas: []*schema.Schema{s}}, atlasrealm.Options{})
	if err != nil {
		t.Fatal(err)
	}
	got, err := emit.DBML(db)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Enum users_status_enum {", "id \"bigint unsigned\" [pk, increment]",
		"email varchar(255) [not null, note: 'login']", "status users_status_enum [not null, default: 'on']",
		"(id, line) [pk]", "email [unique, name: 'ux_email']", "`email <> ''` [name: 'ck_email']",
		"created timestamp [default: `CURRENT_TIMESTAMP`]", "Note: 'people'",
		"Ref fk_user: users.id <? orders.user_id [delete: cascade]",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
}
