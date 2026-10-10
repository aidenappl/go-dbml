Enum "UserRoleEnum" {
  "User"
  "Admin"
}

Table "accounts" {
  "id" int
  "name" "Var Char"
  "role" UserRoleEnum
  "role3" "Schema Name".UserRoleEnum
  "role4" SchemaName.UserRoleEnum
  "role5" "Schema Name"."User Role Enum"
}

Enum "SchemaName"."UserRoleEnum" {
  "User"
  "Admin"
}

Enum "Schema Name"."UserRoleEnum" {
  "User"
  "Admin"
}

Enum "Schema Name"."User Role Enum" {
  "User"
  "Admin"
}
