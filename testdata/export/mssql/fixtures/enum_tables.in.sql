CREATE TABLE [jobs] (
  [id] integer PRIMARY KEY,
  [status] nvarchar(255) CHECK ([status] IN ('created', 'running', 'done', 'failed', 'wait for validation'))
)
GO

CREATE TABLE [orders] (
  [id] int UNIQUE PRIMARY KEY,
  [status] nvarchar(255) CHECK ([status] IN ('created', 'pending', 'processing', 'completed')),
  [created_at] nvarchar(255)
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'This is a column note',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'jobs',
@level2type = N'Column', @level2name = 'status';
GO
