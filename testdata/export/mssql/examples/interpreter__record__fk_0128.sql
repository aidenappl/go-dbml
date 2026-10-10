CREATE TABLE [users] (
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

CREATE TABLE [employees] (
  [id] int PRIMARY KEY,
  [manager_id] int,
  [name] nvarchar(255)
)
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [employees] ADD FOREIGN KEY ([manager_id]) REFERENCES [employees] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [employees] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice');
GO
INSERT INTO [posts] ([id], [user_id], [title])
VALUES
  (1, 1, 'Valid Post'),
  (2, 999, 'Invalid Post');
GO
INSERT INTO [employees] ([id], [manager_id], [name])
VALUES
  (1, NULL, 'CEO'),
  (2, 1, 'Manager'),
  (3, 999, 'Invalid Manager Reference');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [employees] WITH CHECK CHECK CONSTRAINT ALL;
GO