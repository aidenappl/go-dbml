CREATE TABLE [orders] (
  [id] int PRIMARY KEY,
  [total] float,
  [placed] date,
  [status] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [orders] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [orders] ([id], [total], [placed], [status])
VALUES
  (1, 10.5, '2024-03-01', 'paid'),
  (2, 20, '2024-03-02', NULL);
GO

-- Re-enable constraint checks
ALTER TABLE [orders] WITH CHECK CHECK CONSTRAINT ALL;
GO