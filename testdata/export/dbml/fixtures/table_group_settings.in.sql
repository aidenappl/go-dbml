Table "table1" {
  "id" int [pk]
}

TableGroup "group1" [color: #06c] {
  "table1"
  Note: '''note inside setting list of a table group
'''
}
