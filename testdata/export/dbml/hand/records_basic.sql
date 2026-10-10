Enum "user_kind" {
  "admin"
  "guest"
}

Table "users" {
  "id" int [pk, increment]
  "name" varchar
  "score" decimal(5,2)
  "active" bool
  "created" timestamp
  "kind" user_kind
}

Records users(id, name, score, active, created, kind) {
  1, 'Alice', 9.5, true, '2024-01-02T03:04:05', null
  2, 'Bob', -125, false, null, 'guest'
  3, `now()`, 0, true, '2024-01-02', null
}
