CREATE TABLE [users] (
  [id] int PRIMARY KEY IDENTITY(1, 1),
  [name] nvarchar(255) NOT NULL,
  [email] nvarchar(255) UNIQUE,
  [status] nvarchar(255) UNIQUE NOT NULL
)
GO
