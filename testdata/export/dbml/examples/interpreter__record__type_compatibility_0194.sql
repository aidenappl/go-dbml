Enum "status" {
  "active"
  "inactive"
  "pending"
}

Table "users" {
  "id" int [pk]
  "status" status
}

Records users(id, status) {
  1, null
  2, 'inactive'
  3, null
  4, null
}
