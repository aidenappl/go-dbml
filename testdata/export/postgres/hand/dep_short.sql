CREATE TABLE "raw_orders" (
  "id" int,
  "amount" int
);

CREATE TABLE "stg_orders" (
  "id" int,
  "revenue" int
);

CREATE TABLE "fct_orders" (
  "revenue" int
);
