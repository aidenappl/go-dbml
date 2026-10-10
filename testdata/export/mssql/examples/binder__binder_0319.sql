CREATE TABLE [orders] (
  [id] int PRIMARY KEY,
  [status] nvarchar(255) CHECK ([status] IN ('pending', 'processing', 'completed')) DEFAULT 'pending'
)
GO
