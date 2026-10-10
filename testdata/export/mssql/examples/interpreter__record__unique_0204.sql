CREATE TABLE [user_settings] (
  [user_id] int,
  [category] nvarchar(255),
  [value] nvarchar(255)
)
GO

CREATE UNIQUE INDEX [user_settings_index_0] ON [user_settings] ("user_id", "category")
GO

-- Disable constraint checks for tables with data
ALTER TABLE [user_settings] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [user_settings] ([user_id], [category], [value])
VALUES
  (1, NULL, 'default'),
  (1, NULL, 'another default'),
  (1, 'theme', 'dark');
GO

-- Re-enable constraint checks
ALTER TABLE [user_settings] WITH CHECK CHECK CONSTRAINT ALL;
GO