CREATE TABLE [users] (
  [id] int,
  [quantity] int,
  [name] varchar(5),
  [price] decimal(5,2),
  [active] boolean,
  [status] nvarchar(255) CHECK ([status] IN ('active', 'inactive'))
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [quantity], [name], [price], [active], [status])
VALUES
  (invalid, 10, 'VeryLongName', 12345.67, 0, 'invalid_status');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO