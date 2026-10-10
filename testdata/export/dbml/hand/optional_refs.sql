Table "a" {
  "id" int [pk]
  "b_id" int
  "c_id" int
  "d_id" int
}

Table "b" {
  "id" int [pk]
}

Table "c" {
  "id" int [pk]
}

Table "d" {
  "id" int [pk]
}

Ref:"b"."id" ?< "a"."b_id"

Ref:"a"."c_id" ?<> "c"."id"

Ref:"a"."d_id" ?-? "d"."id"

Ref:"b"."id" <? "a"."id"
