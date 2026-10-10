Table "core"."users" as U {
  "id" int
}

Table "orders" {
  "id" int
}

TableGroup "g1" {
  "orders"
}

TableGroup "g2" {
  "core"."users"
}

Note n1 {
  'x'
}
