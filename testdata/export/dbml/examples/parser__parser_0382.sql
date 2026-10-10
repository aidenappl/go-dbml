Table "users" {
  "id" int [pk]
  "email" varchar

  Indexes {
    email [unique]
  }
}
