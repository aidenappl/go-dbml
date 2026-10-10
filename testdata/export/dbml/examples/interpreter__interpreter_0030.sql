Table "users" [headerColor: #3498db] {
  "id" integer [pk, increment]
  "email" varchar(255) [unique, not null]
  "created_at" timestamp
  "updated_at" timestamp
  Note: 'Users table'
}

Table "orders" {
  "id" integer [pk]
  "user_id" integer
  "total" decimal(10,2)
  "created_at" timestamp
  "updated_at" timestamp
}

TableGroup "core" {
  "users"
  "orders"
}

Ref:"users"."id" < "orders"."user_id"
