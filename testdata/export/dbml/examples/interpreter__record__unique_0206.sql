Table "users" {
  "id" int [pk]
  "email" varchar [unique]
  "username" varchar [unique]
}

Table "products" {
  "id" int [pk]
  "sku" varchar [unique]
}

Table "accounts" {
  "id" int [pk]
  "account_number" int [unique]
}

Records users(id, email, username) {
  1, 'alice@example.com', 'alice'
  2, 'bob@example.com', 'bob'
}

Records products(id, sku) {
  1, 'PROD-001'
  2, 'PROD-002'
}

Records accounts(id, account_number) {
  1, 0
  2, -1
  3, 1000
}
