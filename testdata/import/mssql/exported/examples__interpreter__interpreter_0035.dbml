Table "users" {
  "id" int
}

Table "posts" {
  "user_id" int
}

Ref "post_author":"users"."id" < "posts"."user_id" [update: no action, delete: cascade]
