Table "shop"."products" {
  "id" int [pk]
  "category_id" int
}

Table "shop"."orders" {
  "product_id" int
  "category" int
}

Ref:"shop"."products".("id", "category_id") < "shop"."orders".("product_id", "category")
