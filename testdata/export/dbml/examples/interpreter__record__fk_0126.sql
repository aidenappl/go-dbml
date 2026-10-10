Table "categories" {
  "id" int [pk]
  "name" varchar
}

Table "products" {
  "id" int [pk]
  "category_id" int
  "name" varchar
}

Ref:"categories"."id" ?< "products"."category_id"

Records categories(id, name) {
  1, 'Electronics'
}

Records products(id, category_id, name) {
  1, 1, 'Laptop'
  2, null, 'Uncategorized Item'
}
