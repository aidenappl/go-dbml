CREATE TABLE [users] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [posts] (
  [created_at] timestamp DEFAULT (now()),
  [title] nvarchar(255),
  [owner_id] int,
  [id] bigint,
  CONSTRAINT [posts.pos] CHECK (id > 0)
)
GO

CREATE INDEX [posts_index_0] ON [posts] ("created_at")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'p',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts';
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([owner_id]) REFERENCES [users] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [posts] ([id], [title], [owner_id], [created_at])
VALUES
  (1, 'a', 1, '2024-01-01T00:00:00'),
  (1, 'b', 2, NULL);
GO

-- Re-enable constraint checks
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO