Table "bookings" {
  "id" integer
  "country" varchar
  "booking_date" date
  "created_at" timestamp

  Indexes {
    (id, country) [pk]
    created_at [name: "created_at_index"]
    booking_date
    (country, booking_date) [unique]
    booking_date [type: hash]
    `id*2`
    (`id*3`, `getdate()`)
    (`id*3`, id)
  }
}
