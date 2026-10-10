CREATE TABLE [data] (
  [id] int,
  [active] boolean
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([id], [active])
VALUES
  (1, 1),
  (2, 0),
  (3, 1),
  (4, 0),
  (5, 1),
  (6, 0),
  (7, 1),
  (8, 0),
  (9, 1),
  (10, 0),
  (11, 1),
  (12, 0),
  (13, 1),
  (14, 0);
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO