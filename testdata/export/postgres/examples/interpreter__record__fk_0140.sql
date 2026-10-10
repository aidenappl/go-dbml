CREATE TABLE "students" (
  "id" int PRIMARY KEY
);

CREATE TABLE "courses" (
  "id" int PRIMARY KEY
);

CREATE TABLE "students_courses" (
  "students_id" int,
  "courses_id" int,
  PRIMARY KEY ("students_id", "courses_id")
);

ALTER TABLE "students_courses" ADD FOREIGN KEY ("students_id") REFERENCES "students" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "students_courses" ADD FOREIGN KEY ("courses_id") REFERENCES "courses" ("id") DEFERRABLE INITIALLY IMMEDIATE;


-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "students" ("id")
VALUES
  (1),
  (2),
  (3);
INSERT INTO "courses" ("id")
VALUES
  (10),
  (20);

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;