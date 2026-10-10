CREATE TABLE [pulse_aggregates] (
  [bucket] int(10) NOT NULL,
  [period] mediumint(8) NOT NULL,
  [type] varchar(255) NOT NULL,
  [key] text NOT NULL,
  [key_hash] binary(16),
  [aggregate] varchar(255) NOT NULL,
  [value] decimal(20,2) NOT NULL,
  [count] int(10)
)
GO

CREATE TABLE [users] (
  [id] serial PRIMARY KEY,
  [email] varchar(255) UNIQUE NOT NULL,
  [first_name] varchar(255) NOT NULL,
  [last_name] varchar(255) NOT NULL,
  [active] boolean NOT NULL,
  [root] boolean NOT NULL DEFAULT (false)
)
GO

CREATE TABLE [user_roles] (
  [id] serial PRIMARY KEY,
  [role_name] varchar(255) UNIQUE NOT NULL,
  [description] varchar(255)
)
GO

CREATE TABLE [user_role_links] (
  [user_id] bigint NOT NULL,
  [role_id] bigint NOT NULL,
  PRIMARY KEY ([user_id], [role_id])
)
GO

CREATE TABLE [table1] (
  [id] int PRIMARY KEY,
  [name] varchar(255) UNIQUE NOT NULL
)
GO

CREATE TABLE [a] (
  [id] varchar(22)
)
GO

CREATE TABLE [example] (
  [id] varchar(255)
)
GO

CREATE UNIQUE INDEX [pulse_aggregates_bucket_period_type_aggregate_key_hash_unique] ON [pulse_aggregates] ("bucket", "period", "type", "aggregate", "key_hash")
GO

CREATE INDEX [pulse_aggregates_period_bucket_index] ON [pulse_aggregates] ("period", "bucket")
GO

CREATE INDEX [pulse_aggregates_type_index] ON [pulse_aggregates] ("type")
GO

CREATE INDEX [pulse_aggregates_period_type_aggregate_bucket_index] ON [pulse_aggregates] ("period", "type", "aggregate", "bucket")
GO

CREATE INDEX [idx_users_email] ON [users] ("email")
GO

CREATE INDEX [table1_index_5] ON [table1] ("name")
GO

ALTER TABLE [user_role_links] ADD CONSTRAINT [user_role_links_fk_1] FOREIGN KEY ([user_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [user_role_links] ADD CONSTRAINT [user_role_links_fk_2] FOREIGN KEY ([role_id]) REFERENCES [user_roles] ([id])
GO
