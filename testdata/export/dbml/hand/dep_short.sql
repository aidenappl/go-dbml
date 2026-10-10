Table "raw_orders" {
  "id" int
  "amount" int
}

Table "stg_orders" {
  "id" int
  "revenue" int
}

Table "fct_orders" {
  "revenue" int
}

Dep: "raw_orders" -> "stg_orders"

Dep: "stg_orders" -> "fct_orders"

Dep: "raw_orders"."amount" -> "stg_orders"."revenue"

Dep: "stg_orders" -> "fct_orders"."revenue"
