CREATE TABLE "Mixed Case" (
  "Id" INT GENERATED ALWAYS AS IDENTITY,
  "Value" NUMERIC,
  "computed" INT GENERATED ALWAYS AS ("Id" * 2) STORED,
  bits bit varying(8),
  iv interval day to second,
  arr integer ARRAY[4],
  t time(3) with time zone
);
CREATE TABLE child () INHERITS ("Mixed Case");
CREATE TABLE part (id int, region text) PARTITION BY LIST (region);
CREATE TABLE part_eu PARTITION OF part FOR VALUES IN ('eu');
CREATE VIEW v AS SELECT * FROM part;
CREATE MATERIALIZED VIEW mv AS SELECT 1;
ALTER TABLE "Mixed Case" ADD PRIMARY KEY ("Id");
ALTER TABLE part ADD CONSTRAINT part_pk PRIMARY KEY (id, region);
COMMENT ON TABLE part IS NULL;
