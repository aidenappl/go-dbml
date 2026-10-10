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

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice'),
  (1, 'Bob');
GO
INSERT INTO [countries] ([code], [name])
VALUES
  ('US', 'United States'),
  ('US', 'USA');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] WITH CHECK CHECK CONSTRAINT ALL;
GO