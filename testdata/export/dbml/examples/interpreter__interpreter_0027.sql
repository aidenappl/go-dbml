Table "users" {
  "id" int
}

TableGroup "mixed" {
  "users"
  "ecommerce"."users"
}

Table "ecommerce"."users" {
  "id" int
}
