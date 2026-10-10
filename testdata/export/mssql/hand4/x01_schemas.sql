CREATE SCHEMA [app]
GO

CREATE SCHEMA [other]
GO

CREATE TABLE [app].[users] (
  [id] int PRIMARY KEY IDENTITY(1, 1),
  [st] nvarchar(255) CHECK ([st] IN ('active', 'on hold')) NOT NULL,
  [st2] nvarchar(255) CHECK ([st2] IN ('a', 'b')),
  [cur] public.money_t,
  [cur2] other.money_t,
  [name] nvarchar(255) DEFAULT 'x'
)
GO

CREATE TABLE [posts] (
  [id] bigint PRIMARY KEY IDENTITY(1, 1),
  [user_id] int,
  [t] timestamp with time zone DEFAULT (now())
)
GO

CREATE TABLE [other].[tags] (
  [id] smallint PRIMARY KEY IDENTITY(1, 1),
  [name] My Type
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'It''s users',
@level0type = N'Schema', @level0name = 'app',
@level1type = N'Table',  @level1name = 'users';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Bob''s name',
@level0type = N'Schema', @level0name = 'app',
@level1type = N'Table',  @level1name = 'users',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'posts',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'posts';
GO

ALTER TABLE [posts] ADD FOREIGN KEY ([user_id]) REFERENCES [app].[users] ([id])
GO

CREATE TABLE [app].[users_tags] (
  [users_id] int,
  [tags_id] smallint,
  PRIMARY KEY ([users_id], [tags_id])
);
GO

ALTER TABLE [app].[users_tags] ADD FOREIGN KEY ([users_id]) REFERENCES [app].[users] ([id]);
GO

ALTER TABLE [app].[users_tags] ADD FOREIGN KEY ([tags_id]) REFERENCES [other].[tags] ([id]);
GO


CREATE TABLE [posts_tags] (
  [posts_id] bigint,
  [tags_id] smallint,
  PRIMARY KEY ([posts_id], [tags_id])
);
GO

ALTER TABLE [posts_tags] ADD FOREIGN KEY ([posts_id]) REFERENCES [posts] ([id]);
GO

ALTER TABLE [posts_tags] ADD FOREIGN KEY ([tags_id]) REFERENCES [other].[tags] ([id]);
GO


ALTER TABLE [other].[tags] ADD CONSTRAINT [r1] FOREIGN KEY ([id]) REFERENCES [posts] ([user_id]) ON DELETE CASCADE ON UPDATE SET NULL
GO
