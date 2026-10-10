Enum "status" {
  "active"
  "churned"
  "inactive"
}

Table "user" {
  "name" text
  "id" int
  "status" status [default: 'active']
  "gender" demographic.gender [default: 'male']
  "age_type" demographic."age segment" [default: 'young_adult']
}

Enum "demographic"."gender" {
  "male"
  "female"
}

Enum "demographic"."age segment" {
  "toddler"
  "children"
  "teenager"
  "young_adult"
  "elder"
}
