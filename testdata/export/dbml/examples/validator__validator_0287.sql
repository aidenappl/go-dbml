Table "a" {
  "id" int
  "b_id" int
}

Table "b" {
  "id" int
  "a_id" int
}

Ref:"b"."id" < "a"."b_id"

Ref:"a"."id" < "b"."a_id"
