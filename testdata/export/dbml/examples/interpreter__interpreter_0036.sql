Table "users" {
  "id" int
}

Table "posts" {
  "user_id" int
}

Ref "fk_user":"users"."id" < "posts"."user_id"
