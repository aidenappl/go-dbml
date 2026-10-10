CREATE TABLE "orders" (
  "status" varchar
);

COMMENT ON COLUMN "orders"."status" IS '💸 1 = processing,
✔️ 2 = shipped,
❌ 3 = cancelled,
😔 4 = refunded
';
