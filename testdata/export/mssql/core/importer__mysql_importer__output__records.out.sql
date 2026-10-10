CREATE TABLE [users] (
  [id] integer PRIMARY KEY,
  [name] varchar(255),
  [email] varchar(255)
)
GO

CREATE TABLE [posts] (
  [id] integer PRIMARY KEY,
  [user_id] integer,
  [title] varchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [email])
VALUES
  (1, 'Alice', 'alice@example.com'),
  (2, 'Bob', 'bob@example.com'),
  (3, 'Charlie', 'charlie@example.com'),
  (4, 'Dave', 'dave@example.com');
GO
INSERT INTO [posts] ([id], [user_id], [title])
VALUES
  (1, 1, 'First Post'),
  (2, 1, 'Second Post'),
  (3, 2, 'Bob Post');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO