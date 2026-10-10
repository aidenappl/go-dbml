CREATE TABLE [father] (
  [obj_id] char(50) UNIQUE PRIMARY KEY
)
GO

CREATE TABLE [child] (
  [obj_id] char(50) UNIQUE PRIMARY KEY,
  [father_obj_id] char(50)
)
GO

ALTER TABLE [child] ADD FOREIGN KEY ([father_obj_id]) REFERENCES [father] ([obj_id])
GO
