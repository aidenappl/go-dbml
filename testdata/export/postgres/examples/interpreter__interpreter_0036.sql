CREATE TABLE "users" (
  "id" int
);

CREATE TABLE "posts" (
  "user_id" int
);

ALTER TABLE "posts" ADD CONSTRAINT "fk_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
