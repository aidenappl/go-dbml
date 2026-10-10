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

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [user_id] int,
  [title] nvarchar(255)
)
GO

CREATE TABLE [cities] (
  [id] int PRIMARY KEY,
  [country_code] varchar(2),
  [name] nvarchar(255)
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [cities] ADD FOREIGN KEY ([country_code]) REFERENCES [countries] ([code])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [cities] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice');
GO
INSERT INTO [countries] ([code], [name])
VALUES
  ('US', 'United States');
GO
INSERT INTO [posts] ([id], [user_id], [title])
VALUES
  (1, 1, 'Valid Post'),
  (2, 999, 'Invalid FK');
GO
INSERT INTO [cities] ([id], [country_code], [name])
VALUES
  (1, 'US', 'New York'),
  (2, 'FR', 'Paris');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [cities] WITH CHECK CHECK CONSTRAINT ALL;
GO