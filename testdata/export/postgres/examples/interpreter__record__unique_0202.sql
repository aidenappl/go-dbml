CREATE TABLE "user_profiles" (
  "user_id" int,
  "profile_type" varchar,
  "data" text
);

CREATE UNIQUE INDEX ON "user_profiles" ("user_id", "profile_type");

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "user_profiles" ("user_id", "profile_type", "data")
VALUES
  (1, 'work', 'Software Engineer'),
  (1, 'personal', 'Loves hiking'),
  (2, 'work', 'Designer');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;