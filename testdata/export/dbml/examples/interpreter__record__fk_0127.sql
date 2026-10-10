Table "users" {
  "id" int [pk]
  "name" varchar
}

Table "user_profiles" {
  "id" int [pk]
  "user_id" int [unique]
  "bio" text
}

Table "departments" {
  "id" int [pk]
  "name" varchar
}

Table "employees" {
  "id" int [pk]
  "dept_id" int
  "name" varchar
}

Ref:"user_profiles"."user_id" - "users"."id"

Ref:"departments"."id" < "employees"."dept_id"

Records users(id, name) {
  1, 'Alice'
  2, 'Bob'
}

Records user_profiles(id, user_id, bio) {
  1, 1, 'Alice\'s bio'
  2, 3, 'Invalid user'
}

Records departments(id, name) {
  1, 'Engineering'
}

Records employees(id, dept_id, name) {
  1, 1, 'Alice'
  2, 999, 'Bob with invalid dept'
}
