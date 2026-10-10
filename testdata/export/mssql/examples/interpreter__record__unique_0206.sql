CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [email] nvarchar(255) UNIQUE,
  [username] nvarchar(255) UNIQUE
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [sku] nvarchar(255) UNIQUE
)
GO

CREATE TABLE [accounts] (
  [id] int PRIMARY KEY,
  [account_number] int UNIQUE
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [accounts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [email], [username])
VALUES
  (1, 'alice@example.com', 'alice'),
  (2, 'bob@example.com', 'bob');
GO
INSERT INTO [products] ([id], [sku])
VALUES
  (1, 'PROD-001'),
  (2, 'PROD-002');
GO
INSERT INTO [accounts] ([id], [account_number])
VALUES
  (1, 0),
  (2, -1),
  (3, 1000);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [accounts] WITH CHECK CHECK CONSTRAINT ALL;
GO