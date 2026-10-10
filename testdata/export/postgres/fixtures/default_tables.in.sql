CREATE TABLE "orders" (
  "id" int PRIMARY KEY DEFAULT 123,
  "user_id" int UNIQUE NOT NULL,
  "status" varchar DEFAULT 'Completed',
  "created_at" varchar DEFAULT (now())
);

CREATE TABLE "order_items" (
  "order_id" int,
  "product_id" int,
  "quantity" int
);

CREATE TABLE "products" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "merchant_id" int NOT NULL DEFAULT -1,
  "price" float DEFAULT -123.12,
  "stock" boolean DEFAULT true,
  "expiration" date DEFAULT (current_date + interval 1 year)
);
