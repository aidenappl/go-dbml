Table "A" {
  "id" integer
  "code" number
}

Table "B" {
  "id" integer
  "code" number
}

Ref:"B"."id" < "A"."id" [update: cascade, delete: no action]

Ref:"B"."code" < "A"."code" [update: set null, delete: set default]
