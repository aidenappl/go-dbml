Enum "order_status" {
  "pending" [note: 'Initial state']
  "processing"
  "completed"
  "cancelled"
}

Table "users" [headerColor: #3498db] {
  "id" integer [pk, increment]
  "email" varchar(255) [unique, not null, note: 'User email']
  "status" varchar [default: 'active']
  "created_at" timestamp [default: `now()`]
  "updated_at" timestamp

  Indexes {
    email [unique, name: "idx_users_email"]
  }
  Note: 'Main users table'
}

Table "orders" {
  "id" integer [pk, increment]
  "user_id" integer
  "status" order_status [default: 'pending']
  "total" decimal(10,2) [not null]
  "created_at" timestamp [default: `now()`]
  "updated_at" timestamp
}

TableGroup "ecommerce" [color: #2ecc71] {
  "users"
  "orders"
  Note: 'Core e-commerce tables'
}

Ref:"users"."id" < "orders"."user_id"
