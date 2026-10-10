Table "auth"."users" {
  "id" int [pk]
  "email" varchar
}

Table "auth"."roles" {
  "id" int [pk]
  "name" varchar
}

Table "auth"."sessions" {
  "id" int [pk]
}

TableGroup "auth_core" {
  "auth"."users"
  "auth"."roles"
}
