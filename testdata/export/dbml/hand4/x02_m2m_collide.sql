Table "a" {
  "id" int [pk]
  "b_id" int
}

Table "b" {
  "id" int [pk]
}

Table "a_b" {
  "x" int
}

Table "c" {
  "a_id" int
  "id" int [pk]
}

Ref:"a"."id" <> "b"."id"

Ref:"a"."b_id" <> "b"."id"

Ref:"a".("id", "b_id") <> "c".("a_id", "id")
