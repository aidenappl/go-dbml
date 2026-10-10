Table "a" {
  "id" int [pk]
  "b_id" int [not null]
  "n" int
}

Table "b" {
  "id" int [pk]
  "u" int [unique]
}

Ref:"a"."b_id" <? "b"."id"

Ref "r2":"a"."n" <? "b"."u" [delete: cascade]
