CREATE SCHEMA [C##dbx]
GO

CREATE TABLE [C##dbx].[users] (
  [user_id] int PRIMARY KEY IDENTITY(1, 1),
  [username] varchar2(40) UNIQUE NOT NULL,
  [password] varchar2(255) NOT NULL
)
GO

CREATE TABLE [C##dbx].[category] (
  [cat_id] int PRIMARY KEY IDENTITY(1, 1),
  [cat_name] nvarchar2(50) NOT NULL,
  [super_cat_id] int,
  [created_date] timestamp DEFAULT (current_timestamp)
)
GO

CREATE TABLE [C##dbx].[product] (
  [product_id] int PRIMARY KEY IDENTITY(1, 1),
  [product_name] nvarchar2(255) NOT NULL,
  [current_price] float NOT NULL,
  [expired_date] timestamp NOT NULL,
  [cat_id] int NOT NULL
)
GO

CREATE TABLE [C##dbx].[wishlist] (
  [user_id] int NOT NULL,
  [product_id] int NOT NULL,
  PRIMARY KEY ([user_id], [product_id])
)
GO

CREATE INDEX [FULLTEXT_INDEX_CATEGORY] ON [C##dbx].[category] ("cat_name")
GO

CREATE INDEX [FULLTEXT_INDEX_PRODUCT] ON [C##dbx].[product] ("product_name")
GO

ALTER TABLE [C##dbx].[category] ADD CONSTRAINT [FK_CATEGORY_SUPER] FOREIGN KEY ([super_cat_id]) REFERENCES [C##dbx].[category] ([cat_id])
GO

ALTER TABLE [C##dbx].[product] ADD CONSTRAINT [FK_PRODUCT_CATEGORY] FOREIGN KEY ([cat_id]) REFERENCES [C##dbx].[category] ([cat_id])
GO

ALTER TABLE [C##dbx].[wishlist] ADD CONSTRAINT [FK_WISHLIST_PRODUCT] FOREIGN KEY ([product_id]) REFERENCES [C##dbx].[product] ([product_id])
GO

ALTER TABLE [C##dbx].[wishlist] ADD CONSTRAINT [FK_WISHLIST_USER] FOREIGN KEY ([user_id]) REFERENCES [C##dbx].[users] ([user_id])
GO
