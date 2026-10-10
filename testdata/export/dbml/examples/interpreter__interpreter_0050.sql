Table "t" {
  "first_name" varchar
  "last_name" varchar

  Indexes {
    (`lower(first_name)`, last_name)
  }
}
