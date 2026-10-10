CREATE TABLE [data] (
  [id] int,
  [quantity] bigint,
  [serial_num] smallint,
  [price] decimal(10,2),
  [rate] float,
  [amount] double
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [quantity], [serial_num], [price], [rate], [amount])
VALUES
  (1, 1000, 5, 99.99, 3.14159, 101325.5),
  (2, -500, -10, -50, -2.5, -40),
  (3, 0, 0, 0, 0, 0.001);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO