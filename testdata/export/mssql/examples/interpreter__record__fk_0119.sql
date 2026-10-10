CREATE TABLE [merchants] (
  [id] int,
  [country_code] nvarchar(255),
  PRIMARY KEY ([id], [country_code])
)
GO

CREATE TABLE [orders] (
  [id] int PRIMARY KEY,
  [merchant_id] int,
  [country] nvarchar(255),
  [amount] decimal
)
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([merchant_id], [country]) REFERENCES [merchants] ([id], [country_code])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [merchants] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [orders] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [merchants] ([id], [country_code])
VALUES
  (1, 'US'),
  (1, 'UK'),
  (2, 'US');
GO
INSERT INTO [orders] ([id], [merchant_id], [country], [amount])
VALUES
  (1, 1, 'US', 100),
  (2, 1, 'UK', 200.5),
  (3, 2, 'US', 50);
GO

-- Re-enable constraint checks
ALTER TABLE [merchants] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [orders] WITH CHECK CHECK CONSTRAINT ALL;
GO