Table "raw" {
  "id" int
  "amount" int
}

Table "stg" {
  "id" int
  "revenue" int
}

Table "fct" {
  "total" int
}

Dep: "raw" -> "stg"

Dep: "raw"."amount" -> "stg"."revenue"

Dep: "stg"."revenue" -> "fct"."total"

Dep {
  "raw" -> "fct"
  "stg" -> "fct"

  note: 'n'
  materialized: 'table'
  priority: 3
}
