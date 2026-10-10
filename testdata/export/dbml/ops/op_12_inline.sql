Table "b" {
  "id" int [pk]
}

Table "a" {
  "id" int [pk]
  "b_id" int
  "self" int
}

Ref:"a"."b_id" <? "b"."id"

Ref:"a"."self" <? "a"."id"

Records b(id) {
  1
}

Records a(id, b_id, self) {
  1, 1, 1
  2, 9, null
}
