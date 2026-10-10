CREATE TABLE "products" (
  "id" uuid PRIMARY KEY,
  "kind" text NOT NULL,
  "name" text
);

CREATE UNIQUE INDEX "products_kind_name_key" ON "products" ("kind", "name");
