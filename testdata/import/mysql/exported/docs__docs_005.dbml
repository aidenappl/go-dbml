Table "users" {
  "id" integer [pk]
  "username" varchar(255) [unique, not null]
  "full_name" varchar(255) [not null]
  "gender" varchar(1) [not null]
  "source" varchar(255) [default: 'direct']
  "created_at" timestamp [default: `now()`]
  "rating" integer [default: 10]
}
