Table "merchants" {
  "id" int
  "country_code" varchar

  Indexes {
    (id, country_code) [pk]
  }
}

Table "orders" {
  "id" int [pk]
  "merchant_id" int
  "country" varchar
}

Ref:"merchants".("id", "country_code") < "orders".("merchant_id", "country")

Records merchants(id, country_code) {
  1, 'US'
  2, 'UK'
}

Records orders(id, merchant_id, country) {
  1, 1, 'US'
  2, 1, 'UK'
}
