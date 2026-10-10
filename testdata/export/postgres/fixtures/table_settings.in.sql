CREATE TABLE "user" (
  "id" int PRIMARY KEY,
  "name" string
);

CREATE TABLE "country" (
  "id" int PRIMARY KEY,
  "name" string NOT NULL
);

CREATE TABLE "product" (
  "id" int PRIMARY KEY,
  "name" string,
  "price" decimal NOT NULL
);

CREATE TABLE "merchant" (
  "id" int PRIMARY KEY,
  "user_id" int,
  "product_id" int,
  "address" string
);

COMMENT ON TABLE "country" IS 'name is required';

COMMENT ON TABLE "product" IS 'product must have price';

COMMENT ON TABLE "merchant" IS 'merchants sell a lot';

ALTER TABLE "merchant" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "merchant" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id") DEFERRABLE INITIALLY IMMEDIATE;
