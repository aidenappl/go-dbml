CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [categories] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY,
  [title] nvarchar(255),
  [user_id] int,
  [category_id] int
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([category_id]) REFERENCES [categories] ([id])
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [categories] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice');
GO
INSERT INTO [categories] ([id], [name])
VALUES
  (1, 'Tech');
GO
INSERT INTO [posts] ([id], [title], [user_id], [category_id])
VALUES
  (1, 'Valid Post', 1, 1),
  (2, 'Invalid Category', 1, 999),
  (3, 'Invalid User', 999, 1);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [categories] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO