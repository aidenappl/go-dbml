Table "merchants" {
  "id" int [pk]
  "country_code" varchar
}

Table "orders" {
  "merchant_id" int
  "country" varchar
}

Ref:"merchants".("id", "country_code") < "orders".("merchant_id", "country")
