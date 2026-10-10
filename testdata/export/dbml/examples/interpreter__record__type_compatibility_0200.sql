Enum "status" {
  "active"
  "inactive"
}

Table "users" {
  "id" int
  "quantity" int
  "name" varchar(5)
  "price" decimal(5,2)
  "active" boolean
  "status" status
}

Records users(id, quantity, name, price, active, status) {
  `invalid`, 10, 'VeryLongName', 12345.67, `not_bool`, 'invalid_status'
}
