Table "users" {
  "id" int
}

Table "orders" {
  "id" int
}

TableGroup "Inventory" {
  "users"
}

TableGroup "Reporting" {
  "orders"
}
