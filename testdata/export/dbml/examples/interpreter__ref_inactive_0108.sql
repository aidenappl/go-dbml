Table "users" {
  "id" int
}

Table "posts" {
  "user_id" int
}

Ref:"users"."id" < "posts"."user_id" [delete: cascade]
