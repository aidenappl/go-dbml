CREATE TABLE "data" (
  "col1" varchar,
  "col2" boolean,
  "col3" varchar
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "data" ("col1", "col2", "col3")
VALUES
  (NULL, TRUE, 'false'),
  (NULL, TRUE, 'false'),
  (NULL, TRUE, 'false'),
  (NULL, TRUE, 'false');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;