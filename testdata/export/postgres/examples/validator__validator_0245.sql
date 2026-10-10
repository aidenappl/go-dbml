CREATE TABLE "users" (
  "id" int
);

CREATE TABLE "roles" (
  "id" int
);

CREATE TABLE "user_roles" (
  "user_id" int
);

ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "roles" ("id") DEFERRABLE INITIALLY IMMEDIATE;
