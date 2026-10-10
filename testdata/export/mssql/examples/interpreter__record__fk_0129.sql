CREATE TABLE [follows] (
  [following_user_id] integer,
  [followed_user_id] integer,
  [created_at] timestamp
)
GO

CREATE TABLE [users] (
  [id] integer PRIMARY KEY,
  [username] nvarchar(255)
)
GO

ALTER TABLE [follows] ADD FOREIGN KEY ([following_user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [follows] ADD FOREIGN KEY ([followed_user_id]) REFERENCES [users] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [follows] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [follows] ([following_user_id], [followed_user_id], [created_at])
VALUES
  (1, 2, '2026-01-01');
GO

-- Re-enable constraint checks
ALTER TABLE [follows] WITH CHECK CHECK CONSTRAINT ALL;
GO