CREATE TABLE [users] (
  [status] nvarchar(255) CHECK ([status] IN ('value', 'other')) DEFAULT 'value'
)
GO
