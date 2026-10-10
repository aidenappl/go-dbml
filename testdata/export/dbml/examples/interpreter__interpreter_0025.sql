Table "users" {
  "id" int
}

Table "posts" {
  "id" int
}

TableGroup "social" {
  "users"
  "posts"
}
