Table "users" [headerColor: #3498DB] {
  "id" integer [pk]
  "username" varchar(255) [unique, not null]
}

Note note [color: #FF5733] {
  'One line note'
}

Note note2 {
  '''# Title
body
'''
}

Note note3 [color: none] {
  'Note without color'
}
