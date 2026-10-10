CREATE TABLE "bookings" (
  "id" integer,
  "country" varchar,
  "booking_date" date,
  "created_at" timestamp,
  PRIMARY KEY ("id", "country")
);

CREATE INDEX "created_at_index" ON "bookings" ("created_at");

CREATE INDEX ON "bookings" ("booking_date");

CREATE UNIQUE INDEX ON "bookings" ("country", "booking_date");

CREATE INDEX ON "bookings" USING HASH ("booking_date");

CREATE INDEX ON "bookings" ((id*2));

CREATE INDEX ON "bookings" ((id*3), (getdate()));

CREATE INDEX ON "bookings" ((id*3), "id");
