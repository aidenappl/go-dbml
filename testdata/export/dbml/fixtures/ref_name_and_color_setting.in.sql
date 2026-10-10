Table "b" [headerColor: #aaaaaa] {
  "id" int
  "c_id" int
}

Table "c" {
  "id" int
  "b_id" int
}

Ref "short_ref":"b"."id" < "c"."id"

Ref "long_ref":"c"."id" < "b"."c_id"
