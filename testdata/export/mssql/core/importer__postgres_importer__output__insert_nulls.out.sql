CREATE TABLE [test_table] (
  [id] INT PRIMARY KEY,
  [name] VARCHAR(100),
  [email] VARCHAR(100),
  [age] INT,
  [active] BOOLEAN DEFAULT (true)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [test_table] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [test_table] ([id], [name], [email], [age], [active])
VALUES
  (1, 'Alice', 'alice@example.com', 25, 1),
  (2, 'Bob', NULL, 30, 0),
  (3, NULL, NULL, NULL, NULL),
  (4, 'Charlie', NULL, NULL, 1);
GO

-- Re-enable constraint checks
ALTER TABLE [test_table] WITH CHECK CHECK CONSTRAINT ALL;
GO