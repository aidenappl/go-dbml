//go:build !js

package atlasrealm_test

import (
	"strings"
	"testing"

	"ariga.io/atlas/sql/schema"

	dbml "github.com/aidenappl/go-dbml"
	"github.com/aidenappl/go-dbml/export"
	"github.com/aidenappl/go-dbml/importer"
	"github.com/aidenappl/go-dbml/importer/atlasrealm"
)

// Identity stands in for postgres.Identity (matched by type name).
type Identity struct{ schema.Attr }

func col(name string, t schema.Type, raw string, null bool, attrs ...schema.Attr) *schema.Column {
	return &schema.Column{Name: name, Type: &schema.ColumnType{Type: t, Raw: raw, Null: null}, Attrs: attrs}
}

// A MariaDB 11.4-shaped realm: backtick-quoted expressions everywhere,
// every inline ENUM typed "enum", a foreign key into a schema that was not
// inspected.
func mariadbRealm() *schema.Realm {
	app := &schema.Schema{Name: "app"}
	billing := &schema.Schema{Name: "billing"} // not part of the realm
	accounts := &schema.Table{Name: "accounts", Schema: billing}
	accID := col("id", &schema.IntegerType{T: "int"}, "int(11)", false)
	accounts.Columns = []*schema.Column{accID}

	users := &schema.Table{Name: "users", Schema: app}
	id := col("id", &schema.IntegerType{T: "int", Unsigned: true}, "int(10) unsigned", false, &AutoIncrement{})
	name := col("name", &schema.StringType{T: "varchar", Size: 64}, "varchar(64)", false)
	key := col("key", &schema.StringType{T: "varchar", Size: 16}, "varchar(16)", true)
	status := col("status", &schema.EnumType{T: "enum", Values: []string{"on", "off"}}, "enum('on','off')", false)
	full := col("full", &schema.StringType{T: "varchar", Size: 100}, "varchar(100)", true,
		&schema.GeneratedExpr{Expr: "concat(`name`,' ',`key`)", Type: "VIRTUAL"})
	created := col("created", &schema.TimeType{T: "datetime"}, "datetime", true)
	created.Default = &schema.RawExpr{X: "current_timestamp()"}
	acc := col("account_id", &schema.IntegerType{T: "int"}, "int(11)", true)
	users.Columns = []*schema.Column{id, name, key, status, full, created, acc}
	users.PrimaryKey = &schema.Index{Parts: []*schema.IndexPart{{C: id}}}
	users.Indexes = []*schema.Index{{Name: "ix_lower", Parts: []*schema.IndexPart{{X: &schema.RawExpr{X: "lcase(`name`)"}}}}}
	users.Attrs = []schema.Attr{&schema.Check{Name: "ck_name", Expr: "char_length(`name`) > 0 and `key` <> 'x\\'y'"}}
	users.ForeignKeys = []*schema.ForeignKey{{Symbol: "fk_account", Table: users, Columns: []*schema.Column{acc},
		RefTable: accounts, RefColumns: []*schema.Column{accID}}}

	orders := &schema.Table{Name: "orders", Schema: app}
	oid := col("id", &schema.IntegerType{T: "int"}, "int(11)", false, &AutoIncrement{})
	ostatus := col("status", &schema.EnumType{T: "enum", Values: []string{"new", "paid"}}, "enum('new','paid')", false)
	uid := col("user_id", &schema.IntegerType{T: "int", Unsigned: true}, "int(10) unsigned", false)
	orders.Columns = []*schema.Column{oid, ostatus, uid}
	orders.PrimaryKey = &schema.Index{Parts: []*schema.IndexPart{{C: oid}}}
	orders.ForeignKeys = []*schema.ForeignKey{{Symbol: "fk_user", Table: orders, Columns: []*schema.Column{uid},
		RefTable: users, RefColumns: []*schema.Column{id}, OnDelete: schema.Cascade}}
	app.Tables = []*schema.Table{users, orders}
	return &schema.Realm{Schemas: []*schema.Schema{app}}
}

func TestImportMariaDBRealm(t *testing.T) {
	res, err := atlasrealm.Import(mariadbRealm(), atlasrealm.Options{})
	if err != nil {
		t.Fatalf("%v\n%s", err, res.DBML)
	}
	for _, want := range []string{
		"database_type: 'MySQL'",
		"Enum users_status_enum {", "Enum orders_status_enum {",
		"status users_status_enum [not null]", "status orders_status_enum [not null]",
		`generated: 'concat(name,'' '',"key")'`,
		"`lcase(name)` [name: 'ix_lower']",
		"`char_length(name) > 0 and \"key\" <> 'x''y'` [name: 'ck_name']",
		"created datetime [default: `current_timestamp()`]",
		"Ref fk_user: users.id <? orders.user_id [delete: cascade]",
	} {
		if !strings.Contains(res.DBML, want) && !strings.Contains(strings.ReplaceAll(res.DBML, "\\'", "''"), want) {
			t.Errorf("missing %q in\n%s", want, res.DBML)
		}
	}
	if strings.Contains(res.DBML, "fk_account") || len(res.Warnings) != 1 || !strings.Contains(res.Warnings[0], "fk_account on app.users left out: it references billing.accounts") {
		t.Errorf("warnings %q\n%s", res.Warnings, res.DBML)
	}
	if g := res.Database.Tables[0].Fields[4].Generated; g == nil || g.Stored {
		t.Errorf("generated %+v", g)
	}
	// exports, translated or not
	for _, d := range export.Dialects {
		if _, err := export.SQLWith(res.Database, d, export.Options{TranslateTypes: true}); err != nil {
			t.Errorf("%s: %v", d, err)
		}
	}
	// the DBML round-trips
	if db, diags := dbml.Parse(res.DBML); db == nil {
		t.Errorf("DBML does not compile: %v", diags)
	}
}

// A PostgreSQL 16-shaped realm: a named enum type shared by two tables,
// two inspected schemas with a foreign key between them, identity, a cast
// in a CHECK, an expression index and a stored generated column.
func TestImportPostgresRealm(t *testing.T) {
	pub := &schema.Schema{Name: "public"}
	sales := &schema.Schema{Name: "sales"}
	mood := &schema.EnumType{T: "mood", Values: []string{"happy", "sad"}, Schema: pub}
	people := &schema.Table{Name: "people", Schema: pub}
	pid := col("id", &schema.IntegerType{T: "bigint"}, "bigint", false, &Identity{})
	email := col("email", &schema.StringType{T: "character varying", Size: 255}, "character varying(255)", false)
	m := col("mood", mood, "mood", true)
	slug := col("slug", &schema.StringType{T: "text"}, "text", true, &schema.GeneratedExpr{Expr: "lower(email)", Type: "STORED"})
	people.Columns = []*schema.Column{pid, email, m, slug}
	people.PrimaryKey = &schema.Index{Parts: []*schema.IndexPart{{C: pid}}}
	people.Indexes = []*schema.Index{{Name: "people_email_lower", Unique: true, Parts: []*schema.IndexPart{{X: &schema.RawExpr{X: "lower((email)::text)"}}}}}
	people.Attrs = []schema.Attr{&schema.Check{Name: "email_nonempty", Expr: "((email)::text <> ''::text)"}}

	deals := &schema.Table{Name: "deals", Schema: sales}
	did := col("id", &schema.IntegerType{T: "integer"}, "integer", false, &Identity{})
	owner := col("owner_id", &schema.IntegerType{T: "bigint"}, "bigint", false)
	dm := col("mood", mood, "mood", true)
	deals.Columns = []*schema.Column{did, owner, dm}
	deals.PrimaryKey = &schema.Index{Parts: []*schema.IndexPart{{C: did}}}
	deals.ForeignKeys = []*schema.ForeignKey{{Symbol: "deals_owner_id_fkey", Table: deals, Columns: []*schema.Column{owner},
		RefTable: people, RefColumns: []*schema.Column{pid}, OnDelete: schema.SetNull}}
	pub.Tables, sales.Tables = []*schema.Table{people}, []*schema.Table{deals}

	res, err := atlasrealm.Import(&schema.Realm{Schemas: []*schema.Schema{pub, sales}}, atlasrealm.Options{Dialect: importer.Postgres})
	if err != nil {
		t.Fatalf("%v\n%s", err, res.DBML)
	}
	if len(res.Database.Enums) != 1 || res.Database.Enums[0].Name != "mood" || len(res.Warnings) != 0 {
		t.Errorf("enums %+v warnings %q\n%s", res.Database.Enums, res.Warnings, res.DBML)
	}
	for _, want := range []string{"database_type: 'PostgreSQL'", "Ref deals_owner_id_fkey: people.id <? sales.deals.owner_id [delete: set null]",
		"`lower((email)::text)` [unique, name: 'people_email_lower']", "generated_stored: 'true'", "id bigint [pk, increment]"} {
		if !strings.Contains(res.DBML, want) {
			t.Errorf("missing %q in\n%s", want, res.DBML)
		}
	}
	if _, err := export.SQL(res.Database, export.Postgres); err != nil {
		t.Error(err)
	}
}
