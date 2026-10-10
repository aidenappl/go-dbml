CREATE TABLE "raw_users" (
  "id" int PRIMARY KEY,
  "email" varchar,
  "name" varchar,
  "created_at" timestamp,
  "raw_data" json
);

CREATE TABLE "raw_orders" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "ordered_at" timestamp,
  "status" varchar,
  "raw_data" json
);

CREATE TABLE "raw_order_items" (
  "id" int PRIMARY KEY,
  "order_id" int,
  "product_id" int,
  "quantity" int,
  "price" decimal
);

CREATE TABLE "stg_users" (
  "user_id" int PRIMARY KEY,
  "email" varchar,
  "name" varchar,
  "created_at" timestamp
);

CREATE TABLE "stg_orders" (
  "order_id" int PRIMARY KEY,
  "user_id" int,
  "ordered_at" timestamp,
  "status" varchar
);

CREATE TABLE "stg_order_items" (
  "order_item_id" int PRIMARY KEY,
  "order_id" int,
  "product_id" int,
  "quantity" int,
  "unit_price" decimal,
  "line_total" decimal
);

CREATE TABLE "dim_users" (
  "user_id" int PRIMARY KEY,
  "email" varchar,
  "name" varchar,
  "first_order_at" timestamp,
  "total_orders" int,
  "total_spent" decimal,
  "created_at" timestamp
);

CREATE TABLE "fct_orders" (
  "order_id" int PRIMARY KEY,
  "user_id" int,
  "ordered_at" timestamp,
  "status" varchar,
  "total_items" int,
  "order_total" decimal
);

ALTER TABLE "raw_orders" ADD FOREIGN KEY ("user_id") REFERENCES "raw_users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "raw_order_items" ADD FOREIGN KEY ("order_id") REFERENCES "raw_orders" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "stg_orders" ADD FOREIGN KEY ("user_id") REFERENCES "stg_users" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "stg_order_items" ADD FOREIGN KEY ("order_id") REFERENCES "stg_orders" ("order_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "fct_orders" ADD FOREIGN KEY ("user_id") REFERENCES "dim_users" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;
