CREATE TABLE [b] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [a] (
  [id] int PRIMARY KEY,
  [b_id] int,
  [self] int
)
GO

ALTER TABLE [a] ADD FOREIGN KEY ([b_id]) REFERENCES [b] ([id])
GO

ALTER TABLE [a] ADD FOREIGN KEY ([self]) REFERENCES [a] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [b] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [a] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [b] ([id])
VALUES
  (1);
GO
INSERT INTO [a] ([id], [b_id], [self])
VALUES
  (1, 1, 1),
  (2, 9, NULL);
GO

-- Re-enable constraint checks
ALTER TABLE [b] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [a] WITH CHECK CHECK CONSTRAINT ALL;
GO