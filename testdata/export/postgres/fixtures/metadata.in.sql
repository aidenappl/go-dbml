CREATE SCHEMA "sales";

CREATE TABLE "users" (
  "id" int PRIMARY KEY,
  "name" varchar
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY,
  "user_id" int
);

CREATE TABLE "sales"."orders" (
  "id" int PRIMARY KEY
);

COMMENT ON TABLE "users" IS 'this will override';

COMMENT ON TABLE "posts" IS '# Heading 1
    code block
# Heading 2
';
