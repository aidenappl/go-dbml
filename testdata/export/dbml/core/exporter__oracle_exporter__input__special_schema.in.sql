Enum "order_status" {
  "created"
  "pending"
  "processing"
  "waiting_for_payment"
  "success"
  "cancel"
}

Table "users" {
  "id" "integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY"
  "username" nvarchar2(255) [not null]
  "fullname" nvarchar2(255) [not null]

  Indexes {
    (`LOWER("fullname")`, username) [name: "IDX_USERS"]
  }
  Note: 'This table has a expression index'
}

Table "orders" {
  "id" "integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY"
  "status" order_status
  "note" nclob [default: '']
  "user_id" integer [not null, note: 'Store user uuid']
  "created_at" timestamp [default: `current_timestamp`]
  "product_id" int
  "supplier_id" int

  Indexes {
    (id, user_id) [type: btree, unique, name: "IDX_ORDER_USER"]
  }
}

Table "order_details" {
  "order_id" INT [unique, not null]
  "unit" nvarchar2(50) [not null]
  "quantity" number [not null]
  "price" float [not null]
}

Table "products" {
  "product_id" int
  "supplier_id" int

  Indexes {
    (product_id, supplier_id) [pk]
  }
}

Ref:"products".("product_id", "supplier_id") < "orders".("product_id", "supplier_id")

Ref:"users"."id" < "orders"."user_id" [delete: cascade]

Ref:"orders"."id" < "order_details"."order_id" [update: restrict, delete: cascade]
