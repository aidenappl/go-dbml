Table "users" as U [headerColor: #fff] {
  "id" integer [pk, increment]
  "email" varchar(255) [unique, not null]
  "status" user_status [default: 'active']

  Indexes {
    email [unique, name: "idx_email"]
    (id, email) [pk]
  }
  Note: 'Main users table'
}
