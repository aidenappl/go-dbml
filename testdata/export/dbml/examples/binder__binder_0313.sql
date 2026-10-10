Table "users" {
  "first_name" varchar
  "last_name" varchar
  "email" varchar

  Indexes {
    (first_name, last_name)
    email [unique]
  }
}
