Table "parent" {
  "id" int
}

Table "child" {
  "parent_id" int
}

Ref:"parent"."id" < "child"."parent_id"
