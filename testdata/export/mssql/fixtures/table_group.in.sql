CREATE TABLE [users] (
  [id] int,
  [full_name] nvarchar(255),
  [created_at] timestamp,
  [country_code] int
)
GO

CREATE TABLE [merchants] (
  [id] int,
  [merchant_name] nvarchar(255),
  [country_code] int,
  [created at] nvarchar(255),
  [admin_id] int
)
GO

ALTER TABLE [merchants] ADD FOREIGN KEY ([admin_id]) REFERENCES [users] ([id])
GO
