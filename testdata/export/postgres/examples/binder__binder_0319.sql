CREATE TYPE "order_status" AS ENUM (
  'pending',
  'processing',
  'completed'
);

CREATE TABLE "orders" (
  "id" int PRIMARY KEY,
  "status" order_status DEFAULT 'pending'
);
