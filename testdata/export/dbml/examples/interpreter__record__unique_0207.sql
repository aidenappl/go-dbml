Table "users" {
  "id" int [pk]
  "email" varchar [unique]
}

Table "products" {
  "id" int [pk]
  "sku" varchar [unique]
}

Records users(id, email) {
  1, 'alice@example.com'
  2, 'alice@example.com'
}

Records products(id, sku) {
  1, 'PROD-001'
  2, 'PROD-001'
}
