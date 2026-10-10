CREATE TABLE [order_items] (
  [order_id] int,
  [product_id] int,
  [quantity] int,
  PRIMARY KEY ([order_id], [product_id])
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [order_items] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [order_items] ([order_id], [product_id], [quantity])
VALUES
  (1, 100, 2),
  (1, NULL, 3),
  (NULL, 101, 4);
GO

-- Re-enable constraint checks
ALTER TABLE [order_items] WITH CHECK CHECK CONSTRAINT ALL;
GO