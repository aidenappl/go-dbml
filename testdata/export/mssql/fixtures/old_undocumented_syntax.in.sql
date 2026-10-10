CREATE TABLE [customer] (
  [customer_id] SMALLINT PRIMARY KEY NOT NULL IDENTITY(1, 1),
  [store_id] TINYINT NOT NULL,
  [first_name] VARCHAR(45) NOT NULL,
  [last_name] VARCHAR(45) NOT NULL DEFAULT (false),
  [email] VARCHAR(50),
  [address_id] SMALLINT NOT NULL,
  [active] BOOLEAN NOT NULL DEFAULT (true),
  [create_date] DATETIME NOT NULL,
  [last_update] TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
)
GO

CREATE TABLE [cities] (
  [id] integer PRIMARY KEY,
  [name] e DEFAULT 'hello',
  [country_id] integer
)
GO

CREATE TABLE [country] (
  [id] integer,
  [cities] string[]
)
GO

CREATE TABLE [citites] (
  [id] integer,
  [name] string
)
GO

CREATE INDEX [citites_index_0] ON [citites] ("id", "name")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'sasasa',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'cities';
GO
