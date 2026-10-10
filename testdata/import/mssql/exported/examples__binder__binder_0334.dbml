Table "users" {
  "id" int [pk]
}

Table "posts" {
  "id" int [pk]
}

Table "user_posts" {
  "user_id" int
  "post_id" int
}

Ref:"users"."id" < "user_posts"."user_id"

Ref:"posts"."id" < "user_posts"."post_id"
