CREATE TABLE [user] (
  [id] int PRIMARY KEY,
  [name] string
)
GO

CREATE TABLE [country] (
  [id] int PRIMARY KEY,
  [name] string NOT NULL
)
GO

CREATE TABLE [product] (
  [id] int PRIMARY KEY,
  [name] string,
  [price] decimal NOT NULL
)
GO

CREATE TABLE [merchant] (
  [id] int PRIMARY KEY,
  [user_id] int,
  [product_id] int,
  [address] string
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'name is required',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'country';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'product must have price',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'product';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'merchants sell a lot',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'merchant';
GO

ALTER TABLE [merchant] ADD FOREIGN KEY ([user_id]) REFERENCES [user] ([id])
GO

ALTER TABLE [merchant] ADD FOREIGN KEY ([product_id]) REFERENCES [product] ([id])
GO
