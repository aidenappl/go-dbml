CREATE TABLE [products] (
  [id] integer PRIMARY KEY,
  [name] nvarchar(255),
  [price] decimal(10,2),
  [in_stock] bit
)
GO

CREATE TABLE [orders] (
  [id] integer PRIMARY KEY,
  [product_id] integer,
  [quantity] integer
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [orders] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [products] ([id], [name], [price], [in_stock])
VALUES
  (1, 'Laptop', 999.99, 1),
  (2, 'Mouse', 29.99, 1),
  (3, 'Keyboard', 79.99, 0),
  (4, 'Monitor', 149.99, 1);
GO
INSERT INTO [orders] ([id], [product_id], [quantity])
VALUES
  (1, 1, 2),
  (2, 2, 5),
  (3, 1, 1);
GO

-- Re-enable constraint checks
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [orders] WITH CHECK CHECK CONSTRAINT ALL;
GO