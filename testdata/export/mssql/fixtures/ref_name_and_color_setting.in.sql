CREATE TABLE [b] (
  [id] int,
  [c_id] int
)
GO

CREATE TABLE [c] (
  [id] int,
  [b_id] int
)
GO

ALTER TABLE [c] ADD CONSTRAINT [short_ref] FOREIGN KEY ([id]) REFERENCES [b] ([id])
GO

ALTER TABLE [b] ADD CONSTRAINT [long_ref] FOREIGN KEY ([c_id]) REFERENCES [c] ([id])
GO
