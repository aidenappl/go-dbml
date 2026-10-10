Enum "job_status" {
  "created2" [note: 'abcdef']
  "running2"
  "done2"
  "failure2"
}

Enum "gender" {
  "male2"
  "female2"
}

Table "users" {
  "id" int [pk]
  "name" varchar
  "pjs" job_status
  "pjs2" job_status
  "pg" schemaB.gender
  "pg2" gender
}

Table "products" {
  "id" int [pk]
  "name" varchar
}

Table "b" [headerColor: #aaaaaa] {
  "id" int
  "c_id" int
}

Table "c" {
  "id" int
  "b_id" int
}

Table "unary_expression" {
  "id" int(-1) [default: -2]
  "id2" int(--1) [default: -2]
  "id3" int(+-+---+0.1) [default: 7.2225]
}

Table "injected_unary_expression" {
  "id" int
  "id2" int(--1) [default: -2]
  "id3" int(+-+---+0.1) [default: -7.2225]
}

TableGroup "tablegroup_name" {
  "products"
  "users"
  "ecommerce"."users"
  "schemaA"."products"
}

Ref:"ecommerce"."users"."id" < "schemaA"."products"."name"

Ref:"users"."id" < "schemaA"."locations"."name"

Ref:"users"."id" < "ecommerce"."users"."id"

Ref "name_optional":"users"."name" < "ecommerce"."users"."id"

Ref "short_ref":"b"."id" < "c"."id"

Ref "long_ref":"c"."id" < "b"."c_id"

Enum "schemaB"."gender" {
  "male"
  "female"
}

Table "ecommerce"."users" as EU {
  "id" int [pk]
  "name" varchar
  "ejs" job_status
  "ejs2" job_status
  "eg" schemaB.gender
  "eg2" gender
}

Table "schemaA"."products" as A {
  "id" int [pk]
  "name" varchar
}

Table "schemaA"."locations" {
  "id" int [pk]
  "name" varchar
}

Note note {
  'One line note'
}

Note note2 {
  '''# Title
body
'''
}
