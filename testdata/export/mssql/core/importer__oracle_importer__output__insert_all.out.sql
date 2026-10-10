CREATE TABLE [products] (
  [id] NUMBER PRIMARY KEY,
  [name] VARCHAR2(100) NOT NULL,
  [category] VARCHAR2(50),
  [price] NUMBER(10, 2)
)
GO

CREATE TABLE [suppliers] (
  [id] NUMBER PRIMARY KEY,
  [name] VARCHAR2(100) NOT NULL,
  [country] VARCHAR2(50)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [products] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [suppliers] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [products] ([id], [name], [category], [price])
VALUES
  (1, 'Laptop', 'Electronics', 999.99),
  (2, 'Mouse', 'Electronics', 29.99),
  (3, 'Keyboard', 'Electronics', 79.99),
  (4, 'Monitor', 'Electronics', 299.99),
  (5, 'Desk', NULL, NULL),
  (6, 'Chair', NULL, 149.99),
  (7, UPPER('tablet'), NULL, 100 * 5),
  (8, 'Headphones', NULL, 50.00 + 49.99);
GO
INSERT INTO [suppliers] ([id], [name], [country])
VALUES
  (1, 'TechCorp', 'USA'),
  (2, 'GlobalTech', 'Germany'),
  (3, 'AsiaSupply', 'China');
GO

-- Re-enable constraint checks
ALTER TABLE [products] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [suppliers] WITH CHECK CHECK CONSTRAINT ALL;
GO