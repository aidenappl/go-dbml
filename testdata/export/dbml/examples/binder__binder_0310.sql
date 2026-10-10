Table "users" {
  "id" int [pk]
}

Table "posts" {
  "user_id" int
}

Table "comments" {
  "user_id" int
}

Table "likes" {
  "user_id" int
}

Ref:"users"."id" < "posts"."user_id"

Ref:"users"."id" < "comments"."user_id"

Ref:"users"."id" < "likes"."user_id"
