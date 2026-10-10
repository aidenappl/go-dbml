Table "users" {
  "id" int [pk]
}

Table "posts" {
  "user_id" int [not null]
}

Ref:"users"."id" < "posts"."user_id"
