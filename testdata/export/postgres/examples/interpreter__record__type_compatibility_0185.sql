CREATE TABLE "messages" (
  "id" int,
  "text" varchar(20),
  "short" varchar(10)
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "messages" ("id", "text", "short")
VALUES
  (1, 'Hello', 'Test'),
  (2, '😀😁😂😃😄', '😀😁😂');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;