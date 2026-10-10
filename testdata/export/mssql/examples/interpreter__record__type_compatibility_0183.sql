CREATE TABLE [users] (
  [id] int,
  [name] varchar(5),
  [code] char(3),
  [title] nvarchar(10),
  [shortcode] nchar(4),
  [description] character varying(8)
)
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name], [code], [title], [shortcode], [description])
VALUES
  (1, 'Alice', 'ABC', 'Manager', 'MGMT', 'Valid'),
  (2, 'VeryLongName', 'ABCD', 'VeryLongTitle', 'TOOLONG', 'TooLongDesc');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO