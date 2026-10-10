Enum "order_status" {
  "pending"
  "processing"
  "completed"
}

Table "orders" {
  "id" int [pk]
  "status" order_status [default: 'pending']
}
