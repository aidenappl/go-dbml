CREATE TABLE [products] (
  [id] int,
  [region] nvarchar(255),
  PRIMARY KEY ([id], [region])
)
GO

CREATE TABLE [orders] (
  [id] int PRIMARY KEY,
  [product_id] int,
  [region] nvarchar(255)
)
GO

ALTER TABLE [orders] ADD FOREIGN KEY ([product_id], [region]) REFERENCES [products] ([id], [region])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [orders] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [products] ([id], [region])
VALUES
  (1, 'US'),
  (NULL, 'EU');
GO
INSERT INTO [orders] ([id], [product_id], [region])
VALUES
  (1, 1, 'US');
GO

-- Re-enable constraint checks
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [orders] WITH CHECK CHECK CONSTRAINT ALL;
GO