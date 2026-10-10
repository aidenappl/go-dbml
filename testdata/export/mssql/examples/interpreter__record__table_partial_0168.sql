CREATE TABLE [products] (
  [id] int PRIMARY KEY,
  [created_at] timestamp,
  [name] nvarchar(255),
  [price] decimal
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [products] ([id], [name], [price])
VALUES
  (1, 'Widget', 9.99),
  (2, 'Gadget', 19.99);
GO

-- Re-enable constraint checks
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO