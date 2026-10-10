Table "users" {
  "created_at" timestamp
  "updated_at" timestamp
  "id" int [pk]
  "name" varchar
}

Records users(created_at, updated_at, id, name) {
  '2024-01-01T00:00:00+07:00', '2024-01-01T00:00:00+07:00', 1, 'Alice'
  '2024-01-02T00:00:00+07:00', '2024-01-02T00:00:00+07:00', 2, 'Bob'
}
