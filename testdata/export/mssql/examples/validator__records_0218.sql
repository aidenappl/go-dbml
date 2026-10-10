CREATE TABLE [data] (
  [int_col] int,
  [float_col] decimal(10,2),
  [bool_col] boolean,
  [str_col] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([int_col], [float_col], [bool_col], [str_col])
VALUES
  (1, 3.14, 1, 'hello'),
  (2, -2.5, 0, 'world');
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO