CREATE TABLE [Order Items] (
  [Item ID] int PRIMARY KEY,
  [Type] nvarchar(255) CHECK ([Type] IN ('a b', 'c')),
  [user] USER,
  [qty] Int,
  [select] nvarchar(255)
)
GO
