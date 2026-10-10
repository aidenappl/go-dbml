CREATE TABLE "Bảng Người Dùng" (
  "mã số" int PRIMARY KEY,
  "名前" varchar,
  "emoji 😀" text
);

CREATE TABLE "order items" (
  "user id" int
);

COMMENT ON COLUMN "Bảng Người Dùng"."emoji 😀" IS 'emoji 🎉 inside';

ALTER TABLE "order items" ADD FOREIGN KEY ("user id") REFERENCES "Bảng Người Dùng" ("mã số") DEFERRABLE INITIALLY IMMEDIATE;
