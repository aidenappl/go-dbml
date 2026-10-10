CREATE TABLE [users] (
  [id] int IDENTITY(1, 1),
  [tenant_id] int,
  PRIMARY KEY ([id], [tenant_id])
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([tenant_id])
VALUES
  (1),
  (1);
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO