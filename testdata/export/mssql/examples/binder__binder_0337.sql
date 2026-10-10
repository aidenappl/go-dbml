CREATE TABLE [merchants] (
  [id] int PRIMARY KEY,
  [country_code] nvarchar(255)
)
GO

CREATE TABLE [orders] (
  [merchant_id] int,
  [country] nvarchar(255)
)
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([merchant_id], [country]) REFERENCES [merchants] ([id], [country_code])
GO
