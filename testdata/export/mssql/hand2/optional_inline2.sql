CREATE TABLE [x] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [y] (
  [x_id] int,
  [x2] int
)
GO

CREATE TABLE [z] (
  [x3] int
)
GO

ALTER TABLE [y] ADD FOREIGN KEY ([x_id]) REFERENCES [x] ([id])
GO

ALTER TABLE [y] ADD FOREIGN KEY ([x2]) REFERENCES [x] ([id])
GO

CREATE TABLE [x_z] (
  [x_id] int,
  [z_x3] int,
  PRIMARY KEY ([x_id], [z_x3])
);
GO

ALTER TABLE [x_z] ADD FOREIGN KEY ([x_id]) REFERENCES [x] ([id]);
GO

ALTER TABLE [x_z] ADD FOREIGN KEY ([z_x3]) REFERENCES [z] ([x3]);
GO

