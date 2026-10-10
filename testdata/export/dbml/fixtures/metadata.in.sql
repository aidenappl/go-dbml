Table "users" {
  "id" int [pk]
  "name" varchar
  Note: 'this will override'
}

Table "posts" {
  "id" int [pk]
  "user_id" int
  Note: '''# Heading 1
    code block
# Heading 2
'''
}

TableGroup "g1" {
  "users"
  "posts"
}

Table "sales"."orders" {
  "id" int [pk]
}

Note overview {
  'a top-level note'
}

Metadata Table "users" {
  owner: 'scott'
  color: '#aaa'
}

Metadata Column "users"."id" {
  pii: 'true'
  masking: 'partial'
}

Metadata TableGroup "g1" {
  team: 'data'
}

Metadata Note overview {
  author: 'docs'
}
