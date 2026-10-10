Table "a" {
  "id" int [pk]
  "b_id" int [not null]
  "n" int
}

Table "b" {
  "id" int [pk]
  "u" int [unique]
}

Ref:"b"."id" ?< "a"."b_id"

Ref "r2":"b"."u" ?< "a"."n" [delete: cascade]
