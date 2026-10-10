Table "users" {
  "email" varchar

  Indexes {
    email [type: btree, unique, name: "idx_email"]
  }
}
