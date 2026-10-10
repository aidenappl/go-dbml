CREATE TABLE "user_settings" (
  "user_id" int,
  "category" varchar,
  "value" varchar
);

CREATE UNIQUE INDEX ON "user_settings" ("user_id", "category");

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "user_settings" ("user_id", "category", "value")
VALUES
  (1, NULL, 'default'),
  (1, NULL, 'another default'),
  (1, 'theme', 'dark');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;