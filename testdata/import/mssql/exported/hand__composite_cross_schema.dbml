Table "core"."orders" {
  "merchant_id" int
  "country_code" char(2)
}

Table "billing"."merchants" {
  "id" int
  "country_code" char(2)

  Indexes {
    (id, country_code) [pk]
  }
}

Ref "fk_comp":"billing"."merchants".("id", "country_code") < "core"."orders".("merchant_id", "country_code") [update: no action, delete: cascade]

Ref:"billing"."merchants"."id" < "core"."orders"."merchant_id" [delete: set null]
