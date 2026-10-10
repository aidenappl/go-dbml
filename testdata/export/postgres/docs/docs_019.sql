CREATE TABLE "users" (
  "id" int PRIMARY KEY NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "email" varchar UNIQUE,
  "name" varchar,
  "delete_status" boolean NOT NULL,
  "deleted_at" timestamp DEFAULT (now())
);

CREATE UNIQUE INDEX ON "users" ("email");
