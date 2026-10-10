CREATE TABLE [articles] (
  [id] int,
  [content] text,
  [description] nvarchar(255)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [articles] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [articles] ([id], [content], [description])
VALUES
  (1, 'This is a very long text content that can be arbitrarily long without any length restrictions because text type does not have a length parameter', 'This is a very long description that can be arbitrarily long');
GO

-- Re-enable constraint checks
ALTER TABLE [articles] WITH CHECK CHECK CONSTRAINT ALL;
GO