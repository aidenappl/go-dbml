Table "users" {
  "id" int [pk]
}

Table "products" {
  "id" int [pk]
  "user_id" int [not null]
}

Ref:"C##test01"."users"."id" < "products"."user_id"

Ref:"users"."id" < "C##test01"."products"."user_id"

Ref:"C##test01"."users"."id" < "C##test01"."products"."user_id"

Ref:"C##test02"."users"."id" <> "products"."id"

Ref:"users"."id" <> "C##test02"."products"."id"

Ref:"C##test01"."users"."id" <> "C##test02"."products"."id"

Table "C##test01"."users" {
  "id" int [pk]
}

Table "C##test01"."products" {
  "id" int [pk]
  "user_id" int [not null]
}

Table "C##test02"."users" {
  "id" int [pk]
}

Table "C##test02"."products" {
  "id" int [pk]
  "user_id" int [not null]
}
