Table "products" {
  "id" int
  "region" varchar

  Indexes {
    (id, region) [pk]
  }
}

Table "categories" {
  "id" int
  "region" varchar

  Indexes {
    (id, region) [pk]
  }
}

Ref:"products".("id", "region") <> "categories".("id", "region")

Records products(id, region) {
  1, 'US'
  2, 'US'
}

Records categories(id, region) {
  1, 'US'
  3, 'EU'
}
