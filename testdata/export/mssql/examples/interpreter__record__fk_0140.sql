CREATE TABLE [students] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [courses] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [students_courses] (
  [students_id] int,
  [courses_id] int,
  PRIMARY KEY ([students_id], [courses_id])
);
GO

ALTER TABLE [students_courses] ADD FOREIGN KEY ([students_id]) REFERENCES [students] ([id]);
GO

ALTER TABLE [students_courses] ADD FOREIGN KEY ([courses_id]) REFERENCES [courses] ([id]);
GO


-- Disable constraint checks for tables with data
ALTER TABLE [students] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [courses] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [students] ([id])
VALUES
  (1),
  (2),
  (3);
GO
INSERT INTO [courses] ([id])
VALUES
  (10),
  (20);
GO

-- Re-enable constraint checks
ALTER TABLE [students] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [courses] WITH CHECK CHECK CONSTRAINT ALL;
GO