Table "t1" as T [headerColor: #3498DB] {
  "id" integer [unique, pk, not null, increment, default: 1]
  "name" varchar(255) [default: 'x']
  "ratio" float [default: -1.5]
  "flag" bool [default: false]
  "ts" timestamp [default: `now()`]
  "arr" "int[]"
  "json" jsonb

  Indexes {
    id [type: hash, unique, name: "idx_id"]
    (`lower(name)`, name)
    name [pk]
  }
  Note: 'n'
}
