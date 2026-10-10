Table "very_long_users" as U {
  "id" int [pk]
}

Table "orders" {
  "uid" int
  "uid2" int
}

TableGroup "g1" [color: #abc] {
  "very_long_users"
  "orders"
}

Ref:"very_long_users"."id" < "orders"."uid"

Ref:"orders"."uid2" - "very_long_users"."id"

Note sticky {
  'hello'
}
