Enum "products_status" {
  "out_of_stock"
  "in_stock" [note: 'In stock']
  "running_low" [note: 'less than 20']
}

Enum "products_status2" {
  "out_of_stock"
  "in_stock"
  "running_low" [note: 'less than 20']
}

Table "orders" {
  "id" int [pk]
  "user_id" int [unique, not null]
  "status" varchar [note: 'Status of an order']
  "created_at" varchar [note: 'When order created']

  Indexes {
    id [type: hash]
    created_at
  }
}
