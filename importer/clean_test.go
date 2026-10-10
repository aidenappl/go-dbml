package importer_test

import (
	"strings"
	"testing"

	"github.com/aidenappl/go-dbml/emit"
	"github.com/aidenappl/go-dbml/export"
	"github.com/aidenappl/go-dbml/importer"
	"github.com/aidenappl/go-dbml/model"
)

func cleanImport(t *testing.T, src string, f importer.Format) (*model.Database, string) {
	t.Helper()
	r, err := importer.ImportWith(src, f, importer.ImportOptions{IncludeRecords: true})
	if err != nil {
		t.Fatal(err)
	}
	if r.Database == nil {
		t.Fatalf("does not compile: %v\n%s", r.Diagnostics, r.DBML)
	}
	return r.Database, r.DBML
}

func table(t *testing.T, db *model.Database, name string) *model.Table {
	t.Helper()
	for _, tb := range db.Tables {
		if tb.Name == name {
			return tb
		}
	}
	t.Fatalf("no table %q", name)
	return nil
}

func column(t *testing.T, tb *model.Table, name string) *model.Column {
	t.Helper()
	for _, c := range tb.Fields {
		if c.Name == name {
			return c
		}
	}
	t.Fatalf("no column %q in %s", name, tb.Name)
	return nil
}

func TestCleanMySQLDump(t *testing.T) {
	src := "-- MySQL dump 10.13  Distrib 8.0.36\n" +
		"/*!40101 SET NAMES utf8mb4 */;\n" +
		"DROP TABLE IF EXISTS `order items`;\n" +
		"CREATE TABLE `order items` (\n" +
		"  `id` bigint unsigned NOT NULL AUTO_INCREMENT,\n" +
		"  `qty` int NOT NULL DEFAULT '1',\n" +
		"  `sku` varchar(32) COLLATE utf8mb4_bin NOT NULL COMMENT 'it''s the \\\\ SKU',\n" +
		"  `label` varchar(64) DEFAULT 'O\\'Brien',\n" +
		"  `we``ird` int DEFAULT NULL,\n" +
		"  PRIMARY KEY (`id`),\n" +
		"  KEY `idx_lower_sku` ((lower(`sku`))),\n" +
		"  CONSTRAINT `qty_positive` CHECK ((`qty` > 0))\n" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Line ''items''';\n" +
		"INSERT INTO `order items` VALUES (1,2,'a\\'b','x',NULL);\n"
	db, text := cleanImport(t, src, importer.MySQL)
	tb := table(t, db, "order items")
	if tb.Note == nil || tb.Note.Value != "Line 'items'" {
		t.Errorf("table note %+v", tb.Note)
	}
	if n := column(t, tb, "sku").Note; n == nil || n.Value != `it's the \ SKU` {
		t.Errorf("column note %+v", n)
	}
	if d := column(t, tb, "label").DBDefault; d == nil || d.Value != "O'Brien" {
		t.Errorf("default %+v", d)
	}
	column(t, tb, "we`ird")
	if len(tb.Checks) != 1 || tb.Checks[0].Expression != "(qty > 0)" {
		t.Errorf("checks %+v", tb.Checks)
	}
	found := false
	for _, ix := range tb.Indexes {
		for _, c := range ix.Columns {
			if c.Type == "expression" && strings.Contains(c.Value, "lower(sku)") {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expression index missing:\n%s", text)
	}
	if v := db.Records[0].Values[0][2].Value; v != "a'b" {
		t.Errorf("record value %v", v)
	}
}

func TestCleanPgDump(t *testing.T) {
	src := `CREATE TABLE public.notes (
    id integer NOT NULL,
    body text DEFAULT 'it''s' NOT NULL,
    "Odd ""Name""" text,
    tag text DEFAULT E'a\nb'
);
COMMENT ON TABLE public.notes IS 'Notes, it''s true';
COMMENT ON COLUMN public.notes.body IS 'the ''body''';
`
	db, _ := cleanImport(t, src, importer.Postgres)
	tb := table(t, db, "notes")
	if d := column(t, tb, "body").DBDefault; d == nil || d.Value != "it's" {
		t.Errorf("default %+v", d)
	}
	if d := column(t, tb, "tag").DBDefault; d == nil || d.Value != "a\nb" {
		t.Errorf("E'' default %+v", d)
	}
	column(t, tb, `Odd "Name"`)
	if tb.Note == nil || tb.Note.Value != "Notes, it's true" {
		t.Errorf("note %+v", tb.Note)
	}
	if n := column(t, tb, "body").Note; n == nil || n.Value != "the 'body'" {
		t.Errorf("column note %+v", n)
	}
}

func TestCleanMSSQLExtendedProperty(t *testing.T) {
	src := "CREATE TABLE [dbo].[Customers](\n\t[Id] [int] IDENTITY(1,1) NOT NULL,\n\t[Na]]me] [nvarchar](50) NULL DEFAULT (N'O''Neil'),\n" +
		"\tCONSTRAINT [CK_Id] CHECK ([Id] > 0)\n)\nGO\n" +
		"EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'Customer master, it''s' , @level0type=N'SCHEMA',@level0name=N'dbo', @level1type=N'TABLE',@level1name=N'Customers'\nGO\n" +
		"EXEC sp_addextendedproperty N'MS_Description', N'The name', N'SCHEMA', N'dbo', N'TABLE', N'Customers', N'COLUMN', N'Na]me'\nGO\n" +
		"ALTER TABLE [dbo].[Customers] NOCHECK CONSTRAINT ALL\nGO\n"
	db, _ := cleanImport(t, src, importer.MSSQL)
	tb := table(t, db, "Customers")
	if tb.Note == nil || tb.Note.Value != "Customer master, it's" {
		t.Errorf("table note %+v", tb.Note)
	}
	c := column(t, tb, "Na]me")
	if c.Note == nil || c.Note.Value != "The name" {
		t.Errorf("positional column note %+v", c.Note)
	}
	if c.DBDefault == nil || c.DBDefault.Value != "O'Neil" {
		t.Errorf("default %+v", c.DBDefault)
	}
	if len(tb.Checks) != 1 || tb.Checks[0].Expression != "Id > 0" {
		t.Errorf("checks %+v", tb.Checks)
	}
}

func TestCleanSchemaRB(t *testing.T) {
	src := `ActiveRecord::Schema[7.1].define(version: 2024_05_01_120000) do
  create_table "users", comment: "People", force: :cascade do |t|
    t.string "email", limit: 320, null: false, default: "", comment: "login"
    t.decimal "balance", precision: 10, scale: 2, default: 0.0
    t.boolean "admin", default: false
    t.timestamps
    t.index ["email"], name: "index_users_on_email", unique: true
  end
  create_table "posts", id: :uuid do |t|
    t.references :user, null: false
    # a comment inside the block
    t.datetime "published_at", default: -> { "CURRENT_TIMESTAMP" }
  end
  create_table "tags", id: false do |t|
    t.string "name"
  end
  add_foreign_key "posts", "users", on_delete: :cascade, name: "fk_posts_users"
end
`
	db, text := cleanImport(t, src, importer.SchemaRB)
	users := table(t, db, "users")
	if id := column(t, users, "id"); id.Type.TypeName != "bigint" || id.PK == nil || !*id.PK {
		t.Errorf("implicit id %+v", id)
	}
	if users.Note == nil || users.Note.Value != "People" {
		t.Errorf("table comment %+v", users.Note)
	}
	email := column(t, users, "email")
	if email.Type.TypeName != "varchar(320)" || email.NotNull == nil || !*email.NotNull || email.Note == nil || email.Note.Value != "login" {
		t.Errorf("email %+v", email)
	}
	if b := column(t, users, "balance"); b.Type.TypeName != "decimal(10,2)" {
		t.Errorf("balance type %q", b.Type.TypeName)
	}
	column(t, users, "created_at")
	if len(users.Indexes) != 1 || users.Indexes[0].Unique == nil {
		t.Errorf("indexes %+v", users.Indexes)
	}
	posts := table(t, db, "posts")
	if column(t, posts, "id").Type.TypeName != "uuid" {
		t.Errorf("uuid id")
	}
	if d := column(t, posts, "published_at").DBDefault; d == nil || d.Type != "expression" {
		t.Errorf("lambda default %+v", d)
	}
	for _, f := range table(t, db, "tags").Fields {
		if f.Name == "id" {
			t.Errorf("id: false still has id")
		}
	}
	if len(db.Refs) != 1 || db.Refs[0].Name == nil || *db.Refs[0].Name != "fk_posts_users" {
		t.Errorf("refs %+v\n%s", db.Refs, text)
	}
	if s, err := emit.DBML(db); err != nil || !strings.Contains(s, "Ref fk_posts_users: users.id <? posts.user_id [delete: cascade]") {
		t.Errorf("emit %v:\n%s", err, s)
	}
	// pre-5.1 schemas get integer ids
	old, _ := cleanImport(t, "ActiveRecord::Schema.define(version: 20160101000000) do\n  create_table \"a\" do |t|\n    t.string \"n\"\n  end\nend\n", importer.SchemaRB)
	if column(t, table(t, old, "a"), "id").Type.TypeName != "integer" {
		t.Errorf("old id type")
	}
}

func TestRecordsOptIn(t *testing.T) {
	src := "CREATE TABLE t (a int);\nINSERT INTO t VALUES (1),(2),(3);\n"
	r, err := importer.ImportWith(src, importer.MySQL, importer.ImportOptions{})
	if err != nil || len(r.Database.Records) != 0 {
		t.Fatalf("records must be skipped by default: %v %+v", err, r.Database.Records)
	}
	r, _ = importer.ImportWith(src, importer.MySQL, importer.ImportOptions{IncludeRecords: true, MaxRecordRows: 2})
	if n := len(r.Database.Records[0].Values); n != 2 || !r.RecordsTruncated {
		t.Fatalf("cap: %d rows, truncated=%v", n, r.RecordsTruncated)
	}
	if _, err := importer.ImportWith(strings.Repeat(" ", importer.MaxInputBytes+1), importer.MySQL, importer.ImportOptions{}); err != importer.ErrInputTooLarge {
		t.Fatalf("want ErrInputTooLarge, got %v", err)
	}
}

// Rails abstract types become the adapter's SQL types; create_enum,
// t.enum, t.virtual, t.references ... foreign_key: {to_table:} are kept.
func TestCleanSchemaRBFidelity(t *testing.T) {
	src := `ActiveRecord::Schema[7.1].define(version: 2024_05_01_000000) do
  create_enum "mood", ["it's ok", "sad"]
  create_table "accts", comment: "it's" do |t|
    t.string "name", default: "O'Brien", null: false
    t.enum "mood", enum_type: "mood", default: "sad"
    t.virtual "full", type: :string, as: "name || 'x'", stored: true
    t.references "parent", foreign_key: { to_table: :accts, on_delete: :nullify }
    t.boolean "ok", default: true
    t.float "f"
    t.jsonb "props", default: {}
    t.integer "small", limit: 2
    t.timestamps
  end
end
`
	db, text := cleanImport(t, src, importer.SchemaRB)
	a := table(t, db, "accts")
	want := map[string]string{"id": "bigint", "name": "varchar", "mood": "mood", "full": "varchar", "parent_id": "bigint", "ok": "boolean",
		"f": "double precision", "props": "jsonb", "small": "smallint", "created_at": "timestamp"}
	for c, typ := range want {
		if got := column(t, a, c).Type.TypeName; got != typ {
			t.Errorf("%s: type %q, want %q", c, got, typ)
		}
	}
	if g := column(t, a, "full").Generated; g == nil || g.Expression != "name || 'x'" || !g.Stored {
		t.Errorf("virtual: %+v", g)
	}
	if len(db.Enums) != 1 || len(db.Enums[0].Values) != 2 {
		t.Errorf("enums %+v", db.Enums)
	}
	if len(db.Refs) != 1 || db.Refs[0].OnDelete == nil || *db.Refs[0].OnDelete != "set null" || db.Refs[0].Endpoints[1].TableName != "accts" {
		t.Errorf("refs %+v\n%s", db.Refs, text)
	}
	if len(a.Indexes) != 1 || a.Indexes[0].Columns[0].Value != "parent_id" {
		t.Errorf("reference index %+v", a.Indexes)
	}
	// MySQL schemas get MySQL types
	my, _ := cleanImport(t, "ActiveRecord::Schema[7.1].define(version: 2024_01_01_000000) do\n  create_table \"a\", charset: \"utf8mb4\" do |t|\n    t.string \"s\"\n    t.integer \"i\", limit: 1, unsigned: true\n    t.boolean \"b\"\n    t.datetime \"d\", precision: 6\n  end\nend\n", importer.SchemaRB)
	at := table(t, my, "a")
	for c, typ := range map[string]string{"s": "varchar(255)", "i": "tinyint unsigned", "b": "tinyint(1)", "d": "datetime(6)"} {
		if got := column(t, at, c).Type.TypeName; got != typ {
			t.Errorf("mysql %s: type %q, want %q", c, got, typ)
		}
	}
}

// MySQL clean mode keeps UNSIGNED, FULLTEXT, prefix key parts and applies
// ALTER TABLE ADD/MODIFY/CHANGE/DROP.
func TestCleanMySQLFidelity(t *testing.T) {
	src := "CREATE TABLE `a` (\n  `id` int unsigned NOT NULL AUTO_INCREMENT,\n  `title` varchar(255) NOT NULL,\n  `body` text,\n  PRIMARY KEY (`id`),\n  KEY `ix_title` (`title`(20) DESC),\n  FULLTEXT KEY `ft` (`title`,`body`)\n);\n" +
		"ALTER TABLE `a` ADD COLUMN `n` int DEFAULT NULL AFTER `id`, ADD INDEX `ix_n` (`n`);\n" +
		"ALTER TABLE `a` MODIFY `title` varchar(400) NOT NULL;\n" +
		"ALTER TABLE `a` CHANGE `body` `content` mediumtext;\n" +
		"ALTER TABLE `a` ADD COLUMN `gone` int, DROP COLUMN `gone`;\n"
	db, text := cleanImport(t, src, importer.MySQL)
	a := table(t, db, "a")
	var names []string
	for _, f := range a.Fields {
		names = append(names, f.Name+":"+f.Type.TypeName)
	}
	if got := strings.Join(names, ","); got != "id:int unsigned,n:int,title:varchar(400),content:mediumtext" {
		t.Errorf("columns %s\n%s", got, text)
	}
	byName := map[string]*model.Index{}
	for _, ix := range a.Indexes {
		if ix.Name != nil {
			byName[*ix.Name] = ix
		}
	}
	if ix := byName["ix_title"]; ix == nil || ix.Columns[0].Value != "title(20)" {
		t.Errorf("prefix index %+v", ix)
	}
	if ix := byName["ft"]; ix == nil || ix.Type == nil || *ix.Type != "fulltext" {
		t.Errorf("fulltext %+v", ix)
	}
	if byName["ix_n"] == nil {
		t.Errorf("ALTER ADD INDEX lost")
	}
	sql, err := export.SQL(db, export.MySQL)
	if err != nil || !strings.Contains(sql, "CREATE FULLTEXT INDEX `ft` ON `a` (`title`, `content`)") && !strings.Contains(sql, "CREATE FULLTEXT INDEX `ft`") || !strings.Contains(sql, "(`title`(20))") {
		t.Errorf("export %v\n%s", err, sql)
	}
}
