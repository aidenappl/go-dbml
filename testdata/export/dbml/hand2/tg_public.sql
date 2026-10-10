Table "users" {
  "id" int
}

Table "users2" {
  "id" int
}

TableGroup "g" {
  "users"
  "users2"
}

Ref:"users2"."id" < "users"."id"
