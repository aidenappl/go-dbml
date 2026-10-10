CREATE TABLE [categories] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [category_id] int,
  [name] nvarchar(255)
)
GO

ALTER TABLE [products] ADD FOREIGN KEY ([category_id]) REFERENCES [categories] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [categories] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [categories] ([id], [name])
VALUES
  (1, 'Electronics');
GO
INSERT INTO [products] ([id], [category_id], [name])
VALUES
  (1, 1, 'Laptop'),
  (2, NULL, 'Uncategorized Item');
GO

-- Re-enable constraint checks
ALTER TABLE [categories] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO