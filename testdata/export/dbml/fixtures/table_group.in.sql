Table "users" as U {
  "id" int
  "full_name" varchar
  "created_at" timestamp
  "country_code" int
}

Table "merchants" {
  "id" int
  "merchant_name" varchar
  "country_code" int
  "created at" varchar
  "admin_id" int
}

TableGroup "g1" {
  "users"
  "merchants"
}

Ref:"users"."id" < "merchants"."admin_id"
