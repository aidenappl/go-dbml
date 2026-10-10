Table "students" {
  "id" int [pk]
}

Table "courses" {
  "id" int [pk]
}

Ref:"students"."id" <> "courses"."id"

Records students(id) {
  1
  2
  3
}

Records courses(id) {
  10
  20
}
