CREATE TABLE [orders] (
  [id] SERIAL PRIMARY KEY IDENTITY(1, 1),
  [user_id] int UNIQUE NOT NULL,
  [status] nvarchar(255) CHECK ([status] IN ('created', 'running', 'done', 'failure')),
  [created_at] nvarchar(255),
  [modified_at] timestamp(2)
)
GO

CREATE TABLE [order_items] (
  [order_id] int,
  [product_id] int,
  [quantity] int DEFAULT (1),
  [created_at] timestamp,
  [modified_at] time
)
GO

CREATE TABLE [products] (
  [id] int,
  [name] nvarchar(255),
  [merchant_id] int NOT NULL,
  [price] int,
  [status] nvarchar(255) CHECK ([status] IN ('Out of Stock', 'In Stock')),
  [created_at] datetime DEFAULT (now()),
  [modified_at] timestamp,
  PRIMARY KEY ([id], [name])
)
GO

CREATE TABLE [users] (
  [id] int PRIMARY KEY,
  [full_name] nvarchar(255),
  [email] nvarchar(255) UNIQUE,
  [gender] nvarchar(255),
  [date_of_birth] nvarchar(255),
  [created_at] nvarchar(255),
  [modified_at] time(2),
  [country_code] int
)
GO

CREATE TABLE [merchants] (
  [id] int PRIMARY KEY,
  [merchant_name] nvarchar(255),
  [country_code] int,
  [created_at] nvarchar(255),
  [modified_at] time,
  [admin_id] int
)
GO

CREATE TABLE [countries] (
  [code] int PRIMARY KEY,
  [name] nvarchar(255),
  [continent_name] nvarchar(255)
)
GO

CREATE TABLE [foo] (
  [bar] text[],
  [bar2] int[1],
  [bar3] int[2][3],
  [bar4] int[],
  [bar5] int[2],
  [bar6] text[8],
  [bar7] text[100],
  [bar8] time(2)[],
  [bar9] time(1)[1],
  [bar10] time(1)[],
  [bar11] time[5],
  [bar12] timestamp(2)[10][2][5],
  [bar13] character varying[],
  [bar14] character varying(25)[][2][],
  [bar15] character varying[76]
)
GO

CREATE TABLE [accounts] (
  [id] integer NOT NULL IDENTITY(1, 1),
  [domain] character varying(100),
  [name] character varying(100),
  [slug] character varying(10),
  PRIMARY KEY ([id])
)
GO

CREATE TABLE [tests] (
  [id] bigint NOT NULL,
  [type] character varying DEFAULT ('testing'::charactervarying)
)
GO

CREATE TABLE [users2] (
  [username] character varying NOT NULL,
  [hashed_password] character varying NOT NULL,
  [full_name] character varying NOT NULL,
  [email] character varying NOT NULL,
  [password_changed_at] timestamp NOT NULL DEFAULT ('0001-01-01 00:00:00+00'::timestampwithtimezone),
  [created_at] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [table1] (
  [field] text NOT NULL,
  [field2] text NOT NULL,
  PRIMARY KEY ([field], [field2])
)
GO

CREATE TABLE [table2] (
  [id] int PRIMARY KEY IDENTITY(1, 1),
  [field] text,
  [field2] text
)
GO

CREATE TABLE [test_table] (
  [name] nvarchar(255) NOT NULL,
  [id] nvarchar(255),
  [linked_name] nvarchar(255) NOT NULL,
  PRIMARY KEY ([name], [id])
)
GO

CREATE INDEX [product_status] ON [products] ("merchant_id", "status")
GO

CREATE UNIQUE INDEX [products_index_1] ON [products] ("id")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Products table comment',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'products';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Store user data',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'users';
GO

ALTER TABLE [order_items] ADD FOREIGN KEY ([order_id]) REFERENCES [orders] ([id])
GO

ALTER TABLE [order_items] ADD FOREIGN KEY ([product_id]) REFERENCES [products] ([id])
GO

ALTER TABLE [users] ADD FOREIGN KEY ([country_code]) REFERENCES [countries] ([code])
GO

ALTER TABLE [merchants] ADD FOREIGN KEY ([country_code]) REFERENCES [countries] ([code])
GO

ALTER TABLE [products] ADD FOREIGN KEY ([merchant_id]) REFERENCES [merchants] ([id])
GO

ALTER TABLE [merchants] ADD FOREIGN KEY ([admin_id]) REFERENCES [users] ([id])
GO

ALTER TABLE [table2] ADD FOREIGN KEY ([field], [field2]) REFERENCES [table1] ([field], [field2])
GO

ALTER TABLE [test_table] ADD CONSTRAINT [test_table_fk] FOREIGN KEY ([linked_name], [id]) REFERENCES [test_table] ([name], [id])
GO
