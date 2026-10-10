CREATE TABLE "orders" (
  "id" int DEFAULT 0,
  "code" int,
  PRIMARY KEY ("id", "code")
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "orders" ("code")
VALUES
  (NULL),
  (1);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;