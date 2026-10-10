CREATE TABLE [Bảng Người Dùng] (
  [mã số] int PRIMARY KEY,
  [名前] nvarchar(255),
  [emoji 😀] text
)
GO

CREATE TABLE [order items] (
  [user id] int
)
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'emoji 🎉 inside',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'Bảng Người Dùng',
@level2type = N'Column', @level2name = 'emoji 😀';
GO

ALTER TABLE [order items] ADD FOREIGN KEY ([user id]) REFERENCES [Bảng Người Dùng] ([mã số])
GO
