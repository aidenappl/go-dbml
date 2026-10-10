Table "T" {
  "a" int
  "b" int
  "x" int [pk]
  "shared" varchar
  "c" int
  "y" int
}

Records T(a, b, x, shared, c, y) {
  1, 2, 3, 'shared_value', 4, 5
}
