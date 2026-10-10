Table "event_registrations" {
  "event_id" int
  "attendee_id" int
  "registration_date" timestamp

  Indexes {
    (event_id, attendee_id) [unique]
  }
}

Records event_registrations(event_id, attendee_id, registration_date) {
  1, 100, '2024-01-01'
  1, 101, '2024-01-02'
  2, 100, '2024-01-03'
}
