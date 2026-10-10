CREATE TABLE "products" (
  "price" decimal CHECK (price > 0),
  "discount" decimal,
  CONSTRAINT "disc_nonneg" CHECK (discount >= 0),
  CHECK (discount < price)
);
