Table "a" {
  "id" int
}

TableGroup "grp" [color: #ff0000] {
  "a"
  "s"."b"
  Note: 'inner note'
}

Table "s"."b" {
  "id" int
}
