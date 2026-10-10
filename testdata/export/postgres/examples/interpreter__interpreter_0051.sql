CREATE TABLE "t" (
  "email" varchar
);

CREATE UNIQUE INDEX "idx_email" ON "t" USING BTREE ("email");
