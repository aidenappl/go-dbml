Enum "auth"."role" {
  "admin"
  "user"
}

Table "auth"."users" {
  "id" int [pk]
  "role" auth.role
  "status" app.status
}

Enum "app"."status" {
  "active"
  "inactive"
}

Records auth.users(id, role, status) {
  1, ``, 'active'
  2, ``, ``
}
