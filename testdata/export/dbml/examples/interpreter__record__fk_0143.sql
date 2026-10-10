Table "tags" {
  "id" int [pk]
}

Table "posts" {
  "id" int [pk]
}

Ref:"tags"."id" ?<>? "posts"."id"

Records tags(id) {
  null
  1
}

Records posts(id) {
  null
  10
}
