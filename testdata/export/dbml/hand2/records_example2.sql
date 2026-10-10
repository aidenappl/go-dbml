Table "orders" {
  "id" int [pk]
  "total" float
  "placed" date
  "status" varchar
}

Records orders(id, total, placed, status) {
  1, 10.5, '2024-03-01', 'paid'
  2, 20, '2024-03-02', null
}
