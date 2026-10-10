Table "customer" {
  "customer_id" SMALLINT [pk, not null, increment]
  "store_id" TINYINT [not null]
  "first_name" VARCHAR(45) [not null]
  "last_name" VARCHAR(45) [not null, default: false]
  "email" VARCHAR(50)
  "address_id" SMALLINT [not null]
  "active" BOOLEAN [not null, default: true]
  "create_date" DATETIME [not null]
  "last_update" TIMESTAMP [default: `CURRENT_TIMESTAMP`]
}

Table "cities" {
  "id" integer [pk]
  "name" e [default: 'hello']
  "country_id" integer
  Note: 'sasasa'
}

Table "country" {
  "id" integer
  "cities" "string[]"
}

Table "citites" {
  "id" integer
  "name" string

  Indexes {
    (id, name)
  }
}
