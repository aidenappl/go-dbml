Table "raw_orders" {
  "id" int [pk]
  "user_id" int
  "amount" decimal
  "created_at" timestamp
}

Table "stg_orders" {
  "id" int [pk]
  "user_id" int
  "amount" decimal
}

Table "fct_orders" {
  "id" int [pk]
  "user_id" int
  "revenue" decimal
}

Table "mart_orders" {
  "id" int [pk]
  "total" decimal
}

Dep: "raw_orders" -> "stg_orders"

Dep {
  "stg_orders" -> "fct_orders"

  note: 'Aggregate staging orders into facts'
  materialized: 'table'
  owner: 'data-team'
}

Dep: "stg_orders"."amount" -> "fct_orders"."revenue"

Dep: "raw_orders"."amount" -> "stg_orders"."amount"

Dep: "fct_orders" -> "mart_orders"
