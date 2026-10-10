CREATE TABLE [data] (
  [col1] nvarchar(255),
  [col2] boolean,
  [col3] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [data] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [data] ([col1], [col2], [col3])
VALUES
  (NULL, 1, 'false'),
  (NULL, 1, 'false'),
  (NULL, 1, 'false'),
  (NULL, 1, 'false');
GO

-- Re-enable constraint checks
ALTER TABLE [data] WITH CHECK CHECK CONSTRAINT ALL;
GO