CREATE TABLE [merchants] (
  [id] int,
  [country_code] nvarchar(255),
  PRIMARY KEY ([id], [country_code])
)
GO

CREATE TABLE [orders] (
  [id] int PRIMARY KEY,
  [merchant_id] int,
  [country] nvarchar(255)
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
  (2, 'UK');
GO
INSERT INTO [orders] ([id], [merchant_id], [country])
VALUES
  (1, 1, 'US'),
  (2, 1, 'UK');
GO

-- Re-enable constraint checks
ALTER TABLE [merchants] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [orders] WITH CHECK CHECK CONSTRAINT ALL;
GO