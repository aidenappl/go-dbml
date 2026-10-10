Enum "app"."status" {
  "active"
  "on hold"
}

Table "app"."users" {
  "id" int [pk, increment]
  "st" app.status [not null]
  "st2" status
  "cur" public.money_t
  "cur2" other.other.money_t
  "name" varchar [default: 'x', note: '''Bob\'s name''']
  Note: '''It\'s users'''
}

Enum "status" {
  "a"
  "b"
}

Table "posts" {
  "id" bigint [pk, increment]
  "user_id" int
  "t" "timestamp with time zone" [default: `now()`]
  Note: 'posts'
}

Ref:"app"."users"."id" < "posts"."user_id"

Ref:"app"."users"."id" <> "other"."tags"."id"

Ref:"posts"."id" <> "other"."tags"."id"

Ref "r1":"posts"."user_id" - "other"."tags"."id" [update: set null, delete: cascade]

Table "other"."tags" {
  "id" smallint [pk, increment]
  "name" "My Type"
}
