CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [countries] (
  [code] varchar(2) PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [accounts] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [accounts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice'),
  (2, 'Bob');
GO
INSERT INTO [countries] ([code], [name])
VALUES
  ('US', 'United States'),
  ('UK', 'United Kingdom');
GO
INSERT INTO [accounts] ([id], [name])
VALUES
  (0, 'Default Account'),
  (1, 'User Account'),
  (2, 'Business Account');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [accounts] WITH CHECK CHECK CONSTRAINT ALL;
GO