CREATE TABLE [user_profiles] (
  [user_id] int,
  [profile_type] nvarchar(255),
  [data] text
)
GO

CREATE UNIQUE INDEX [user_profiles_index_0] ON [user_profiles] ("user_id", "profile_type")
GO

-- Disable constraint checks for tables with data
ALTER TABLE [user_profiles] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [user_profiles] ([user_id], [profile_type], [data])
VALUES
  (1, 'work', 'Software Engineer'),
  (1, 'work', 'Updated job title');
GO

-- Re-enable constraint checks
ALTER TABLE [user_profiles] WITH CHECK CHECK CONSTRAINT ALL;
GO