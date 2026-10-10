CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255) UNIQUE
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [sku] nvarchar(255) UNIQUE
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [email])
VALUES
  (1, 'alice@example.com'),
  (2, 'alice@example.com');
GO
INSERT INTO [products] ([id], [sku])
VALUES
  (1, 'PROD-001'),
  (2, 'PROD-001');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO