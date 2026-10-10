CREATE TABLE [users] (
  [id] integer PRIMARY KEY,
  [name] varchar(255),
  [email] varchar(255),
  [active] boolean
)
GO

CREATE TABLE [comments] (
  [id] integer PRIMARY KEY,
  [user_id] integer,
  [content] text
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [comments] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [email], [active])
VALUES
  (1, 'Alice', 'alice@example.com', 1),
  (2, 'Bob', 'bob@example.com', 0),
  (3, 'Charlie', 'charlie@example.com', 1),
  (4, 'Dave', 'dave@example.com', 0);
GO
INSERT INTO [comments] ([id], [user_id], [content])
VALUES
  (1, 1, 'Great post!'),
  (2, 2, 'Nice article'),
  (3, 1, 'Thanks for sharing');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [comments] WITH CHECK CHECK CONSTRAINT ALL;
GO