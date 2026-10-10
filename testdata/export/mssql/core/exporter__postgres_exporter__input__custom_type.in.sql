CREATE SCHEMA [SchemaName]
GO

CREATE SCHEMA [Schema Name]
GO

CREATE TABLE [accounts] (
  [id] int,
  [name] Var Char,
  [role] nvarchar(255) CHECK ([role] IN ('User', 'Admin')),
  [role3] nvarchar(255) CHECK ([role3] IN ('User', 'Admin')),
  [role4] nvarchar(255) CHECK ([role4] IN ('User', 'Admin')),
  [role5] nvarchar(255) CHECK ([role5] IN ('User', 'Admin')),
  [field_keyword] user,
  [field_keyword_quoted] user,
  [field_keyword_dot] user.e,
  [field_keyword_dot_quoted] user.e,
  [field_builtin_capitalized] Int,
  [field_builtin_capitalized_quoted] Int,
  [field_builtin_uppercase] INT,
  [field_builtin_uppercase_quoted] INT,
  [field_custom_with_schema] Random.Rand,
  [field_custom_with_schema_quoted] Random.Rand,
  [field_keyword_capitalized] User,
  [field_keyword_capitalized_quoted] User,
  [field_keyword_uppercase] USER,
  [field_keyword_uppercase_quoted] USER
)
GO
