Table "users" {
  "id" int [increment]
  "tenant_id" int

  Indexes {
    (id, tenant_id) [pk]
  }
}

Records users(tenant_id) {
  1
  1
}
