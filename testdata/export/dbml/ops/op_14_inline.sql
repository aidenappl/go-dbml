Table "b" {
  "id" int [pk]
}

Table "a" {
  "id" int [pk]
  "b_id" int
  "self" int
}

Ref:"b"."id" <>? "a"."b_id"

Ref:"a"."id" <>? "a"."self"

Records b(id) {
  1
}

Records a(id, b_id, self) {
  1, 1, 1
  2, 9, null
}
