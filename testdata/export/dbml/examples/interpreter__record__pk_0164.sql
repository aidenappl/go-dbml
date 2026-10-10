Table "orders" {
  "id" int [default: 0]
  "code" int

  Indexes {
    (id, code) [pk]
  }
}

Records orders(code) {
  1
  1
}
