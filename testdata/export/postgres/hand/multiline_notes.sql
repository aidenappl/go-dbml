CREATE TABLE "t" (
  "c" int
);

COMMENT ON TABLE "t" IS 'table-level note
';

COMMENT ON COLUMN "t"."c" IS 'it''s escaped 
 newline';
