CREATE TABLE [api_key_outer_auth] (
  [id] bigint PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [api_key_outer_id] bigint NOT NULL,
  [flag] varchar(255) NOT NULL,
  [created_at] datetime
)
GO

CREATE TABLE [entity] (
  [id] bigint PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [dateCreated] timestamp DEFAULT (now()),
  [dateEdited] timestamp DEFAULT (now()),
  [createdBy] varchar(40) NOT NULL DEFAULT '',
  [editedBy] varchar(40) NOT NULL DEFAULT '',
  [constituentId] varchar(40) NOT NULL,
  [name] varchar(255) NOT NULL,
  [ein] varchar(16) NOT NULL
)
GO

CREATE TABLE [entity2] (
  [id] bigint PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [dateCreated] timestamp DEFAULT (now()),
  [dateEdited] timestamp DEFAULT (now()),
  [createdBy] varchar(40) NOT NULL DEFAULT '',
  [editedBy] varchar(40) NOT NULL DEFAULT '',
  [constituentId] varchar(40) NOT NULL,
  [name] varchar(255) NOT NULL DEFAULT 'hello',
  [ein] varchar(16) NOT NULL
)
GO

CREATE UNIQUE INDEX [UQE_api_key_outer_auth_api_key_outer_auth_af] ON [api_key_outer_auth] ("api_key_outer_id", "flag")
GO

CREATE INDEX [IDX_api_key_outer_auth_id] ON [api_key_outer_auth] ("id")
GO

CREATE INDEX [IDX_api_key_outer_auth_api_key_outer_id] ON [api_key_outer_auth] ("api_key_outer_id")
GO

CREATE INDEX [IDX_api_key_outer_auth_flag] ON [api_key_outer_auth] ("flag")
GO

CREATE INDEX [IDX_api_key_outer_auth_created_at] ON [api_key_outer_auth] ("created_at")
GO

CREATE UNIQUE INDEX [id] ON [entity] ("id")
GO

CREATE UNIQUE INDEX [constituentId] ON [entity] ("constituentId")
GO

CREATE UNIQUE INDEX [ein] ON [entity] ("ein")
GO

CREATE INDEX [entity_index_0] ON [entity] ("constituentId")
GO

CREATE INDEX [entity_index_1] ON [entity] ("ein")
GO

CREATE UNIQUE INDEX [id] ON [entity2] ("id")
GO

CREATE UNIQUE INDEX [constituentId] ON [entity2] ("constituentId")
GO

CREATE UNIQUE INDEX [ein] ON [entity2] ("ein")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The database table row id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date / time created',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity',
@level2type = N'Column', @level2name = 'dateCreated';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date / time last modified',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity',
@level2type = N'Column', @level2name = 'dateEdited';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The person who created the record',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity',
@level2type = N'Column', @level2name = 'createdBy';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The person who last modified the record',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity',
@level2type = N'Column', @level2name = 'editedBy';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The globally unique identifier for this constituent',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity',
@level2type = N'Column', @level2name = 'constituentId';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The database table row id',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity2',
@level2type = N'Column', @level2name = 'id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date / time created',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity2',
@level2type = N'Column', @level2name = 'dateCreated';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date / time last modified',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity2',
@level2type = N'Column', @level2name = 'dateEdited';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The person who created the record',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity2',
@level2type = N'Column', @level2name = 'createdBy';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The person who last modified the record',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity2',
@level2type = N'Column', @level2name = 'editedBy';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The globally unique identifier for this constituent',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entity2',
@level2type = N'Column', @level2name = 'constituentId';
GO
