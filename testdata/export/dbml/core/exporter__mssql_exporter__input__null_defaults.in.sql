Table "users" {
  "id" int [pk]
  "name" varchar(100)
  "email" varchar(100)
  "age" int
  "status" varchar(50) [not null, default: 'active']
  "count" int [default: 0]
}
