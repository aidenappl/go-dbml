Table "a" {
  "id" int
}

Table "b" {
  "id" int
}

Dep: "b"."id" -> "a"."id"

Dep: "b" -> "a"
