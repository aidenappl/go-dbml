CREATE TABLE "users" (
  "id" int,
  "name" varchar(5),
  "code" char(3),
  "title" nvarchar(10),
  "shortcode" nchar(4),
  "description" "character varying(8)"
);

-- Defer constraint checking for INSERT
BEGIN;
SET CONSTRAINTS ALL DEFERRED;

INSERT INTO "users" ("id", "name", "code", "title", "shortcode", "description")
VALUES
  (1, 'Alice', 'ABC', 'Manager', 'MGMT', 'Valid'),
  (2, 'VeryLongName', 'ABCD', 'VeryLongTitle', 'TOOLONG', 'TooLongDesc');

SET CONSTRAINTS ALL IMMEDIATE;
COMMIT;