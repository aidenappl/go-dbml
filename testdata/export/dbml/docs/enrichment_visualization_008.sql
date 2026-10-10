Table "users" {
  "id" int [pk]
  "name" varchar
  Note: 'scott is the owner'
}

TableGroup "g1" {
  "users"
}

Metadata Table "users" {
  owner: 'scott'
}

Metadata Column "users"."id" {
  pii: 'true'
  masking: 'partial'
}
