Table "users" {
  "id" int [pk]
  "name" varchar
}

Table "posts" {
  "id" int [pk]
  "title" varchar
  "user_id" int
}

Ref:"users"."id" < "posts"."user_id"

Records users(id, name) {
  1, 'Alice'
  2, 'Bob'
}

Records posts(id, title, user_id) {
  1, 'Post 1', 1
  2, 'Post 2', 2
}
