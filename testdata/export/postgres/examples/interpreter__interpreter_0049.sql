CREATE TABLE "t" (
  "email" varchar
);

CREATE INDEX "idx_email_lower" ON "t" ((lower(email)));
