CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [user_profiles] (
  [id] int PRIMARY KEY,
  [user_id] int UNIQUE,
  [bio] text
)
GO

CREATE TABLE [departments] (
  [id] int PRIMARY KEY,
  [name] nvarchar(255)
)
GO

CREATE TABLE [employees] (
  [id] int PRIMARY KEY,
  [dept_id] int,
  [name] nvarchar(255)
)
GO

ALTER TABLE [users] ADD FOREIGN KEY ([id]) REFERENCES [user_profiles] ([user_id])
GO

ALTER TABLE [employees] ADD FOREIGN KEY ([dept_id]) REFERENCES [departments] ([id])
GO

-- Disable constraint checks for tables with data
ALTER TABLE [users] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [user_profiles] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [departments] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [employees] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [users] ([id], [name])
VALUES
  (1, 'Alice'),
  (2, 'Bob');
GO
INSERT INTO [user_profiles] ([id], [user_id], [bio])
VALUES
  (1, 1, 'Alice''s bio'),
  (2, 3, 'Invalid user');
GO
INSERT INTO [departments] ([id], [name])
VALUES
  (1, 'Engineering');
GO
INSERT INTO [employees] ([id], [dept_id], [name])
VALUES
  (1, 1, 'Alice'),
  (2, 999, 'Bob with invalid dept');
GO

-- Re-enable constraint checks
ALTER TABLE [users] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [user_profiles] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [departments] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [employees] WITH CHECK CHECK CONSTRAINT ALL;
GO