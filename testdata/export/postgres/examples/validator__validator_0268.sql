CREATE TABLE "users" (
  "email" varchar
);

CREATE UNIQUE INDEX "idx_email" ON "users" USING BTREE ("email");
