CREATE TABLE [users] (
  [id] integer PRIMARY KEY,
  [name] nvarchar(255),
  [active] boolean
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [active])
VALUES
  (1, 'Alice', NULL),
  (2, 'Bob', NULL),
  (NULL, 'Charlie', 1),
  (3, NULL, 0);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO