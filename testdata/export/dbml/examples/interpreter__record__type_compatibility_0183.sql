Table "users" {
  "id" int
  "name" varchar(5)
  "code" char(3)
  "title" nvarchar(10)
  "shortcode" nchar(4)
  "description" "character varying(8)"
}

Records users(id, name, code, title, shortcode, description) {
  1, 'Alice', 'ABC', 'Manager', 'MGMT', 'Valid'
  2, 'VeryLongName', 'ABCD', 'VeryLongTitle', 'TOOLONG', 'TooLongDesc'
}
