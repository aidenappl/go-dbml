CREATE TABLE "follows" (
  "following_user_id" integer,
  "followed_user_id" integer,
  "created_at" timestamp
);

CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "role" varchar,
  "created_at" timestamp
);

CREATE TABLE "posts" (
  "id" integer PRIMARY KEY,
  "title" varchar,
  "body" text,
  "user_id" integer,
  "status" varchar,
  "created_at" timestamp
);

COMMENT ON TABLE "follows" IS '# Heading 1
    code block
# Heading 2
  * 1
  * 1
';

COMMENT ON TABLE "users" IS '# Heading 1
    code block
# Heading 2
  * 1
  * 2
';

COMMENT ON TABLE "posts" IS '# Heading 1
    code block
# Heading 2
  * 1
    * 1.1
    * 1.2
  * 2
';

COMMENT ON COLUMN "posts"."body" IS '# Heading 1
    code block
# Heading 2
  * 1
  * 2
  ';

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "follows" ADD FOREIGN KEY ("following_user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "follows" ADD FOREIGN KEY ("followed_user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
