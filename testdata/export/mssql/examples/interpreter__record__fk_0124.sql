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

CREATE TABLE [transactions] (
  [id] int PRIMARY KEY,
  [account_id] int,
  [amount] decimal
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [cities] ADD FOREIGN KEY ([country_code]) REFERENCES [countries] ([code])
GO

ALTER TABLE [transactions] ADD FOREIGN KEY ([account_id]) REFERENCES [accounts] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [accounts] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [cities] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [transactions] NOCHECK CONSTRAINT ALL;
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
INSERT INTO [posts] ([id], [user_id], [title])
VALUES
  (1, 1, 'Alice Post'),
  (2, 1, 'Another Post'),
  (3, 2, 'Bob Post');
GO
INSERT INTO [cities] ([id], [country_code], [name])
VALUES
  (1, 'US', 'New York'),
  (2, 'UK', 'London');
GO
INSERT INTO [transactions] ([id], [account_id], [amount])
VALUES
  (1, 0, 100),
  (2, 1, 50),
  (3, 2, 25);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [countries] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [accounts] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [cities] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [transactions] WITH CHECK CHECK CONSTRAINT ALL;
GO