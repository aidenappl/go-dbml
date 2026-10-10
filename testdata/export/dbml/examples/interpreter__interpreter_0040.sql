Table "a" {
  "id" int
}

Table "b" {
  "a_id" int
}

Ref:"a"."id" <> "b"."a_id"
