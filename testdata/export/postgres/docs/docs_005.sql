CREATE TABLE "users" (
  "id" integer PRIMARY KEY,
  "username" varchar(255) UNIQUE NOT NULL,
  "full_name" varchar(255) NOT NULL,
  "gender" varchar(1) NOT NULL,
  "source" varchar(255) DEFAULT 'direct',
  "created_at" timestamp DEFAULT (now()),
  "rating" integer DEFAULT 10
);
