CREATE SCHEMA [core]
GO

CREATE SCHEMA [billing]
GO

CREATE TABLE [core].[orders] (
  [merchant_id] int,
  [country_code] char(2)
)
GO

CREATE TABLE [billing].[merchants] (
  [id] int,
  [country_code] char(2),
  PRIMARY KEY ([id], [country_code])
)
GO

ALTER TABLE [core].[orders] ADD CONSTRAINT [fk_comp] FOREIGN KEY ([merchant_id], [country_code]) REFERENCES [billing].[merchants] ([id], [country_code]) ON DELETE CASCADE ON UPDATE NO ACTION
GO

ALTER TABLE [core].[orders] ADD FOREIGN KEY ([merchant_id]) REFERENCES [billing].[merchants] ([id]) ON DELETE SET NULL
GO
