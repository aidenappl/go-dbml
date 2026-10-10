CREATE TYPE "User Type" AS ENUM (
  'a b',
  'c'
);

CREATE TABLE "Order Items" (
  "Item ID" int PRIMARY KEY,
  "Type" "User Type",
  "user" "USER",
  "qty" Int,
  "select" varchar
);
