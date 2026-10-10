CREATE TABLE "data" (
  "active" boolean,
  "verified" bool
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("active", "verified")
VALUES
  (TRUE, FALSE),
  (FALSE, TRUE);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;