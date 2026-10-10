Table "users" {
  "id" int [pk, increment]
  "name" varchar [not null]
  "email" varchar [unique]
  "status" varchar [unique, not null]
}
