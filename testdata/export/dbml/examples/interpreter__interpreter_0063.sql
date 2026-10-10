Table "users" {
  "id" int
}

Table "posts" {
  "id" int
}

TableGroup "auth_tables" {
  "users"
}

TableGroup "content_tables" {
  "posts"
}
