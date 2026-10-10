Table "father" {
  "obj_id" char(50) [unique, pk]
}

Table "child" {
  "obj_id" char(50) [unique, pk]
  "father_obj_id" char(50)
}

Ref:"father"."obj_id" - "child"."father_obj_id"
