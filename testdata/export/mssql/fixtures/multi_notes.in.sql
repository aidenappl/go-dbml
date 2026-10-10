CREATE TABLE [orders] (
  [id] int PRIMARY KEY,
  [status] nvarchar(255) CHECK ([status] IN ('created', 'pending', 'processing', 'completed')),
  [created_at] nvarchar(255)
)
GO

CREATE TABLE [bookings] (
  [id] integer,
  [country] nvarchar(255),
  [booking_date] date,
  [created_at] timestamp,
  PRIMARY KEY ([id], [country])
)
GO

CREATE INDEX [created_at_index] ON [bookings] ("created_at")
GO

CREATE INDEX [bookings_index_1] ON [bookings] ("booking_date")
GO

CREATE UNIQUE INDEX [bookings_index_2] ON [bookings] ("country", "booking_date")
GO

CREATE INDEX [bookings_index_3] ON [bookings] ("booking_date")
GO

CREATE INDEX [bookings_index_4] ON [bookings] ((id*2))
GO

CREATE INDEX [bookings_index_5] ON [bookings] ((id*3), (getdate()))
GO

CREATE INDEX [bookings_index_6] ON [bookings] ((id*3), "id")
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Note on table orders',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'primary key field',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'orders',
@level2type = N'Column', @level2name = 'id';
GO
