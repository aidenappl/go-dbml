CREATE TABLE "users" (
  "id" int PRIMARY KEY
);

CREATE TABLE "posts" (
  "id" int PRIMARY KEY
);

CREATE TABLE "user_posts" (
  "user_id" int,
  "post_id" int
);

ALTER TABLE "user_posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "user_posts" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") DEFERRABLE INITIALLY IMMEDIATE;
