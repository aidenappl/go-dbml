Table "a" {
  "id" int [pk]
  "code" int [unique]

  Indexes {
    (id, code) [unique]
  }
}

Table "b" {
  "a_id" int [not null]
  "a_code" int
}

Ref:"a".("id", "code") < "b".("a_id", "a_code")
