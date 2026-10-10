CREATE TABLE [nodes] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255),
  [parent_id] int
)
GO

ALTER TABLE [nodes] ADD FOREIGN KEY ([parent_id]) REFERENCES [nodes] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [nodes] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [nodes] ([id], [name], [parent_id])
VALUES
  (1, 'Root', NULL),
  (2, 'Child 1', 1),
  (3, 'Invalid Child', 999);
GO

-- Re-enable constraint checks
ALTER TABLE [nodes] WITH CHECK CHECK CONSTRAINT ALL;
GO