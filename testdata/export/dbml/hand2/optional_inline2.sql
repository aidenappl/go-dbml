Table "x" {
  "id" int [pk]
}

Table "y" {
  "x_id" int
  "x2" int
}

Table "z" {
  "x3" int
}

Ref:"x"."id" ?< "y"."x_id"

Ref:"x"."id" -? "y"."x2"

Ref:"x"."id" ?<> "z"."x3"
