Table "users" {
  "id" int [pk]
}

Table "posts" {
  "user_id" int
}

Table "orders" {
  "user_id" int
}

Ref "r1":"users"."id" < "posts"."user_id"

Ref "r2":"users"."id" < "orders"."user_id"
