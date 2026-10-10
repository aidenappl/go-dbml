CREATE TABLE "users" (
  "identity_data" jsonb,
  "email" text
);

CREATE TABLE "json_operators" (
  "id" SERIAL PRIMARY KEY,
  "data" jsonb,
  "field_text" text,
  "field_json" jsonb,
  "nested_text" text,
  "nested_json" jsonb
);

CREATE TABLE "other_operators" (
  "id" SERIAL PRIMARY KEY,
  "a" int,
  "b" int,
  "left_shift" int,
  "right_shift" int,
  "gte" boolean,
  "lte" boolean,
  "neq" boolean,
  "as_text" text
);
