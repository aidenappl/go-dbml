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
  "status" varchar
}

Ref:"merchants".("id", "country_code") ?< "orders".("merchant_id", "country")

Records merchants(id, country_code) {
  1, 'US'
}

Records orders(id, merchant_id, country, status) {
  1, 1, 'US', 'confirmed'
  2, null, 'UK', 'pending'
  3, 1, null, 'processing'
}
