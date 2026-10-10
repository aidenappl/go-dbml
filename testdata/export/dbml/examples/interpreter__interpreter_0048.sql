Table "t" {
  "first_name" varchar
  "last_name" varchar

  Indexes {
    (first_name, last_name) [name: "idx_fullname"]
  }
}
