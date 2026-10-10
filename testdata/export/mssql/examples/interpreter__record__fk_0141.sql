CREATE TABLE [tags] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [posts] (
  [id] int PRIMARY KEY
)
GO

CREATE TABLE [tags_posts] (
  [tags_id] int,
  [posts_id] int,
  PRIMARY KEY ([tags_id], [posts_id])
);
GO

ALTER TABLE [tags_posts] ADD FOREIGN KEY ([tags_id]) REFERENCES [tags] ([id]);
GO

ALTER TABLE [tags_posts] ADD FOREIGN KEY ([posts_id]) REFERENCES [posts] ([id]);
GO


-- Disable constraint checks for tables with data
ALTER TABLE [tags] NOCHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] NOCHECK CONSTRAINT ALL;
GO

INSERT INTO [tags] ([id])
VALUES
  (1);
GO
INSERT INTO [posts] ([id])
VALUES
  (10),
  (20);
GO

-- Re-enable constraint checks
ALTER TABLE [tags] WITH CHECK CHECK CONSTRAINT ALL;
GO
ALTER TABLE [posts] WITH CHECK CHECK CONSTRAINT ALL;
GO