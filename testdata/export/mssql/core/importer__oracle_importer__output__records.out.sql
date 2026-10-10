CREATE TABLE [departments] (
  [id] integer PRIMARY KEY,
  [name] nvarchar2(255),
  [budget] number(10, 2)
)
GO

CREATE TABLE [employees] (
  [id] integer PRIMARY KEY,
  [dept_id] integer,
  [name] nvarchar2(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [departments] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [employees] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [departments] ([id], [name], [budget])
VALUES
  (1, 'Engineering', 500000),
  (2, 'Marketing', 300000),
  (3, 'Sales', 400000),
  (4, 'HR', 250000);
GO
INSERT INTO [employees] ([id], [dept_id], [name])
VALUES
  (1, 1, 'John Doe'),
  (2, 1, 'Jane Smith'),
  (3, 2, 'Bob Johnson');
GO

-- Re-enable constraint checks
ALTER TABLE [departments] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [employees] WITH CHECK CHECK CONSTRAINT ALL;
GO