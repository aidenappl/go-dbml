Table "a" {
  "id" int
  "id2" int
}

Table "users" [headerColor: #ff0000] {
  "created_at" timestamp [default: `now()`]
  "updated_at" timestamp [default: `now()`]
  "name" varchar
  "email" varchar
  "id" int
  "delete_status" boolean [not null]
  "deleted_at" timestamp [default: `now()`]
  "refColumn1" bit
  "refColumn2" bool

  Indexes {
    email [unique]
    id [pk]
    (email, id) [unique]
    email [unique]
  }
}

Table "testColWithTypeHaveArgs" {
  "id" int
  "callExpColInTest1" int(4) [note: 'simple column with call expression as type']
  "callExpColInTest2" decimal(10,2) [default: 0.1, note: 'column with type that have arguments and default value being a float number']
}

Table "table_no_column_with_partial" {
  "email" varchar [unique]

  Indexes {
    email [unique]
  }
}

Ref:"a"."id2" < "users"."refColumn1"

Ref:"a"."a"."id" < "users"."refColumn2"

Table "a"."a" {
  "id" int
  "id2" int
}
