Enum "status" {
  "active" [note: 'Active']
}

Table "users" {
  "email" varchar
}

Metadata Table "users" {
  owner: 'bob'
  pii: 'yes'
  accent: '#ff00ff'
  domain: 'identity'
}

Metadata Column "users"."email" {
  sensitivity: 'high'
}
