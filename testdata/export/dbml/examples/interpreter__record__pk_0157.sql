Table "users" {
  "id" int [pk]
  "name" varchar
}

Table "countries" {
  "code" varchar(2) [pk]
  "name" varchar
}

Records users(id, name) {
  1, 'Alice'
  1, 'Bob'
}

Records countries(code, name) {
  'US', 'United States'
  'US', 'USA'
}
