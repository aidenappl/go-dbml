Table "users" {
  "id" int
}

Table "orders" {
  "id" int
}

Table "events" {
  "id" int
}

Table "reports" {
  "id" int
}

Dep: "users" -> "orders"

Dep: "orders" -> "events"

Dep: "users" -> "reports"

Dep: "events" -> "users"
