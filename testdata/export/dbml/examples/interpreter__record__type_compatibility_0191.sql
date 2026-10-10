Table "users" {
  "id" int [pk, increment]
  "status" varchar [not null, default: 'active']
}

Records users(id, status) {
  null, null
  null, 'inactive'
}
