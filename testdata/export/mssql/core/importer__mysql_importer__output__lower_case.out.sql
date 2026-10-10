CREATE TABLE [orders] (
  [id] int PRIMARY KEY IDENTITY(1, 1),
  [user_id] int UNIQUE NOT NULL,
  [status] nvarchar(255) CHECK ([status] IN ('created', 'running', 'done', 'failure')),
  [created_at] varchar(255)
)
GO

CREATE TABLE [order_items] (
  [order_id] int,
  [product_id] int,
  [quantity] int DEFAULT (1)
)
GO

CREATE TABLE [products] (
  [id] int,
  [name] varchar(255),
  [merchant_id] int NOT NULL,
  [price] int,
  [status] nvarchar(255) CHECK ([status] IN ('Out of Stock', 'In Stock')),
  [created_at] datetime DEFAULT (now()),
  PRIMARY KEY ([id], [price])
)
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [full_name] varchar(255),
  [email] varchar(255) UNIQUE,
  [gender] varchar(255),
  [date_of_birth] varchar(255),
  [created_at] varchar(255),
  [country_code] int
)
GO

CREATE TABLE [merchants] (
  [id] int PRIMARY KEY,
  [merchant_name] varchar(255),
  [country_code] int,
  [created_at] varchar(255),
  [admin_id] int
)
GO

CREATE TABLE [countries] (
  [code] int PRIMARY KEY,
  [name] varchar(255),
  [continent_name] varchar(255)
)
GO

CREATE TABLE [composite_service_item] (
  [composite_service_item_id] int(11) PRIMARY KEY NOT NULL,
  [service_id] int(11) NOT NULL,
  [ticket_item_id] int(11) NOT NULL,
  [discount] decimal(10,0) NOT NULL,
  [original_price] decimal(10,0) NOT NULL,
  [patient_price] decimal(10,0) NOT NULL,
  [insurance_price] decimal(10,0) NOT NULL
)
GO

CREATE TABLE [Countries] (
  [Id] int PRIMARY KEY NOT NULL,
  [Name] varchar(32) NOT NULL
)
GO

CREATE TABLE [States] (
  [Id] int NOT NULL,
  [CountryId] int NOT NULL,
  [Name] varchar(32) NOT NULL,
  PRIMARY KEY ([Id], [CountryId])
)
GO

CREATE TABLE [Cities] (
  [Id] int NOT NULL,
  [CountryId] int NOT NULL,
  [StateId] int NOT NULL,
  [Name] varchar(32) NOT NULL,
  PRIMARY KEY ([Id], [StateId], [CountryId])
)
GO

CREATE INDEX [product_status] ON [products] ("merchant_id", "status")
GO

CREATE UNIQUE INDEX [products_index_1] ON [products] ("id")
GO

CREATE INDEX [ticket_item_id] ON [composite_service_item] ("ticket_item_id")
GO

CREATE INDEX [service_id] ON [composite_service_item] ("service_id")
GO

CREATE INDEX [service_id, ticket_item_id] ON [composite_service_item] ("service_id", "ticket_item_id")
GO

CREATE INDEX [composite_service_item_id] ON [composite_service_item] ("composite_service_item_id")
GO

CREATE INDEX [test test test key] ON [composite_service_item] ("service_id", "patient_price")
GO

CREATE INDEX [IX_States_CountryId] ON [States] ("CountryId")
GO

CREATE INDEX [IX_Cities_CountryId] ON [Cities] ("CountryId")
GO

CREATE INDEX [IX_Cities_StateId_CountryId] ON [Cities] ("StateId", "CountryId")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Notes about products table',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Products price field',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products',
@level2type = N'Column', @level2name = 'price';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'User basic information',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

ALTER TABLE [States] ADD CONSTRAINT [FK_States_Countries_CountryId] FOREIGN KEY ([CountryId]) REFERENCES [Countries] ([Id]) ON DELETE CASCADE
GO

ALTER TABLE [Cities] ADD CONSTRAINT [FK_Cities_Countries_CountryId] FOREIGN KEY ([CountryId]) REFERENCES [Countries] ([Id]) ON DELETE CASCADE
GO

ALTER TABLE [Cities] ADD CONSTRAINT [FK_Cities_States_StateId_CountryId] FOREIGN KEY ([StateId], [CountryId]) REFERENCES [States] ([Id], [CountryId]) ON DELETE CASCADE
GO

ALTER TABLE [order_items] ADD FOREIGN KEY ([order_id]) REFERENCES [orders] ([id])
GO

ALTER TABLE [order_items] ADD FOREIGN KEY ([product_id]) REFERENCES [products] ([id])
GO

ALTER TABLE [users] ADD FOREIGN KEY ([country_code]) REFERENCES [countries] ([code])
GO

ALTER TABLE [merchants] ADD FOREIGN KEY ([country_code]) REFERENCES [countries] ([code])
GO

ALTER TABLE [products] ADD FOREIGN KEY ([merchant_id]) REFERENCES [merchants] ([id])
GO

ALTER TABLE [merchants] ADD FOREIGN KEY ([admin_id]) REFERENCES [users] ([id])
GO
