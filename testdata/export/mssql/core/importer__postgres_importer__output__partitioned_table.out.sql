CREATE TABLE [id_to_uuid_migration_map] (
  [column_name] key_column_name_type NOT NULL,
  [int_key] bigint NOT NULL,
  [uuid_key] uuid,
  PRIMARY KEY ([column_name], [int_key])
)
GO

CREATE UNIQUE INDEX [id_to_uuid_migration__uuid_uniq] ON [id_to_uuid_migration_map] ("column_name", "uuid_key")
GO
