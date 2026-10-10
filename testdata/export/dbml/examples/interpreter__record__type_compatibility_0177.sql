Table "data" {
  "id" int
  "active" boolean
}

Records data(id, active) {
  1, `not_a_bool`
  2, `99`
  3, `-1`
  4, `invalid`
  5, `2`
}
