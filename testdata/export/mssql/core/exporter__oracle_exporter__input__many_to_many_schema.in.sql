CREATE SCHEMA [C##test]
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [C##test].[users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [C##test].[products] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [users_products] (
  [users_id] int,
  [products_id] int,
  PRIMARY KEY ([users_id], [products_id])
);
GO

ALTER TABLE [users_products] ADD FOREIGN KEY ([users_id]) REFERENCES [users] ([id]);
GO

ALTER TABLE [users_products] ADD FOREIGN KEY ([products_id]) REFERENCES [products] ([id]);
GO


CREATE TABLE [C##test].[users_products(1)] (
  [users_id] int,
  [products_id] int,
  PRIMARY KEY ([users_id], [products_id])
);
GO

ALTER TABLE [C##test].[users_products(1)] ADD FOREIGN KEY ([users_id]) REFERENCES [C##test].[users] ([id]);
GO

ALTER TABLE [C##test].[users_products(1)] ADD FOREIGN KEY ([products_id]) REFERENCES [C##test].[products] ([id]);
GO

