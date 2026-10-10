Table "users" {
  "id" int
}

Table "posts" {
  "id" int
}

TableGroup "main" {
  "users"
  "posts"
}
