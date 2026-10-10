CREATE SCHEMA [C##test01]
GO

CREATE SCHEMA [C##test02]
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [user_id] int NOT NULL
)
GO

CREATE TABLE [C##test01].[users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [C##test01].[products] (
  [id] int PRIMARY KEY,
  [user_id] int NOT NULL
)
GO

CREATE TABLE [C##test02].[users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [C##test02].[products] (
  [id] int PRIMARY KEY,
  [user_id] int NOT NULL
)
GO

ALTER TABLE [products] ADD FOREIGN KEY ([user_id]) REFERENCES [C##test01].[users] ([id])
GO

ALTER TABLE [C##test01].[products] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [C##test01].[products] ADD FOREIGN KEY ([user_id]) REFERENCES [C##test01].[users] ([id])
GO

CREATE TABLE [C##test02].[users_products] (
  [users_id] int,
  [products_id] int,
  PRIMARY KEY ([users_id], [products_id])
);
GO

ALTER TABLE [C##test02].[users_products] ADD FOREIGN KEY ([users_id]) REFERENCES [C##test02].[users] ([id]);
GO

ALTER TABLE [C##test02].[users_products] ADD FOREIGN KEY ([products_id]) REFERENCES [products] ([id]);
GO


CREATE TABLE [users_products(1)] (
  [users_id] int,
  [products_id] int,
  PRIMARY KEY ([users_id], [products_id])
);
GO

ALTER TABLE [users_products(1)] ADD FOREIGN KEY ([users_id]) REFERENCES [users] ([id]);
GO

ALTER TABLE [users_products(1)] ADD FOREIGN KEY ([products_id]) REFERENCES [C##test02].[products] ([id]);
GO


CREATE TABLE [C##test01].[users_products(2)] (
  [users_id] int,
  [products_id] int,
  PRIMARY KEY ([users_id], [products_id])
);
GO

ALTER TABLE [C##test01].[users_products(2)] ADD FOREIGN KEY ([users_id]) REFERENCES [C##test01].[users] ([id]);
GO

ALTER TABLE [C##test01].[users_products(2)] ADD FOREIGN KEY ([products_id]) REFERENCES [C##test02].[products] ([id]);
GO

