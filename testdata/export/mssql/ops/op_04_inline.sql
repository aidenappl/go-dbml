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

CREATE TABLE [b_a] (
  [b_id] int,
  [a_b_id] int,
  PRIMARY KEY ([b_id], [a_b_id])
);
GO

ALTER TABLE [b_a] ADD FOREIGN KEY ([b_id]) REFERENCES [b] ([id]);
GO

ALTER TABLE [b_a] ADD FOREIGN KEY ([a_b_id]) REFERENCES [a] ([b_id]);
GO


CREATE TABLE [a_a] (
  [a_id] int,
  [a_self] int,
  PRIMARY KEY ([a_id], [a_self])
);
GO

ALTER TABLE [a_a] ADD FOREIGN KEY ([a_id]) REFERENCES [a] ([id]);
GO

ALTER TABLE [a_a] ADD FOREIGN KEY ([a_self]) REFERENCES [a] ([self]);
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