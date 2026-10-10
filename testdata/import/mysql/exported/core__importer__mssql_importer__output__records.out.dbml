Table "products" {
  "id" integer [pk]
  "name" nvarchar(255)
  "price" decimal(10,2)
  "in_stock" bit
}

Table "orders" {
  "id" integer [pk]
  "product_id" integer
  "quantity" integer
}

Records products(id, name, price, in_stock) {
  1, 'Laptop', 999.99, true
  2, 'Mouse', 29.99, true
  3, 'Keyboard', 79.99, false
  4, 'Monitor', 149.99, true
}

Records orders(id, product_id, quantity) {
  1, 1, 2
  2, 2, 5
  3, 1, 1
}
