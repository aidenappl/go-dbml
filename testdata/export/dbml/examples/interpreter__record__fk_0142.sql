Table "tags" {
  "id" int [pk]
}

Table "posts" {
  "id" int [pk]
}

Ref:"tags"."id" <>? "posts"."id"

Records tags(id) {
  1
  2
}

Records posts(id) {
  10
}
